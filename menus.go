package main

import (
	"github.com/bcicen/ctop/config"
	"github.com/bcicen/ctop/widgets"
	ui "github.com/gizak/termui"
)

var helpDialog = []string{
	"[h] - open this help dialog",
	"[s] - select container sort field",
	"[r] - reverse container sort order",
	"[q] - exit ctop",
}

func HelpMenu() {
	ResetView()
	defer ResetView()

	m := widgets.NewMenu()
	m.TextFgColor = ui.ColorWhite
	m.BorderLabel = "Help"
	m.BorderFg = ui.ColorCyan
	m.AddItems(widgets.NewMenuItems(helpDialog)...)
	ui.Render(m)
	ui.Handle("/sys/kbd/", func(ui.Event) {
		ui.StopLoop()
	})
	ui.Loop()
}

func FilterMenu() {
	ui.DefaultEvtStream.ResetHandlers()
	defer ResetView()

	i := widgets.NewInput()
	i.TextFgColor = ui.ColorWhite
	i.BorderLabel = "Filter"
	i.BorderFg = ui.ColorCyan
	i.SetY(ui.TermHeight() - i.Height)
	ui.Render(i)
	i.InputHandlers()
	ui.Handle("/sys/kbd/<enter>", func(ui.Event) {
		config.Update("filterStr", i.Data)
		ui.StopLoop()
	})
	ui.Loop()
}

func SortMenu() {
	ResetView()
	defer ResetView()

	m := widgets.NewMenu()
	m.Selectable = true
	m.SortItems = true
	m.TextFgColor = ui.ColorWhite
	m.BorderLabel = "Sort Field"
	m.BorderFg = ui.ColorCyan

	m.AddItems(widgets.NewMenuItems(SortFields())...)

	// set cursor position to current sort field
	current := config.Get("sortField")
	for n, item := range m.Items {
		if item.Val == current {
			m.CursorPos = n
		}
	}

	ui.Render(m)
	m.NavigationHandlers()
	ui.Handle("/sys/kbd/<enter>", func(ui.Event) {
		config.Update("sortField", m.Items[m.CursorPos].Val)
		ui.StopLoop()
	})
	ui.Loop()
}
