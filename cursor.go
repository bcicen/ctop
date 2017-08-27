package main

import (
	"math"

	"github.com/bcicen/ctop/connector"
	ui "github.com/gizak/termui"
	"github.com/bcicen/ctop/entity"
	"github.com/bcicen/ctop/config"
)

type GridCursor struct {
	selectedID         string // id of currently selected container
	filteredContainers entity.Containers
	filteredNodes      entity.Nodes
	filteredServices   entity.Services
	filteredTasks      entity.Tasks
	cSource            connector.Connector
	isScrolling        bool // toggled when actively scrolling
}

func (gc *GridCursor) LenNodes() int      { return len(gc.filteredNodes) }
func (gc *GridCursor) LenServices() int   { return len(gc.filteredServices) }
func (gc *GridCursor) LenTasks() int      { return len(gc.filteredTasks) }
func (gc *GridCursor) LenContainers() int { return len(gc.filteredContainers) }

func (gc *GridCursor) Selected() (entity.Entity, string) {
	idx, type_entity := gc.Idx()
	if idx < gc.LenEntity(type_entity) {
		return gc.entity(type_entity, idx), type_entity
	}
	return nil, type_entity
}

func (gc *GridCursor) SelectedContainer() *entity.Container {
	idx, _ := gc.Idx()
	if idx < gc.LenContainers() {
		return gc.filteredContainers[idx]
	}
	return nil
}

// Refresh node from source
func (gc *GridCursor) RefreshNodes() (lenChanged bool) {
	oldLen := gc.LenNodes()

	// Containers filtered by display bool
	gc.filteredNodes = entity.Nodes{}
	var cursorVisible bool
	for _, n := range gc.cSource.AllNodes() {
		if n.Display {
			if n.Id == gc.selectedID {
				cursorVisible = true
			}
			gc.filteredNodes = append(gc.filteredNodes, n)
		}
	}

	if oldLen != gc.LenNodes() {
		lenChanged = true
	}

	if !cursorVisible {
		gc.Reset()
	}
	if gc.selectedID == "" {
		gc.Reset()
	}
	return lenChanged
}

func (gc *GridCursor) RefreshTasks() (lenChanged bool) {
	oldLen := gc.LenTasks()

	gc.filteredTasks = entity.Tasks{}
	var cursorVisible bool
	for _, t := range gc.cSource.AllTasks() {
		if t.Display {
			if t.Id == gc.selectedID {
				cursorVisible = true
			}
			gc.filteredTasks = append(gc.filteredTasks, t)
		}
	}

	if oldLen != gc.LenTasks() {
		lenChanged = true
	}

	if !cursorVisible {
		gc.Reset()
	}
	if gc.selectedID == "" {
		gc.Reset()
	}
	return lenChanged
}

func (gc *GridCursor) RefreshServices() (lenChanged bool) {
	oldLen := gc.LenServices()

	gc.filteredServices = entity.Services{}
	var cursorVisible bool
	for _, s := range gc.cSource.AllServices() {
		if s.Display {
			if s.Id == gc.selectedID {
				cursorVisible = true
			}
			gc.filteredServices = append(gc.filteredServices, s)
		}
	}

	if oldLen != gc.LenServices() {
		lenChanged = true
	}

	if !cursorVisible {
		gc.Reset()
	}
	if gc.selectedID == "" {
		gc.Reset()
	}
	return lenChanged
}

func (gc *GridCursor) RefreshContainers() (lenChanged bool) {
	oldLen := gc.LenContainers()

	gc.filteredContainers = entity.Containers{}
	var cursorVisible bool
	for _, c := range gc.cSource.AllContainers() {
		if c.Display {
			if c.Id == gc.selectedID {
				cursorVisible = true
			}
			gc.filteredContainers = append(gc.filteredContainers, c)
		}
	}

	if oldLen != gc.LenContainers() {
		lenChanged = true
	}

	if !cursorVisible {
		gc.Reset()
	}
	if gc.selectedID == "" {
		gc.Reset()
	}
	return lenChanged
}

// Set an initial cursor position, if possible
func (gc *GridCursor) Reset() {
	if config.GetSwitchVal("swarmMode") {
		for _, n := range gc.cSource.AllNodes() {
			n.Widgets.Name.UnHighlight()
		}
		for _, s := range gc.cSource.AllServices() {
			s.Widgets.Name.UnHighlight()
		}
		for _, t := range gc.cSource.AllTasks() {
			t.Widgets.Name.UnHighlight()
		}
		if gc.LenNodes() > 0 {
			gc.selectedID = gc.filteredNodes[0].Id
			gc.filteredNodes[0].Widgets.Name.Highlight()
		}
	} else {
		for _, c := range gc.cSource.AllContainers() {
			c.Widgets.Name.UnHighlight()
		}
		if gc.LenContainers() > 0 {
			gc.selectedID = gc.filteredContainers[0].Id
			gc.filteredContainers[0].Widgets.Name.Highlight()
		}
	}
}

// Return current cursor index
func (gc *GridCursor) Idx() (int, string) {
	if config.GetSwitchVal("swarmMode") {
		for n, node := range gc.filteredNodes {
			if node.Id == gc.selectedID {
				return n, "node"
			}
		}
		for n, service := range gc.filteredServices {
			if service.Id == gc.selectedID {
				return n, "service"
			}
		}
		for n, task := range gc.filteredTasks {
			if task.Id == gc.selectedID {
				return n, "task"
			}
		}
	} else {
		for n, container := range gc.filteredContainers {
			if container.Id == gc.selectedID {
				return n, "container"
			}
		}
	}
	gc.Reset()
	return 0, ""
}

func (gc *GridCursor) ScrollPage() {

	if gc.AllLen() < cGrid.MaxRows() {
		cGrid.Offset = 0
		return
	}

	idx, _ := gc.Idx()

	// page down
	if idx >= cGrid.Offset+cGrid.MaxRows() {
		cGrid.Offset++
		cGrid.Align()
	}
	// page up
	if idx < cGrid.Offset {
		cGrid.Offset--
		cGrid.Align()
	}

}

func (gc *GridCursor) Up() {
	gc.isScrolling = true
	defer func() { gc.isScrolling = false }()

	idx, entity := gc.Idx()
	if idx <= 0 { // already at top
		return
	}
	active := gc.entity(entity, idx)
	next := gc.entity(entity, idx-1)

	active.GetMetaEntity().Widgets.Name.UnHighlight()
	gc.selectedID = next.GetId()
	next.GetMetaEntity().Widgets.Name.Highlight()

	gc.ScrollPage()
	ui.Render(cGrid)
}

func (gc *GridCursor) Down() {
	gc.isScrolling = true
	defer func() { gc.isScrolling = false }()

	idx, entity := gc.Idx()
	if idx >= gc.LenEntity(entity)-1 { // already at bottom
		return
	}

	active := gc.entity(entity, idx)
	next := gc.entity(entity, idx+1)

	active.GetMetaEntity().Widgets.Name.UnHighlight()
	gc.selectedID = next.GetId()
	next.GetMetaEntity().Widgets.Name.Highlight()

	gc.ScrollPage()
	ui.Render(cGrid)
}

func (gc *GridCursor) PgUp() {
	idx, entity := gc.Idx()
	if idx <= 0 { // already at top
		return
	}

	nextidx := int(math.Max(0.0, float64(idx-cGrid.MaxRows())))
	if gc.pgCount() > 0 {
		cGrid.Offset = int(math.Max(float64(cGrid.Offset-cGrid.MaxRows()),
			float64(0)))
	}

	active := gc.entity(entity, idx)
	next := gc.entity(entity, nextidx)

	active.GetMetaEntity().Widgets.Name.UnHighlight()
	gc.selectedID = next.GetId()
	next.GetMetaEntity().Widgets.Name.Highlight()

	cGrid.Align()
	ui.Render(cGrid)
}

func (gc *GridCursor) PgDown() {
	idx, entity := gc.Idx()
	if idx >= gc.LenEntity(entity)-1 { // already at bottom
		return
	}

	nextidx := int(math.Min(float64(gc.LenEntity(entity)-1), float64(idx+cGrid.MaxRows())))
	if gc.pgCount() > 0 {
		cGrid.Offset = int(math.Min(float64(cGrid.Offset+cGrid.MaxRows()),
			float64(gc.LenEntity(entity)-cGrid.MaxRows())))
	}

	active := gc.entity(entity, idx)
	next := gc.entity(entity, nextidx)

	active.GetMetaEntity().Widgets.Name.UnHighlight()
	gc.selectedID = next.GetId()
	next.GetMetaEntity().Widgets.Name.Highlight()

	cGrid.Align()
	ui.Render(cGrid)
}

func (gc *GridCursor) AllLen() int {
	if config.GetSwitchVal("swarmMode") {
		return gc.LenNodes() + gc.LenServices() + gc.LenTasks()
	}
	return gc.LenContainers()
}

func (gc *GridCursor) pgCount() int {
	pages := gc.AllLen() / cGrid.MaxRows()
	if gc.AllLen()%cGrid.MaxRows() > 0 {
		pages++
	}
	return pages
}

func (gc *GridCursor) entity(t string, id int) entity.Entity {
	switch t {
	case "container":
		return gc.filteredContainers[id]
	case "node":
		return gc.filteredNodes[id]
	case "service":
		return gc.filteredServices[id]
	case "task":
		return gc.filteredTasks[id]
	}
	return nil
}

func (gc *GridCursor) LenEntity(t string) int {
	switch t {
	case "node":
		return gc.LenNodes()
	case "service":
		return gc.LenServices()
	case "task":
		return gc.LenTasks()
	default:
		return gc.LenContainers()
	}

}
