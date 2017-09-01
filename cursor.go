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
	filteredId         []string
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

func (gc *GridCursor) Selected() (entity.Entity) {
	idx := gc.Idx()
	if idx < gc.Len() {
		return gc.entity(idx)
	}
	return nil
}

// Refresh node from source
//func (gc *GridCursor) RefreshNodes() (lenChanged bool) {
//	oldLen := gc.LenNodes()
//
//	// Containers filtered by display bool
//	gc.filteredNodes = entity.Nodes{}
//	var cursorVisible bool
//	for _, n := range gc.cSource.AllNodes() {
//		if n.Display {
//			if n.Id == gc.selectedID {
//				cursorVisible = true
//			}
//			gc.filteredNodes = append(gc.filteredNodes, n)
//		}
//	}
//
//	if oldLen != gc.LenNodes() {
//		lenChanged = true
//	}
//
//	if !cursorVisible {
//		gc.Reset()
//	}
//	if gc.selectedID == "" {
//		gc.Reset()
//	}
//	return lenChanged
//}

func (gc *GridCursor) RefreshTasks() (lenChanged bool) {
	oldLen := gc.LenTasks()
	gc.filteredId = []string{}

	gc.filteredTasks = entity.Tasks{}
	var cursorVisible bool
	for _, t := range gc.cSource.AllTasks() {
		if t.Display {
			if t.Id == gc.selectedID {
				cursorVisible = true
			}
			gc.filteredTasks = append(gc.filteredTasks, t)
			//gc.addFilteredId(t)
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
	gc.filteredId = []string{}

	gc.filteredServices = entity.Services{}
	var cursorVisible bool
	for _, s := range gc.cSource.AllServices() {
		if s.Display {
			if s.Id == gc.selectedID {
				cursorVisible = true
			}
			gc.filteredServices = append(gc.filteredServices, s)
			//gc.addFilteredId(s)
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
	gc.filteredId = []string{}

	gc.filteredContainers = entity.Containers{}
	var cursorVisible bool
	for _, c := range gc.cSource.AllContainers() {
		if c.Display {
			if c.Id == gc.selectedID {
				cursorVisible = true
			}
			gc.filteredContainers = append(gc.filteredContainers, c)
			//gc.addFilteredId(c)
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
	} else {
		for _, c := range gc.cSource.AllContainers() {
			c.Widgets.Name.UnHighlight()
		}
		//if gc.LenContainers() > 0 {
		//	gc.selectedID = gc.filteredContainers[0].Id
		//	gc.filteredContainers[0].Widgets.Name.Highlight()
		//}
	}
	if gc.Len() > 0 {
		gc.selectedID = gc.idByIndex(0)
		e := gc.entityById(gc.selectedID)
		if e != nil {
			e.GetMetaEntity().Widgets.Name.Highlight()
		}
	}
}

// Return current cursor index
func (gc *GridCursor) Idx() int {
	log.Debugf(">>>>>FilteredId %s", gc.filteredId)
	n := 0
	for _, k := range gc.filteredId {
		if k == gc.selectedID {
			return n
		}
		n += 1
	}
	gc.Reset()
	return 0
}

func (gc *GridCursor) idByIndex(i int) string {
	for _, k := range gc.filteredId {
		if i != 0 {
			i -= 1
		} else {
			return k
		}
	}
	return ""
}

func (gc *GridCursor) AddFilteredId(e entity.Entity) {
	gc.filteredId = append(gc.filteredId, e.GetId())
}

func (gc *GridCursor) ScrollPage() {

	if gc.AllLen() < cGrid.MaxRows() {
		cGrid.Offset = 0
		return
	}

	idx := gc.Idx()

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

	idx := gc.Idx()
	if idx <= 0 { // already at top
		return
	}
	active := gc.entity(idx)
	next := gc.entity(idx - 1)

	active.GetMetaEntity().Widgets.Name.UnHighlight()
	gc.selectedID = next.GetId()
	next.GetMetaEntity().Widgets.Name.Highlight()

	gc.ScrollPage()
	ui.Render(cGrid)
}

func (gc *GridCursor) Down() {
	gc.isScrolling = true
	defer func() { gc.isScrolling = false }()

	idx := gc.Idx()
	if idx >= gc.Len()-1 { // already at bottom
		return
	}

	active := gc.entity(idx)
	next := gc.entity(idx + 1)

	active.GetMetaEntity().Widgets.Name.UnHighlight()
	gc.selectedID = next.GetId()
	next.GetMetaEntity().Widgets.Name.Highlight()

	gc.ScrollPage()
	ui.Render(cGrid)
}

func (gc *GridCursor) PgUp() {
	idx := gc.Idx()
	if idx <= 0 { // already at top
		return
	}

	nextIdx := int(math.Max(0.0, float64(idx-cGrid.MaxRows())))
	if gc.pgCount() > 0 {
		cGrid.Offset = int(math.Max(float64(cGrid.Offset-cGrid.MaxRows()),
			float64(0)))
	}

	active := gc.entity(idx)
	next := gc.entity(nextIdx)

	active.GetMetaEntity().Widgets.Name.UnHighlight()
	gc.selectedID = next.GetId()
	next.GetMetaEntity().Widgets.Name.Highlight()

	cGrid.Align()
	ui.Render(cGrid)
}

func (gc *GridCursor) PgDown() {
	idx := gc.Idx()
	if idx >= gc.Len()-1 { // already at bottom
		return
	}

	nextIdx := int(math.Min(float64(gc.Len()-1), float64(idx+cGrid.MaxRows())))
	if gc.pgCount() > 0 {
		cGrid.Offset = int(math.Min(float64(cGrid.Offset+cGrid.MaxRows()),
			float64(gc.Len()-cGrid.MaxRows())))
	}

	active := gc.entity(idx)
	next := gc.entity(nextIdx)

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

func (gc *GridCursor) entity(id int) entity.Entity {
	cid := gc.idByIndex(id)
	return gc.entityById(cid)
}

func (gc *GridCursor) entityById(cid string) entity.Entity {
	for _, s := range gc.filteredServices {
		if cid == s.Id {
			return s
		}
	}
	for _, t := range gc.filteredTasks {
		if cid == t.Id {
			return t
		}
	}
	for _, c := range gc.filteredContainers {
		if cid == c.Id {
			return c
		}
	}
	if config.GetSwitchVal("swarmMode") {
		return gc.filteredServices[0]
	} else {
		return gc.filteredContainers[0]
	}
}

func (gc *GridCursor) Len() int {
	return len(gc.filteredId)
}
