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
	containers []*Container // sorted slice of containers
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
	}
	return g
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
	if idx > 0 {
		g.cursorID = g.containers[idx-1].id
		g.redrawCursor()
	}
}

func (g *Grid) cursorDown() {
	idx := g.cursorIdx()
	// increment if possible
	if idx < (len(g.containers) - 1) {
		g.cursorID = g.containers[idx+1].id
		g.redrawCursor()
	}
}

// Redraw the cursor with the currently selected row
func (g *Grid) redrawCursor() {
	for _, c := range g.containers {
		if c.id == g.cursorID {
			c.widgets.Highlight()
		} else {
			c.widgets.UnHighlight()
		}
		ui.Render(ui.Body)
	}
}

func (g *Grid) redrawRows() {
	// reinit body rows
	ui.Body.Rows = []*ui.Row{}
	ui.Clear()

	// build layout
	if config.GetToggle("enableHeader") {
		g.header.SetCount(len(g.containers))
		g.header.SetFilter(config.Get("filterStr"))
		ui.Body.AddRows(g.header.Row())
	}
	ui.Body.AddRows(fieldHeader())
	for _, c := range g.containers {
		if !config.GetToggle("allContainers") && c.state != "running" {
			continue
		}
		ui.Body.AddRows(c.widgets.Row())
	}

	ui.Body.Align()
	ui.Render(ui.Body)
}

func fieldHeader() *ui.Row {
	return ui.NewRow(
		ui.NewCol(1, 0, headerPar("STATUS")),
		ui.NewCol(2, 0, headerPar("NAME")),
		ui.NewCol(2, 0, headerPar("CID")),
		ui.NewCol(2, 0, headerPar("CPU")),
		ui.NewCol(2, 0, headerPar("MEM")),
		ui.NewCol(2, 0, headerPar("NET RX/TX")),
	)
}

func headerPar(s string) *ui.Par {
	p := ui.NewPar(fmt.Sprintf(" %s", s))
	p.Border = false
	p.Height = 2
	p.Width = 20
	p.TextFgColor = ui.ColorWhite
	return p
}

func ResetView() {
	ui.DefaultEvtStream.ResetHandlers()
	ui.Clear()
}

func (g *Grid) ExpandView() {
	ResetView()
	defer ResetView()
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

	// calculate layout
	ui.Body.Align()
	g.redrawCursor()
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
		ui.Body.Width = ui.TermWidth()
		ui.Body.Align()
		ui.Clear()
		ui.Render(ui.Body)
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
