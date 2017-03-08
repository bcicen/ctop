package main

import (
	"github.com/bcicen/ctop/config"
	"github.com/bcicen/ctop/cwidgets/expanded"
	ui "github.com/gizak/termui"
)

func maxRows() int {
	return ui.TermHeight() - 2 - cGrid.Y
}

func RedrawRows(clr bool) {
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

	if clr {
		ui.Clear()
		log.Debugf("screen cleared")
	}
	if config.GetSwitchVal("enableHeader") {
		header.Render()
	}
	cGrid.Align()
	ui.Render(cGrid)
}

func ExpandView(c *Container) {
	ui.Clear()
	ui.DefaultEvtStream.ResetHandlers()
	defer ui.DefaultEvtStream.ResetHandlers()

	ex := expanded.NewExpanded(c.Id)
	c.SetUpdater(ex)

	ex.Align()
	ui.Render(ex)
	ui.Handle("/timer/1s", func(ui.Event) {
		ui.Render(ex)
	})
	ui.Handle("/sys/wnd/resize", func(e ui.Event) {
		ex.SetWidth(ui.TermWidth())
		ex.Align()
		log.Infof("resize: width=%v max-rows=%v", ex.Width, maxRows())
	})
	ui.Handle("/sys/kbd/", func(ui.Event) {
		ui.StopLoop()
	})
	ui.Loop()

	c.SetUpdater(c.Widgets)
}

func RefreshDisplay() {
	needsClear := cursor.RefreshContainers()
	RedrawRows(needsClear)
}

func Display() bool {
	var menu func()
	var expand bool

	cGrid.SetWidth(ui.TermWidth())
	ui.DefaultEvtStream.Hook(logEvent)

	// initial draw
	header.Align()
	cursor.RefreshContainers()
	RedrawRows(true)

	ui.Handle("/sys/kbd/<up>", func(ui.Event) {
		cursor.Up()
	})
	ui.Handle("/sys/kbd/<down>", func(ui.Event) {
		cursor.Down()
	})
	ui.Handle("/sys/kbd/<enter>", func(ui.Event) {
		expand = true
		ui.StopLoop()
	})

	ui.Handle("/sys/kbd/a", func(ui.Event) {
		config.Toggle("allContainers")
		RefreshDisplay()
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
		RedrawRows(true)
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
		RefreshDisplay()
	})

	ui.Handle("/sys/wnd/resize", func(e ui.Event) {
		header.Align()
		cGrid.SetWidth(ui.TermWidth())
		log.Infof("resize: width=%v max-rows=%v", cGrid.Width, maxRows())
		RedrawRows(true)
	})

	ui.Loop()
	if menu != nil {
		menu()
		return false
	}
	if expand {
		ExpandView(cursor.Selected())
		return false
	}
	return true
}
