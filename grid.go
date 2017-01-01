package main

import (
	"fmt"

	"github.com/bcicen/ctop/views"
	ui "github.com/gizak/termui"
)

type Grid struct {
	cursorPos  uint
	containers *ContainerMap
}

func NewGrid() *Grid {
	return &Grid{
		cursorPos:  0,
		containers: NewContainerMap(),
	}
}

// Return sorted list of container IDs
func (g *Grid) CIDs() []string {
	var ids []string
	for _, c := range g.containers.Sorted() {
		ids = append(ids, c.id)
	}
	return ids
}

// Redraw the cursor with the currently selected row
func (g *Grid) Cursor() {
	for n, c := range g.containers.Sorted() {
		if uint(n) == g.cursorPos {
			c.widgets.name.TextFgColor = ui.ColorDefault
			c.widgets.name.TextBgColor = ui.ColorWhite
		} else {
			c.widgets.name.TextFgColor = ui.ColorWhite
			c.widgets.name.TextBgColor = ui.ColorDefault
		}
	}
	ui.Render(ui.Body)
}

func (g *Grid) Redraw() {
	// reinit body rows
	ui.Body.Rows = []*ui.Row{}
	// build layout
	ui.Body.AddRows(header())

	for _, c := range g.containers.Sorted() {
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

	if err := ui.Init(); err != nil {
		panic(err)
	}
	defer ui.Close()
	// calculate layout
	ui.Body.Align()
	g.Cursor()
	ui.Render(ui.Body)

	ui.Handle("/sys/kbd/<up>", func(ui.Event) {
		if g.cursorPos > 0 {
			g.cursorPos -= 1
			g.Cursor()
		}
	})
	ui.Handle("/sys/kbd/<down>", func(ui.Event) {
		if g.cursorPos < (g.containers.Len() - 1) {
			g.cursorPos += 1
			g.Cursor()
		}
	})
	ui.Handle("/sys/kbd/h", func(ui.Event) {
		newView = views.Help
		ui.StopLoop()
	})
	ui.Handle("/sys/kbd/q", func(ui.Event) {
		ui.StopLoop()
	})
	ui.Handle("/timer/1s", func(e ui.Event) {
		g.Redraw()
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
