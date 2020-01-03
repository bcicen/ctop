package menu

import (
	"sort"

	ui "github.com/gizak/termui"
)

type Padding [2]int // x,y padding

type Menu struct {
	ui.Block
	SortItems   bool   // enable automatic sorting of menu items
	Selectable  bool   // whether menu is navigable
	SubText     string // optional text to display before items
	TextFgColor ui.Attribute
	TextBgColor ui.Attribute
	cursorPos   int
	items       Items
	padding     Padding
	toolTip     *ToolTip
}

func NewMenu() *Menu {
	m := &Menu{
		Block:       *ui.NewBlock(),
		TextFgColor: ui.ThemeAttr("menu.text.fg"),
		TextBgColor: ui.ThemeAttr("menu.text.bg"),
		cursorPos:   0,
		padding:     Padding{4, 2},
	}
	m.BorderFg = ui.ThemeAttr("menu.border.fg")
	m.BorderLabelFg = ui.ThemeAttr("menu.label.fg")
	m.X = 1
	return m
}

// Append Item to Menu
func (m *Menu) AddItems(items ...Item) {
	for _, i := range items {
		m.items = append(m.items, i)
	}
	m.refresh()
}

// DelItem removes menu item by value or label
func (m *Menu) DelItem(s string) (success bool) {
	for n, i := range m.items {
		if i.Val == s || i.Label == s {
			m.items = append(m.items[:n], m.items[n+1:]...)
			success = true
			m.refresh()
			break
		}
	}
	return success
}

// ClearItems removes all current menu items
func (m *Menu) ClearItems() {
	m.items = m.items[:0]
}

// Move cursor to an position by Item value or label
func (m *Menu) SetCursor(s string) (success bool) {
	for n, i := range m.items {
		if i.Val == s || i.Label == s {
			m.cursorPos = n
			return true
		}
	}
	return false
}

// SetToolTip sets an optional tooltip string to show at bottom of screen
func (m *Menu) SetToolTip(lines ...string) {
	m.toolTip = NewToolTip(lines...)
}

func (m *Menu) SelectedItem() Item {
	return m.items[m.cursorPos]
}

func (m *Menu) SelectedValue() string {
	return m.items[m.cursorPos].Val
}

func (m *Menu) Buffer() ui.Buffer {
	var cell ui.Cell
	buf := m.Block.Buffer()

	y := m.Y + m.padding[1]

	if m.SubText != "" {
		x := m.X + m.padding[0]
		for i, ch := range m.SubText {
			cell = ui.Cell{Ch: ch, Fg: m.TextFgColor, Bg: m.TextBgColor}
			buf.Set(x+i, y, cell)
		}
		y += 2
	}

	for n, item := range m.items {
		x := m.X + m.padding[0]
		for _, ch := range item.Text() {
			// invert bg/fg colors on currently selected row
			if m.Selectable && n == m.cursorPos {
				cell = ui.Cell{Ch: ch, Fg: ui.ColorBlack, Bg: m.TextFgColor}
			} else {
				cell = ui.Cell{Ch: ch, Fg: m.TextFgColor, Bg: m.TextBgColor}
			}
			buf.Set(x, y+n, cell)
			x++
		}
	}

	if m.toolTip != nil {
		buf.Merge(m.toolTip.Buffer())
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
	if m.cursorPos < (len(m.items) - 1) {
		m.cursorPos++
		ui.Render(m)
	}
}

// Sort menu items(if enabled) and re-calculate window size
func (m *Menu) refresh() {
	if m.SortItems {
		sort.Sort(m.items)
	}
	m.calcSize()
	ui.Render(m)
}

// Set width and height based on menu items
func (m *Menu) calcSize() {
	m.Width = 7 // minimum width

	var height int
	for _, i := range m.items {
		s := i.Text()
		if len(s) > m.Width {
			m.Width = len(s)
		}
		height++
	}

	if m.SubText != "" {
		if len(m.SubText) > m.Width {
			m.Width = len(m.SubText)
		}
		height += 2
	}

	m.Width += (m.padding[0] * 2)
	m.Height = height + (m.padding[1] * 2)
}
