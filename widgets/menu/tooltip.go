package menu

import (
	ui "github.com/gizak/termui"
)

type ToolTip struct {
	ui.Block
	Lines       []string
	TextFgColor ui.Attribute
	TextBgColor ui.Attribute
	padding     Padding
}

func NewToolTip(lines ...string) *ToolTip {
	t := &ToolTip{
		Block:       *ui.NewBlock(),
		Lines:       lines,
		TextFgColor: ui.ThemeAttr("menu.text.fg"),
		TextBgColor: ui.ThemeAttr("menu.text.bg"),
		padding:     Padding{2, 1},
	}
	t.BorderFg = ui.ThemeAttr("menu.border.fg")
	t.BorderLabelFg = ui.ThemeAttr("menu.label.fg")
	t.X = 1
	t.Align()
	return t
}

func (t *ToolTip) Buffer() ui.Buffer {
	var cell ui.Cell
	buf := t.Block.Buffer()

	y := t.Y + t.padding[1]

	for n, line := range t.Lines {
		x := t.X + t.padding[0]
		for _, ch := range line {
			cell = ui.Cell{Ch: ch, Fg: t.TextFgColor, Bg: t.TextBgColor}
			buf.Set(x, y+n, cell)
			x++
		}
	}

	return buf
}

// Set width and height based on screen size
func (t *ToolTip) Align() {
	t.Width = ui.TermWidth() - (t.padding[0] * 2)
	t.Height = len(t.Lines) + (t.padding[1] * 2)
	t.Y = ui.TermHeight() - t.Height

	t.Block.Align()
}
