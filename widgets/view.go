package widgets

import (
	ui "github.com/gizak/termui"
)

type TextView struct {
	ui.Block
	inputStream <-chan string
	render      chan bool
	Text        []string // all the text
	TextOut     []string // text to be displayed
	TextFgColor ui.Attribute
	TextBgColor ui.Attribute
	padding     Padding
}

func NewTextView(lines <-chan string) *TextView {
	t := &TextView{
		Block:       *ui.NewBlock(),
		inputStream: lines,
		render:      make(chan bool),
		Text:        []string{},
		TextOut:     []string{},
		TextFgColor: ui.ThemeAttr("menu.text.fg"),
		TextBgColor: ui.ThemeAttr("menu.text.bg"),
		padding:     Padding{4, 2},
	}

	t.BorderFg = ui.ThemeAttr("menu.border.fg")
	t.BorderLabelFg = ui.ThemeAttr("menu.label.fg")
	t.Height = ui.TermHeight()
	t.Width = ui.TermWidth()

	t.readInputLoop()
	t.renderLoop()
	return t
}

func (t *TextView) Resize() {
	ui.Clear()
	t.Height = ui.TermHeight()
	t.Width = ui.TermWidth()
	t.render <- true
}

func (t *TextView) Buffer() ui.Buffer {
	var cell ui.Cell
	buf := t.Block.Buffer()

	x := t.Block.X + t.padding[0]
	y := t.Block.Y + t.padding[1]

	for _, line := range t.TextOut {
		for _, ch := range line {
			cell = ui.Cell{Ch: ch, Fg: t.TextFgColor, Bg: t.TextBgColor}
			buf.Set(x, y, cell)
			x++
		}
		x = t.Block.X + t.padding[0]
		y++
	}
	return buf
}

func (t *TextView) renderLoop() {
	go func() {
		for range t.render {
			maxWidth := t.Width - (t.padding[0] * 2)
			height := t.Height - (t.padding[1] * 2)
			t.TextOut = []string{}
			for i := len(t.Text) - 1; i >= 0; i-- {
				lines := splitLine(t.Text[i], maxWidth)
				t.TextOut = append(lines, t.TextOut...)
				if len(t.TextOut) > height {
					t.TextOut = t.TextOut[:height]
					break
				}
			}
			ui.Render(t)
		}
	}()
}

func (t *TextView) readInputLoop() {
	go func() {
		for line := range t.inputStream {
			t.Text = append(t.Text, line)
			t.render <- true
		}
		close(t.render)
	}()
}

func splitLine(line string, lineSize int) []string {
	var lines []string
	for {
		if len(line) < lineSize {
			lines = append(lines, line)
			return lines
		}
		lines = append(lines, line[:lineSize])
		line = line[lineSize:]
	}
	return lines
}
