package main

import (
	"github.com/bcicen/ctop/config"
	"github.com/bcicen/ctop/entity"
	"github.com/bcicen/ctop/widgets"
	"github.com/bcicen/ctop/widgets/menu"
	ui "github.com/gizak/termui"
)

var helpDialog = []menu.Item{
	menu.Item{"[a] - toggle display of all containers", ""},
	menu.Item{"[f] - filter displayed containers", ""},
	menu.Item{"[h] - open this help dialog", ""},
	menu.Item{"[H] - toggle ctop header", ""},
	menu.Item{"[s] - select container sort field", ""},
	menu.Item{"[w] - toggle swarm mode", ""},
	menu.Item{"[r] - reverse container sort order", ""},
	menu.Item{"[q] - exit ctop", ""},
	menu.Item{"[m] - open menu container", ""},
}

func HelpMenu() {
	ui.Clear()
	ui.DefaultEvtStream.ResetHandlers()
	defer ui.DefaultEvtStream.ResetHandlers()

	m := menu.NewMenu()
	m.BorderLabel = "Help"
	m.AddItems(helpDialog...)
	ui.Render(m)
	ui.Handle("/sys/kbd/", func(ui.Event) {
		ui.StopLoop()
	})
	ui.Loop()
}

func FilterMenu() {
	ui.DefaultEvtStream.ResetHandlers()
	defer ui.DefaultEvtStream.ResetHandlers()

	i := widgets.NewInput()
	i.BorderLabel = "Filter"
	i.SetY(ui.TermHeight() - i.Height)
	i.Data = config.GetVal("filterStr")
	ui.Render(i)

	// refresh container rows on input
	stream := i.Stream()
	go func() {
		for s := range stream {
			config.Update("filterStr", s)
			RefreshDisplay()
			ui.Render(i)
		}
	}()

	i.InputHandlers()
	ui.Handle("/sys/kbd/<escape>", func(ui.Event) {
		config.Update("filterStr", "")
		ui.StopLoop()
	})
	ui.Handle("/sys/kbd/<enter>", func(ui.Event) {
		config.Update("filterStr", i.Data)
		ui.StopLoop()
	})
	ui.Loop()
}

func SortMenu() {
	ui.Clear()
	ui.DefaultEvtStream.ResetHandlers()
	defer ui.DefaultEvtStream.ResetHandlers()

	m := menu.NewMenu()
	m.Selectable = true
	m.SortItems = true
	m.BorderLabel = "Sort Field"

	for _, field := range entity.SortFields() {
		m.AddItems(menu.Item{field, ""})
	}

	// set cursor position to current sort field
	m.SetCursor(config.GetVal("sortField"))

	HandleKeys("up", m.Up)
	HandleKeys("down", m.Down)
	HandleKeys("exit", ui.StopLoop)

	ui.Handle("/sys/kbd/<enter>", func(ui.Event) {
		config.Update("sortField", m.SelectedItem().Val)
		ui.StopLoop()
	})

	ui.Render(m)
	ui.Loop()
}

func ContainerMenu() {
	if config.GetSwitchVal("swarMode") {
		return
	}

	c := cursor.Selected()
	if c == nil {
		return
	}

	ui.DefaultEvtStream.ResetHandlers()
	defer ui.DefaultEvtStream.ResetHandlers()

	m := menu.NewMenu()
	m.Selectable = true

	m.BorderLabel = "Menu"
	var items []menu.Item
	if c.GetMeta("state") == "running" {
		items = append(items, menu.Item{Val: "stop", Label: "stop"})
	}
	if c.GetMeta("state") == "exited" {
		items = append(items, menu.Item{Val: "start", Label: "start"})
		items = append(items, menu.Item{Val: "remove", Label: "remove"})
	}
	items = append(items, menu.Item{Val: "cancel", Label: "cancel"})

	m.AddItems(items...)
	ui.Render(m)

	HandleKeys("up", m.Up)
	HandleKeys("down", m.Down)
	ui.Handle("/sys/kbd/<enter>", func(ui.Event) {
		switch m.SelectedItem().Val {
		case "start":
			c.(*entity.Container).Start()
			ui.StopLoop()
		case "stop":
			c.(*entity.Container).Stop()
			ui.StopLoop()
		case "remove":
			c.(*entity.Container).Remove()
			ui.StopLoop()
		case "cancel":
			ui.StopLoop()
		}
	})
	ui.Handle("/sys/kbd/", func(ui.Event) {
		ui.StopLoop()
	})
	ui.Loop()
}

func LogMenu() {
	if config.GetSwitchVal("swarmMode") {
		return
	}

	c := cursor.Selected()
	if c == nil {
		return
	}

	ui.DefaultEvtStream.ResetHandlers()
	defer ui.DefaultEvtStream.ResetHandlers()

	logs, quit := logReader(c.(*entity.Container))
	m := widgets.NewTextView(logs)
	m.BorderLabel = "Logs"
	ui.Render(m)

	ui.Handle("/sys/wnd/resize", func(e ui.Event) {
		m.Resize()
		ui.Render(m)
	})
	ui.Handle("/sys/kbd/", func(ui.Event) {
		quit <- true
		ui.StopLoop()
	})
	ui.Loop()
}

func logReader(container *entity.Container) (logs chan string, quit chan bool) {

	logCollector := container.Logs()
	stream := logCollector.Stream()
	logs = make(chan string)
	quit = make(chan bool)

	go func() {
		for {
			select {
			case log := <-stream:
				logs <- log.Message
			case <-quit:
				logCollector.Stop()
				close(logs)
				return
			}
		}
	}()
	return
}
