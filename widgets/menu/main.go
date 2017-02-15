package menu

import (
	"sort"

	ui "github.com/gizak/termui"
)

type Padding [2]int // x,y padding

type Menu struct {
	ui.Block
	SortItems   bool // enable automatic sorting of menu items
	TextFgColor ui.Attribute
	TextBgColor ui.Attribute
	Selectable  bool
	CursorPos   int
	items       Items
	padding     Padding
}

func NewMenu() *Menu {
	return &Menu{
		Block:       *ui.NewBlock(),
		TextFgColor: ui.ThemeAttr("par.text.fg"),
		TextBgColor: ui.ThemeAttr("par.text.bg"),
		CursorPos:   0,
		padding:     Padding{4, 2},
	}
}

func (m *Menu) AddItems(items ...Item) {
	for _, i := range items {
		m.items = append(m.items, i)
	}
	m.refresh()
}

// Remove menu item by value or label
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

// Sort menu items(if enabled) and re-calculate window size
func (m *Menu) refresh() {
	if m.SortItems {
		sort.Sort(m.items)
	}
	m.calcSize()
	ui.Render(m)
}

func (m *Menu) SelectedItem() Item {
	return m.items[m.CursorPos]
}

func (m *Menu) Buffer() ui.Buffer {
	var cell ui.Cell
	buf := m.Block.Buffer()

	for n, item := range m.items {
		x := m.padding[0]
		for _, ch := range item.Text() {
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
	if m.CursorPos < (len(m.items) - 1) {
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

// Set width and height based on menu items
func (m *Menu) calcSize() {
	m.Width = 8 // minimum width

	items := m.items
	for _, i := range m.items {
		s := i.Text()
		if len(s) > m.Width {
			m.Width = len(s)
		}
	}

	m.Width += (m.padding[0] * 2)
	m.Height = len(items) + (m.padding[1] * 2)
}
