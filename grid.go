package main

import (
	"github.com/bcicen/ctop/config"
	"github.com/bcicen/ctop/container"
	"github.com/bcicen/ctop/cwidgets/expanded"
	ui "github.com/gizak/termui"
)

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

	for _, c := range cursor.filtered {
		cGrid.AddRows(c.Widgets)
	}

	if clr {
		ui.Clear()
		log.Debugf("screen cleared")
	}
	if config.GetSwitchVal("enableHeader") {
		ui.Render(header)
	}
	cGrid.Align()
	ui.Render(cGrid)
}

func ExpandView(c *container.Container) {
	ui.Clear()
	ui.DefaultEvtStream.ResetHandlers()
	defer ui.DefaultEvtStream.ResetHandlers()

	ex := expanded.NewExpanded(c.Id)
	c.SetUpdater(ex)

	ex.Align()
	ui.Render(ex)

	HandleKeys("up", ex.Up)
	HandleKeys("down", ex.Down)
	ui.Handle("/sys/kbd/", func(ui.Event) { ui.StopLoop() })

	ui.Handle("/timer/1s", func(ui.Event) { ui.Render(ex) })
	ui.Handle("/sys/wnd/resize", func(e ui.Event) {
		ex.SetWidth(ui.TermWidth())
		ex.Align()
		log.Infof("resize: width=%v max-rows=%v", ex.Width, cGrid.MaxRows())
	})

	ui.Loop()
	c.SetUpdater(c.Widgets)
}

func RefreshDisplay() {
	// skip display refresh during scroll
	if !cursor.isScrolling {
		needsClear := cursor.RefreshContainers()
		RedrawRows(needsClear)
	}
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

	HandleKeys("up", cursor.Up)
	HandleKeys("down", cursor.Down)

	HandleKeys("pgup", cursor.PgUp)
	HandleKeys("pgdown", cursor.PgDown)

	HandleKeys("exit", ui.StopLoop)
	HandleKeys("help", func() {
		menu = HelpMenu
		ui.StopLoop()
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
	ui.Handle("/sys/kbd/H", func(ui.Event) {
		config.Toggle("enableHeader")
		RedrawRows(true)
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
		cursor.ScrollPage()
		cGrid.SetWidth(ui.TermWidth())
		log.Infof("resize: width=%v max-rows=%v", cGrid.Width, cGrid.MaxRows())
		RedrawRows(true)
	})

	ui.Loop()
	if menu != nil {
		menu()
		return false
	}
	if expand {
		c := cursor.Selected()
		if c != nil {
			ExpandView(c)
		}
		return false
	}
	return true
}
