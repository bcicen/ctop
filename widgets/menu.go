package widgets

import (
	ui "github.com/gizak/termui"
)

type Padding [2]int // x,y padding

type Menu struct {
	ui.Block
	Items        []string
	DisplayItems []string
	TextFgColor  ui.Attribute
	TextBgColor  ui.Attribute
	Selectable   bool
	CursorPos    int
	padding      Padding
}

func NewMenu(items []string) *Menu {
	m := &Menu{
		Block:        *ui.NewBlock(),
		Items:        items,
		DisplayItems: []string{},
		TextFgColor:  ui.ThemeAttr("par.text.fg"),
		TextBgColor:  ui.ThemeAttr("par.text.bg"),
		Selectable:   false,
		CursorPos:    0,
		padding:      Padding{4, 2},
	}
	m.calcSize()
	return m
}

func (m *Menu) Buffer() ui.Buffer {
	var cell ui.Cell
	buf := m.Block.Buffer()

	for n, item := range m.displayItems() {
		x := m.padding[0]
		for _, ch := range item {
			// invert bg/fg colors on currently selected row
			if m.Selectable && n == m.CursorPos {
				cell = ui.Cell{Ch: ch, Fg: m.TextBgColor, Bg: m.TextFgColor}
			} else {
				cell = ui.Cell{Ch: ch, Fg: m.TextFgColor, Bg: m.TextBgColor}
			}
			buf.Set(x, n+m.padding[1], cell)
			x++
		}
	}

	return buf
}

func (m *Menu) Up(ui.Event) {
	if m.CursorPos > 0 {
		m.CursorPos--
		ui.Render(m)
	}
}

func (m *Menu) Down(ui.Event) {
	if m.CursorPos < (len(m.Items) - 1) {
		m.CursorPos++
		ui.Render(m)
	}
}

// Setup some default handlers for menu navigation
func (m *Menu) NavigationHandlers() {
	ui.Handle("/sys/kbd/<up>", m.Up)
	ui.Handle("/sys/kbd/<down>", m.Down)
	ui.Handle("/sys/kbd/q", func(ui.Event) { ui.StopLoop() })
}

// override display of items, if given
func (m *Menu) displayItems() []string {
	if len(m.DisplayItems) == len(m.Items) {
		return m.DisplayItems
	}
	return m.Items
}

// Set width and height based on menu items
func (m *Menu) calcSize() {
	m.Width = 8 // minimum width

	items := m.displayItems()
	for _, s := range items {
		if len(s) > m.Width {
			m.Width = len(s)
		}
	}

	m.Width += (m.padding[0] * 2)
	m.Height = len(items) + (m.padding[1] * 2)
}
