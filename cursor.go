package main

import (
	ui "github.com/gizak/termui"
)

type GridCursor struct {
	selectedID string // id of currently selected container
	containers Containers
	cSource    ContainerSource
}

func NewGridCursor() *GridCursor {
	return &GridCursor{
		cSource: NewDockerContainerSource(),
	}
}

func (gc *GridCursor) Len() int             { return len(gc.Filtered()) }
func (gc *GridCursor) Selected() *Container { return gc.containers[gc.Idx()] }

// Return Containers filtered by display bool
func (gc *GridCursor) Filtered() Containers {
	var filtered Containers
	for _, c := range gc.containers {
		if c.display {
			filtered = append(filtered, c)
		}
	}
	return filtered
}

// Refresh containers from source
func (gc *GridCursor) RefreshContainers() (lenChanged bool) {
	oldLen := gc.Len()
	gc.containers = gc.cSource.All()
	if oldLen != gc.Len() {
		lenChanged = true
	}
	if gc.selectedID == "" {
		gc.Reset()
	}
	return lenChanged
}

// Set an initial cursor position, if possible
func (gc *GridCursor) Reset() {
	if gc.Len() > 0 {
		gc.selectedID = gc.containers[0].Id
		gc.containers[0].Widgets.Name.Highlight()
	}
}

// Return current cursor index
func (gc *GridCursor) Idx() int {
	for n, c := range gc.containers {
		if c.Id == gc.selectedID {
			return n
		}
	}
	return 0
}

func (gc *GridCursor) Up() {
	idx := gc.Idx()
	// decrement if possible
	if idx <= 0 {
		return
	}
	active := gc.containers[idx]
	next := gc.containers[idx-1]

	active.Widgets.Name.UnHighlight()
	gc.selectedID = next.Id
	next.Widgets.Name.Highlight()

	ui.Render(cGrid)
}

func (gc *GridCursor) Down() {
	idx := gc.Idx()
	// increment if possible
	if idx >= (gc.Len() - 1) {
		return
	}
	if idx >= maxRows()-1 {
		return
	}
	active := gc.containers[idx]
	next := gc.containers[idx+1]

	active.Widgets.Name.UnHighlight()
	gc.selectedID = next.Id
	next.Widgets.Name.Highlight()
	ui.Render(cGrid)
}
