package main

import (
	"fmt"

	"github.com/bcicen/ctop/views"
	ui "github.com/gizak/termui"
)

type Grid struct {
	cursorID     string // id of currently selected container
	containers   []*Container
	containerMap *ContainerMap
}

func NewGrid() *Grid {
	containerMap := NewContainerMap()
	containers := containerMap.Sorted()
	return &Grid{
		cursorID:     containers[0].id,
		containers:   containers,
		containerMap: containerMap,
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
			c.widgets.name.TextFgColor = ui.ColorDefault
			c.widgets.name.TextBgColor = ui.ColorWhite
		} else {
			c.widgets.name.TextFgColor = ui.ColorWhite
			c.widgets.name.TextBgColor = ui.ColorDefault
		}
		ui.Render(ui.Body)
	}
}

func (g *Grid) redrawRows() {
	// reinit body rows
	ui.Body.Rows = []*ui.Row{}

	// build layout
	ui.Body.AddRows(header())
	for _, c := range g.containers {
		ui.Body.AddRows(c.widgets.MakeRow())
	}

	ui.Body.Align()
	ui.Render(ui.Body)
}

func header() *ui.Row {
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

type View func()

func ResetView() {
	ui.DefaultEvtStream.ResetHandlers()
	ui.Clear()
}

func OpenView(v View) {
	ResetView()
	defer ResetView()
	v()
}

func Display(g *Grid) bool {
	var newView View

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
	ui.Handle("/sys/kbd/h", func(ui.Event) {
		newView = views.Help
		ui.StopLoop()
	})
	ui.Handle("/sys/kbd/q", func(ui.Event) {
		ui.StopLoop()
	})
	ui.Handle("/timer/1s", func(e ui.Event) {
		g.containers = g.containerMap.Sorted() // refresh containers for current sort order
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
		OpenView(newView)
		return false
	}
	return true
}
