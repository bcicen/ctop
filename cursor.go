package main

import (
	"math"

	"github.com/bcicen/ctop/connector"
	"github.com/bcicen/ctop/container"
	ui "github.com/gizak/termui"
)

type GridCursor struct {
	selectedID  string // id of currently selected container
	filtered    container.Containers
	cSuper      *connector.ConnectorSuper
	isScrolling bool // toggled when actively scrolling
}

func (gc *GridCursor) Len() int { return len(gc.filtered) }

func (gc *GridCursor) Selected() *container.Container {
	idx := gc.Idx()
	if idx < gc.Len() {
		return gc.filtered[idx]
	}
	return nil
}

// Refresh containers from source, returning whether the quantity of
// containers has changed and any error
func (gc *GridCursor) RefreshContainers() (bool, error) {
	oldLen := gc.Len()
	gc.filtered = container.Containers{}

	cSource, err := gc.cSuper.Get()
	if err != nil {
		return true, err
	}

	// filter Containers by display bool
	var cursorVisible bool
	for _, c := range cSource.All() {
		if c.Display {
			if c.Id == gc.selectedID {
				cursorVisible = true
			}
			gc.filtered = append(gc.filtered, c)
		}
	}

	if !cursorVisible || gc.selectedID == "" {
		gc.Reset()
	}

	return oldLen != gc.Len(), nil
}

// Set an initial cursor position, if possible
func (gc *GridCursor) Reset() {
	cSource, err := gc.cSuper.Get()
	if err != nil {
		return
	}

	for _, c := range cSource.All() {
		c.Widgets.UnHighlight()
	}
	if gc.Len() > 0 {
		gc.selectedID = gc.filtered[0].Id
		gc.filtered[0].Widgets.Highlight()
	}
}

// Idx returns current cursor index
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

	active.Widgets.UnHighlight()
	gc.selectedID = next.Id
	next.Widgets.Highlight()

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

	active.Widgets.UnHighlight()
	gc.selectedID = next.Id
	next.Widgets.Highlight()

	gc.ScrollPage()
	ui.Render(cGrid)
}

func (gc *GridCursor) PgUp() {
	idx := gc.Idx()
	if idx <= 0 { // already at top
		return
	}

	nextidx := int(math.Max(0.0, float64(idx-cGrid.MaxRows())))
	if gc.pgCount() > 0 {
		cGrid.Offset = int(math.Max(float64(cGrid.Offset-cGrid.MaxRows()),
			float64(0)))
	}

	active := gc.filtered[idx]
	next := gc.filtered[nextidx]

	active.Widgets.UnHighlight()
	gc.selectedID = next.Id
	next.Widgets.Highlight()

	cGrid.Align()
	ui.Render(cGrid)
}

func (gc *GridCursor) PgDown() {
	idx := gc.Idx()
	if idx >= gc.Len()-1 { // already at bottom
		return
	}

	nextidx := int(math.Min(float64(gc.Len()-1), float64(idx+cGrid.MaxRows())))
	if gc.pgCount() > 0 {
		cGrid.Offset = int(math.Min(float64(cGrid.Offset+cGrid.MaxRows()),
			float64(gc.Len()-cGrid.MaxRows())))
	}

	active := gc.filtered[idx]
	next := gc.filtered[nextidx]

	active.Widgets.UnHighlight()
	gc.selectedID = next.Id
	next.Widgets.Highlight()

	cGrid.Align()
	ui.Render(cGrid)
}

// number of pages at current row count and term height
func (gc *GridCursor) pgCount() int {
	pages := gc.Len() / cGrid.MaxRows()
	if gc.Len()%cGrid.MaxRows() > 0 {
		pages++
	}
	return pages
}
