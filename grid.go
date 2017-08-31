package main

import (
	"github.com/bcicen/ctop/config"
	"github.com/bcicen/ctop/cwidgets/single"
	ui "github.com/gizak/termui"
	"github.com/bcicen/ctop/entity"
)

func RedrawRows(clr bool) {
	// reinit body rows
	cGrid.Clear()

	// build layout
	y := 1
	if config.GetSwitchVal("enableHeader") {
		header.SetCount(cursor.LenContainers(), cursor.LenNodes(), cursor.LenServices(), cursor.LenTasks())
		header.SetFilter(config.GetVal("filterStr"))
		y += header.Height()
	}
	cGrid.SetY(y)
	if config.GetSwitchVal("swarmMode") {
		//for _, n := range cursor.filteredNodes {
			//cGrid.AddRows(n.Widgets)
		//}
		for _, s := range cursor.filteredServices {
			cGrid.AddRows(s.Widgets)
		}
		for _, t := range cursor.filteredTasks {
			cGrid.AddRows(t.Widgets)
		}
	} else {
		for _, c := range cursor.filteredContainers {
			cGrid.AddRows(c.Widgets)
		}
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

func SingleView(c *entity.Container) {
	ui.Clear()
	ui.DefaultEvtStream.ResetHandlers()
	defer ui.DefaultEvtStream.ResetHandlers()

	ex := single.NewSingle(c.Id)
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
		if config.GetSwitchVal("swarmMode") {
			needsClear := cursor.RefreshNodes()
			RedrawRows(needsClear)
			needsClear = cursor.RefreshServices()
			RedrawRows(needsClear)
			needsClear = cursor.RefreshTasks()
			RedrawRows(needsClear)
		} else {
			needsClear := cursor.RefreshContainers()
			RedrawRows(needsClear)
		}
	}
}

func Display() bool {
	var menu func()
	var single bool

	cGrid.SetWidth(ui.TermWidth())
	ui.DefaultEvtStream.Hook(logEvent)

	// initial draw
	header.Align()
	if config.GetSwitchVal("swarmMode") {
		cursor.RefreshNodes()
		cursor.RefreshServices()
		cursor.RefreshTasks()
	} else {
		cursor.RefreshContainers()
	}

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
		single = true
		ui.StopLoop()
	})
	ui.Handle("/sys/kbd/a", func(ui.Event) {
		config.Toggle("allContainers")
		RefreshDisplay()
	})
	ui.Handle("/sys/kbd/D", func(ui.Event) {
		e, _ := cursor.Selected()
		dumpContainer(e)
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
	if single {
		c, t := cursor.Selected()
		switch t {
		case "container":
			if c != nil {
				SingleView(cursor.SelectedContainer())
			}
		}
		return false
	}
	return true
}
