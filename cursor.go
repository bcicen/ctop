package main

import (
	"math"

	"github.com/bcicen/ctop/connector"
	"github.com/bcicen/ctop/container"
	ui "github.com/gizak/termui"
)

var enabledConnectors = map[string]func() connector.Connector{
	"docker": connector.NewDocker,
	"runc":   connector.NewRunc,
}

type GridCursor struct {
	selectedID  string // id of currently selected container
	filtered    container.Containers
	cSource     connector.Connector
	isScrolling bool // toggled when actively scrolling
}

func NewGridCursor(connector string) *GridCursor {
	return &GridCursor{
		cSource: enabledConnectors[connector](),
	}
}

func (gc *GridCursor) Len() int { return len(gc.filtered) }

func (gc *GridCursor) Selected() *container.Container {
	idx := gc.Idx()
	if idx < gc.Len() {
		return gc.filtered[idx]
	}
	return nil
}

// Refresh containers from source
func (gc *GridCursor) RefreshContainers() (lenChanged bool) {
	oldLen := gc.Len()

	// Containers filtered by display bool
	gc.filtered = container.Containers{}
	var cursorVisible bool
	for _, c := range gc.cSource.All() {
		if c.Display {
			if c.Id == gc.selectedID {
				cursorVisible = true
			}
			gc.filtered = append(gc.filtered, c)
		}
	}

	if oldLen != gc.Len() {
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
	for _, c := range gc.cSource.All() {
		c.Widgets.Name.UnHighlight()
	}
	if gc.Len() > 0 {
		gc.selectedID = gc.filtered[0].Id
		gc.filtered[0].Widgets.Name.Highlight()
	}
}

// Return current cursor index
func (gc *GridCursor) Idx() int {
	for n, c := range gc.filtered {
		if c.Id == gc.selectedID {
			return n
		}
	}
	gc.Reset()
	return 0
}

func (gc *GridCursor) ScrollPage() {
	// skip scroll if no need to page
	if gc.Len() < cGrid.MaxRows() {
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
	active := gc.filtered[idx]
	next := gc.filtered[idx-1]

	active.Widgets.Name.UnHighlight()
	gc.selectedID = next.Id
	next.Widgets.Name.Highlight()

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
	active := gc.filtered[idx]
	next := gc.filtered[idx+1]

	active.Widgets.Name.UnHighlight()
	gc.selectedID = next.Id
	next.Widgets.Name.Highlight()

	gc.ScrollPage()
	ui.Render(cGrid)
}

func (gc *GridCursor) PgUp() {
	idx := gc.Idx()
	if idx <= 0 { // already at top
		return
	}

	var nextidx int
	nextidx = int(math.Max(0.0, float64(idx-cGrid.MaxRows())))
	cGrid.Offset = int(math.Max(float64(cGrid.Offset-cGrid.MaxRows()),
		float64(0)))

	active := gc.filtered[idx]
	next := gc.filtered[nextidx]

	active.Widgets.Name.UnHighlight()
	gc.selectedID = next.Id
	next.Widgets.Name.Highlight()

	cGrid.Align()
	ui.Render(cGrid)
}

func (gc *GridCursor) PgDown() {
	idx := gc.Idx()
	if idx >= gc.Len()-1 { // already at bottom
		return
	}

	var nextidx int
	nextidx = int(math.Min(float64(gc.Len()-1),
		float64(idx+cGrid.MaxRows())))
	cGrid.Offset = int(math.Min(float64(cGrid.Offset+cGrid.MaxRows()),
		float64(gc.Len()-cGrid.MaxRows())))

	active := gc.filtered[idx]
	next := gc.filtered[nextidx]

	active.Widgets.Name.UnHighlight()
	gc.selectedID = next.Id
	next.Widgets.Name.Highlight()

	cGrid.Align()
	ui.Render(cGrid)
}
