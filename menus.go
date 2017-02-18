package main

import (
	"github.com/bcicen/ctop/config"
	"github.com/bcicen/ctop/widgets"
	"github.com/bcicen/ctop/widgets/menu"
	ui "github.com/gizak/termui"
)

var helpDialog = []menu.Item{
	menu.Item{"[h] - open this help dialog", ""},
	menu.Item{"[s] - select container sort field", ""},
	menu.Item{"[r] - reverse container sort order", ""},
	menu.Item{"[q] - exit ctop", ""},
}

func HelpMenu() {
	ui.DefaultEvtStream.ResetHandlers()
	defer ui.DefaultEvtStream.ResetHandlers()

	m := menu.NewMenu()
	m.TextFgColor = ui.ColorWhite
	m.BorderLabel = "Help"
	m.BorderFg = ui.ColorCyan
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
	ui.DefaultEvtStream.ResetHandlers()
	defer ui.DefaultEvtStream.ResetHandlers()

	m := menu.NewMenu()
	m.Selectable = true
	m.SortItems = true
	m.TextFgColor = ui.ColorWhite
	m.BorderLabel = "Sort Field"
	m.BorderFg = ui.ColorCyan

	for _, field := range SortFields() {
		m.AddItems(menu.Item{field, ""})
	}

	// set cursor position to current sort field
	//current := config.GetVal("sortField")
	//for n, item := range m.Items {
	//if item.Val == current {
	//m.CursorPos = n
	//}
	//}

	ui.Render(m)
	m.NavigationHandlers()
	ui.Handle("/sys/kbd/<enter>", func(ui.Event) {
		config.Update("sortField", m.SelectedItem().Val)
		ui.StopLoop()
	})
	ui.Loop()
}
