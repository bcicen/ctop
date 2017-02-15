package widgets

import (
	"sort"

	ui "github.com/gizak/termui"
)

type Padding [2]int // x,y padding

type MenuItem struct {
	Val   string
	Label string
}

// Use label as display text of item, if given
func (m MenuItem) Text() string {
	if m.Label != "" {
		return m.Label
	}
	return m.Val
}

type MenuItems []MenuItem

// Sort methods for MenuItems
func (m MenuItems) Len() int      { return len(m) }
func (m MenuItems) Swap(a, b int) { m[a], m[b] = m[b], m[a] }
func (m MenuItems) Less(a, b int) bool {
	return m[a].Text() < m[b].Text()
}

type Menu struct {
	ui.Block
	Items       MenuItems
	TextFgColor ui.Attribute
	TextBgColor ui.Attribute
	Selectable  bool
	CursorPos   int
	padding     Padding
}

func NewMenu(items []string) *Menu {
	m := &Menu{
		Block:       *ui.NewBlock(),
		TextFgColor: ui.ThemeAttr("par.text.fg"),
		TextBgColor: ui.ThemeAttr("par.text.bg"),
		Selectable:  false,
		CursorPos:   0,
		padding:     Padding{4, 2},
	}

	for _, s := range items {
		m.Items = append(m.Items, MenuItem{Val: s})
	}
	sort.Sort(m.Items)

	m.calcSize()
	return m
}

func (m *Menu) SetItems(items []MenuItem) {
	m.Items = items
	sort.Sort(m.Items)
	m.calcSize()
}

func (m *Menu) SelectedItem() MenuItem {
	return m.Items[m.CursorPos]
}

func (m *Menu) Buffer() ui.Buffer {
	var cell ui.Cell
	buf := m.Block.Buffer()

	for n, item := range m.Items {
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

// Set width and height based on menu items
func (m *Menu) calcSize() {
	m.Width = 8 // minimum width

	items := m.Items
	for _, i := range m.Items {
		s := i.Text()
		if len(s) > m.Width {
			m.Width = len(s)
		}
	}

	m.Width += (m.padding[0] * 2)
	m.Height = len(items) + (m.padding[1] * 2)
}
