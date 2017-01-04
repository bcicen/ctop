package main

import (
	ui "github.com/gizak/termui"
)

var (
	padding  = 2
	minWidth = 30
)

type Menu struct {
	ui.Block
	Items       []string
	TextFgColor ui.Attribute
	TextBgColor ui.Attribute
	Selectable  bool
	cursorPos   int
}

func NewMenu(items []string) *Menu {
	m := &Menu{
		Block:       *ui.NewBlock(),
		Items:       items,
		TextFgColor: ui.ThemeAttr("par.text.fg"),
		TextBgColor: ui.ThemeAttr("par.text.bg"),
		Selectable:  false,
		cursorPos:   0,
	}
	m.Width = calcWidth(items)
	return m
}

// return dynamic width based on string len in items
func calcWidth(items []string) int {
	maxWidth := 0
	for _, s := range items {
		if len(s) > maxWidth {
			maxWidth = len(s)
		}
	}
	if maxWidth < minWidth {
		maxWidth = minWidth
	}
	return maxWidth + (padding * 2)
}

func (m *Menu) Buffer() ui.Buffer {
	var cell ui.Cell

	buf := m.Block.Buffer()

	for n, item := range m.Items {
		x := 2 // initial offset
		for _, ch := range item {
			// invert bg/fg colors on currently selected row
			if m.Selectable && n == m.cursorPos {
				cell = ui.Cell{Ch: ch, Fg: m.TextBgColor, Bg: m.TextFgColor}
			} else {
				cell = ui.Cell{Ch: ch, Fg: m.TextFgColor, Bg: m.TextBgColor}
			}
			buf.Set(x, n+2, cell)
			x++
		}
	}

	return buf
}

func (m *Menu) Up() {
	if m.cursorPos > 0 {
		m.cursorPos--
		ui.Render(m)
	}
}

func (m *Menu) Down() {
	if m.cursorPos < (len(m.Items) - 1) {
		m.cursorPos++
		ui.Render(m)
	}
}

func HelpMenu(g *Grid) {
	m := NewMenu([]string{
		"[h] - open this help dialog",
		"[q] - exit ctop",
	})
	m.Height = 10
	m.TextFgColor = ui.ColorWhite
	m.BorderLabel = "Help"
	m.BorderFg = ui.ColorCyan
	ui.Render(m)
	ui.Handle("/sys/kbd/", func(ui.Event) {
		ui.StopLoop()
	})
	ui.Loop()
}

func SortMenu(g *Grid) {
	m := NewMenu(SortFields)
	m.Height = 10
	m.Selectable = true
	m.TextFgColor = ui.ColorWhite
	m.BorderLabel = "Sort Field"
	m.BorderFg = ui.ColorCyan
	ui.Render(m)
	ui.Handle("/sys/kbd/<up>", func(ui.Event) {
		m.Up()
	})
	ui.Handle("/sys/kbd/<down>", func(ui.Event) {
		m.Down()
	})
	ui.Handle("/sys/kbd/<enter>", func(ui.Event) {
		g.containerMap.config.sortField = m.Items[m.cursorPos]
		ui.StopLoop()
	})
	ui.Loop()
}
