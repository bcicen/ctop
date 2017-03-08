package main

import (
	"github.com/bcicen/ctop/config"
	"github.com/bcicen/ctop/widgets"
	"github.com/bcicen/ctop/widgets/menu"
	ui "github.com/gizak/termui"
)

var helpDialog = []menu.Item{
	menu.Item{"[a] - toggle display of all containers", ""},
	menu.Item{"[f] - filter displayed containers", ""},
	menu.Item{"[h] - open this help dialog", ""},
	menu.Item{"[H] - toggle cTop header", ""},
	menu.Item{"[s] - select container sort field", ""},
	menu.Item{"[r] - reverse container sort order", ""},
	menu.Item{"[q] - exit ctop", ""},
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

	for _, field := range SortFields() {
		m.AddItems(menu.Item{field, ""})
	}

	// set cursor position to current sort field
	m.SetCursor(config.GetVal("sortField"))

	ui.Render(m)
	m.NavigationHandlers()
	ui.Handle("/sys/kbd/<enter>", func(ui.Event) {
		config.Update("sortField", m.SelectedItem().Val)
		ui.StopLoop()
	})
	ui.Loop()
}
