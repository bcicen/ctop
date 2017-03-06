package main

import (
	"github.com/bcicen/ctop/config"
	ui "github.com/gizak/termui"
)

func maxRows() int {
	return ui.TermHeight() - 2 - cGrid.Y
}

func RedrawRows() {
	// reinit body rows
	cGrid.Clear()

	// build layout
	y := 1
	if config.GetSwitchVal("enableHeader") {
		header.SetCount(cursor.Len())
		header.SetFilter(config.GetVal("filterStr"))
		y += header.Height()
	}
	cGrid.SetY(y)

	var cursorVisible bool
	max := maxRows()
	for n, c := range cursor.containers {
		if n >= max {
			break
		}
		cGrid.AddRows(c.Widgets)
		if c.Id == cursor.selectedID {
			cursorVisible = true
		}
	}

	if !cursorVisible {
		cursor.Reset()
	}

	ui.Clear()
	if config.GetSwitchVal("enableHeader") {
		header.Render()
	}
	cGrid.Align()
	ui.Render(cGrid)
}

//func (g *Grid) ExpandView() {
//ui.Clear()
//ui.DefaultEvtStream.ResetHandlers()
//defer ui.DefaultEvtStream.ResetHandlers()

//container, _ := g.cSource.Get(g.cursorID)
//// copy current widgets to restore on exit view
//curWidgets := container.widgets
//container.Expand()

//ui.Render(container.widgets)
//ui.Handle("/timer/1s", func(ui.Event) {
//ui.Render(container.widgets)
//})
//ui.Handle("/sys/kbd/", func(ui.Event) {
//ui.StopLoop()
//})
//ui.Loop()

//container.widgets = curWidgets
//container.widgets.Reset()
//}

func Display() bool {
	var menu func()

	cGrid.SetWidth(ui.TermWidth())
	ui.DefaultEvtStream.Hook(logEvent)

	// initial draw
	header.Align()
	cursor.RefreshContainers()
	RedrawRows()

	ui.Handle("/sys/kbd/<up>", func(ui.Event) {
		cursor.Up()
	})
	ui.Handle("/sys/kbd/<down>", func(ui.Event) {
		cursor.Down()
	})
	ui.Handle("/sys/kbd/<enter>", func(ui.Event) {
		//c := g.containers[g.cursorIdx()]
		//c.Widgets.ToggleExpand()
		RedrawRows()
	})

	ui.Handle("/sys/kbd/a", func(ui.Event) {
		config.Toggle("allContainers")
		cursor.RefreshContainers()
		RedrawRows()
	})
	ui.Handle("/sys/kbd/D", func(ui.Event) {
		dumpContainer(cursor.Selected())
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
		RedrawRows()
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
		cursor.RefreshContainers()
		RedrawRows()
	})

	ui.Handle("/sys/wnd/resize", func(e ui.Event) {
		header.Align()
		cGrid.SetWidth(ui.TermWidth())
		log.Infof("resize: width=%v max-rows=%v", cGrid.Width, maxRows())
		RedrawRows()
	})

	ui.Loop()
	if menu != nil {
		ui.Clear()
		menu()
		return false
	}
	return true
}
