package main

import (
	"fmt"

	"github.com/bcicen/ctop/config"
	"github.com/bcicen/ctop/widgets"
	ui "github.com/gizak/termui"
)

type Grid struct {
	cursorID   string // id of currently selected container
	cmap       *ContainerMap
	maxRows    int
	containers Containers // sorted slice of containers
	header     *widgets.CTopHeader
}

func NewGrid() *Grid {
	cmap := NewContainerMap()
	g := &Grid{
		cmap:       cmap,
		containers: cmap.All(),
		header:     widgets.NewCTopHeader(),
	}
	// set initial cursor position
	if len(g.containers) > 0 {
		g.cursorID = g.containers[0].id
		g.containers[0].widgets.Highlight()
	}
	return g
}

func (g *Grid) calcMaxRows() {
	g.maxRows = ui.TermHeight() - widgets.CompactHeader.Height - ui.Body.Y
}

// Return current cursor index
func (g *Grid) cursorIdx() int {
	for n, c := range g.containers {
		if c.id == g.cursorID {
			return n
		}
	}
	return 0
}

func (g *Grid) cursorUp() {
	idx := g.cursorIdx()
	// decrement if possible
	if idx <= 0 {
		return
	}
	active := g.containers[idx]
	next := g.containers[idx-1]

	active.widgets.UnHighlight()
	g.cursorID = next.id
	next.widgets.Highlight()

	ui.Render(ui.Body)
}

func (g *Grid) cursorDown() {
	idx := g.cursorIdx()
	// increment if possible
	if idx > (len(g.containers) - 1) {
		return
	}
	if idx >= g.maxRows-1 {
		return
	}
	active := g.containers[idx]
	next := g.containers[idx+1]

	active.widgets.UnHighlight()
	g.cursorID = next.id
	next.widgets.Highlight()
	ui.Render(ui.Body)
}

func (g *Grid) redrawRows() {
	// reinit body rows
	g.calcMaxRows()
	ui.Body.Rows = []*ui.Row{}
	ui.Clear()

	// build layout
	if config.GetSwitchVal("enableHeader") {
		ui.Body.Y = g.header.Height()
		g.header.SetCount(len(g.containers))
		g.header.SetFilter(config.GetVal("filterStr"))
		g.header.Render()
	} else {
		ui.Body.Y = 0
	}
	ui.Body.AddRows(widgets.CompactHeader)
	for n, c := range g.containers.Filter() {
		if n >= g.maxRows {
			break
		}
		ui.Body.AddRows(c.widgets.Row())
	}

	ui.Body.Align()
	resizeIndicator()
	ui.Render(ui.Body)

	// dump aligned widget positions and sizes
	//for i, w := range ui.Body.Rows[1].Cols {
	//log.Infof("w%v: x=%v y=%v w=%v h=%v", i, w.X, w.Y, w.Width, w.Height)
	//}

}

// override Align()'d size for indicator column
func resizeIndicator() {
	xShift := 1
	toWidth := 3
	for _, r := range ui.Body.Rows {
		wDiff := r.Cols[0].Width - (toWidth + xShift)
		// set indicator x, width
		r.Cols[0].SetX(xShift)
		r.Cols[0].SetWidth(toWidth)

		// shift remainder of columns left by wDiff
		for _, c := range r.Cols[1:] {
			c.SetX(c.X - wDiff)
			c.SetWidth(c.Width - wDiff)
		}
	}
}

func (g *Grid) ExpandView() {
	ui.Clear()
	ui.DefaultEvtStream.ResetHandlers()
	defer ui.DefaultEvtStream.ResetHandlers()
	container := g.cmap.Get(g.cursorID)
	container.Expand()
	container.widgets.Render()
	container.Collapse()
}

func logEvent(e ui.Event) {
	var s string
	s += fmt.Sprintf("Type: %s\n", e.Type)
	s += fmt.Sprintf("Path: %s\n", e.Path)
	s += fmt.Sprintf("From: %s\n", e.From)
	s += fmt.Sprintf("To: %s", e.To)
	log.Debugf("new event:\n%s", s)
}

func Display(g *Grid) bool {
	var menu func()
	var expand bool
	var loopIter int

	ui.DefaultEvtStream.Hook(logEvent)

	// initial draw
	g.redrawRows()

	ui.Handle("/sys/kbd/<up>", func(ui.Event) {
		g.cursorUp()
	})
	ui.Handle("/sys/kbd/<down>", func(ui.Event) {
		g.cursorDown()
	})
	ui.Handle("/sys/kbd/<enter>", func(ui.Event) {
		expand = true
		ui.StopLoop()
	})
	ui.Handle("/sys/kbd/a", func(ui.Event) {
		config.Toggle("allContainers")
		g.redrawRows()
	})
	ui.Handle("/sys/kbd/f", func(ui.Event) {
		menu = FilterMenu
		ui.StopLoop()
	})
	ui.Handle("/sys/kbd/h", func(ui.Event) {
		menu = HelpMenu
		ui.StopLoop()
	})
	ui.Handle("/sys/kbd/H", func(ui.Event) {
		config.Toggle("enableHeader")
		g.redrawRows()
	})
	ui.Handle("/sys/kbd/q", func(ui.Event) {
		ui.StopLoop()
	})
	ui.Handle("/sys/kbd/r", func(e ui.Event) {
		config.Toggle("sortReversed")
	})
	ui.Handle("/sys/kbd/s", func(ui.Event) {
		menu = SortMenu
		ui.StopLoop()
	})
	ui.Handle("/timer/1s", func(e ui.Event) {
		loopIter++
		if loopIter%5 == 0 {
			g.cmap.Refresh()
		}
		g.containers = g.cmap.All() // refresh containers for current sort order
		g.redrawRows()
	})

	ui.Handle("/sys/wnd/resize", func(e ui.Event) {
		g.header.Align()
		ui.Body.Width = ui.TermWidth()
		log.Infof("resize: width=%v max-rows=%v", ui.Body.Width, g.maxRows)
		g.redrawRows()
	})

	ui.Loop()
	if menu != nil {
		menu()
		return false
	}
	if expand {
		g.ExpandView()
		return false
	}
	return true
}
