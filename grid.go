package main

import (
	"fmt"

	"github.com/bcicen/ctop/widgets"
	ui "github.com/gizak/termui"
)

type Grid struct {
	cursorID     string // id of currently selected container
	containers   []*Container
	containerMap *ContainerMap
	header       *widgets.CTopHeader
}

func NewGrid() *Grid {
	containerMap := NewContainerMap()
	containers := containerMap.All()
	return &Grid{
		cursorID:     containers[0].id,
		containers:   containers,
		containerMap: containerMap,
		header:       widgets.NewCTopHeader(),
	}
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

	// build layout
	if GlobalConfig["enableHeader"] == "1" {
		g.header.SetCount(len(g.containers))
		ui.Body.AddRows(g.header.Row())
	}
	ui.Body.AddRows(fieldHeader())
	for _, c := range g.containers {
		ui.Body.AddRows(c.widgets.Row())
	}

	ui.Body.Align()
	ui.Render(ui.Body)
}

func fieldHeader() *ui.Row {
	return ui.NewRow(
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

type View func(*Grid)

func ResetView() {
	ui.DefaultEvtStream.ResetHandlers()
	ui.Clear()
}

func (g *Grid) OpenView(v View) {
	ResetView()
	defer ResetView()
	v(g)
}

func (g *Grid) ExpandView() {
	ResetView()
	defer ResetView()
	container := g.containerMap.Get(g.cursorID)
	container.Expand()
	container.widgets.Render()
	container.Collapse()
}

func Display(g *Grid) bool {
	var newView View
	var expand bool

	// calculate layout
	ui.Body.Align()
	g.redrawCursor()
	ui.Render(ui.Body)

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
	ui.Handle("/sys/kbd/h", func(ui.Event) {
		newView = HelpMenu
		ui.StopLoop()
	})
	ui.Handle("/sys/kbd/s", func(ui.Event) {
		newView = SortMenu
		ui.StopLoop()
	})
	ui.Handle("/sys/kbd/q", func(ui.Event) {
		ui.StopLoop()
	})
	ui.Handle("/timer/1s", func(e ui.Event) {
		g.containers = g.containerMap.All() // refresh containers for current sort order
		g.redrawRows()
	})

	ui.Handle("/sys/wnd/resize", func(e ui.Event) {
		ui.Body.Width = ui.TermWidth()
		ui.Body.Align()
		ui.Clear()
		ui.Render(ui.Body)
	})

	ui.Loop()
	if newView != nil {
		g.OpenView(newView)
		return false
	}
	if expand {
		g.ExpandView()
		return false
	}
	return true
}
