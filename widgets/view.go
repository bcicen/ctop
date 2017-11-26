package widgets

import (
	ui "github.com/gizak/termui"
	"fmt"
)

type TextView struct {
	ui.Block
	inputStream <- chan string
	render      chan bool
	Text        []string // all the text
	TextOut     []string // text to be displayed
	TextFgColor ui.Attribute
	TextBgColor ui.Attribute
	padding     Padding
}

func NewTextView(lines <- chan string) *TextView {
	i := &TextView{
		Block:       *ui.NewBlock(),
		inputStream: lines,
		render:      make(chan bool),
		Text:        []string{},
		TextOut:     []string{},
		TextFgColor: ui.ThemeAttr("menu.text.fg"),
		TextBgColor: ui.ThemeAttr("menu.text.bg"),
		padding:     Padding{4, 2},
	}

	i.BorderFg = ui.ThemeAttr("menu.border.fg")
	i.BorderLabelFg = ui.ThemeAttr("menu.label.fg")

	ui.Clear()
	i.Height = ui.TermHeight()
	i.Width = ui.TermWidth()

	i.readInputLoop()
	i.renderLoop()
	return i
}

func (i *TextView) Buffer() ui.Buffer {

	var cell ui.Cell
	buf := i.Block.Buffer()

	x := i.Block.X + i.padding[0]
	y := i.Block.Y + i.padding[1]

	for _, line := range i.TextOut {
		for _, ch := range line {
			cell = ui.Cell{Ch: ch, Fg: i.TextFgColor, Bg: i.TextBgColor}
			buf.Set(x, y, cell)
			x++
		}
		x = i.Block.X + i.padding[0]
		y++
	}
	return buf
}

func (i *TextView) renderLoop() {
	go func() {
		for range i.render {
			size := i.Height - (i.padding[1] * 2)
			if size > len(i.Text) {
				size = len(i.Text)
			}
			i.TextOut = i.Text[len(i.Text) - size:]

			width := i.Width - (i.padding[0] * 2)
			for n := range i.TextOut {
				if len(i.TextOut[n]) > width {
					i.TextOut[n] = fmt.Sprintf("%s...", i.TextOut[n][:width - 3])
				}
			}
			ui.Render(i)
		}
	}()
}

func (i *TextView) readInputLoop() {
	go func() {
		for line := range i.inputStream {
			i.Text = append(i.Text, line)
			i.render <- true
		}
		close(i.render)
	}()
}
