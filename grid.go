package main

import (
	"sort"

	ui "github.com/gizak/termui"
)

type Grid struct {
	containers map[string]*Container
}

func (g *Grid) AddContainer(id string) {
	g.containers[id] = NewContainer(id)
}

// Return sorted list of active container IDs
func (g *Grid) CIDs() []string {
	var ids []string
	for id, _ := range g.containers {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	return ids
}

func (g *Grid) Rows() (rows []*ui.Row) {
	for _, cid := range g.CIDs() {
		c := g.containers[cid]
		rows = append(rows, ui.NewRow(
			ui.NewCol(1, 0, c.widgets.cid),
			ui.NewCol(2, 0, c.widgets.cpu),
			ui.NewCol(2, 0, c.widgets.memory),
			ui.NewCol(2, 0, c.widgets.net),
		))
	}
	return rows
}

func header() *ui.Row {
	// cid
	c1 := ui.NewPar(" CID")
	c1.Border = false
	c1.Height = 2
	c1.Width = 20
	c1.TextFgColor = ui.ColorWhite

	// cpu
	c2 := ui.NewPar(" CPU")
	c2.Border = false
	c2.Height = 2
	c2.Width = 10
	c2.TextFgColor = ui.ColorWhite

	// mem
	c3 := ui.NewPar(" MEM")
	c3.Border = false
	c3.Height = 2
	c3.Width = 10
	c3.TextFgColor = ui.ColorWhite

	// net
	c4 := ui.NewPar(" NET RX/TX")
	c4.Border = false
	c4.Height = 2
	c4.Width = 10
	c4.TextFgColor = ui.ColorWhite
	return ui.NewRow(
		ui.NewCol(1, 0, c1),
		ui.NewCol(2, 0, c2),
		ui.NewCol(2, 0, c3),
		ui.NewCol(2, 0, c4),
	)
}

func Display(g *Grid) {
	if err := ui.Init(); err != nil {
		panic(err)
	}
	defer ui.Close()

	// build layout
	ui.Body.AddRows(header())

	for _, row := range g.Rows() {
		ui.Body.AddRows(row)
	}

	// calculate layout
	ui.Body.Align()

	ui.Render(ui.Body)

	ui.Handle("/sys/kbd/q", func(ui.Event) {
		ui.StopLoop()
	})
	ui.Handle("/timer/1s", func(e ui.Event) {
		//		t := e.Data.(ui.EvtTimer)
		//		i := t.Count
		//		if i > 103 {
		//			ui.StopLoop()
		//			return
		//		}

		ui.Render(ui.Body)
	})

	ui.Handle("/sys/wnd/resize", func(e ui.Event) {
		ui.Body.Width = ui.TermWidth()
		ui.Body.Align()
		ui.Clear()
		ui.Render(ui.Body)
	})

	ui.Loop()
}
