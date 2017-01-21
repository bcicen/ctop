package widgets

import (
	"strings"

	ui "github.com/gizak/termui"
)

var (
	input_chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ-_."
)

type Input struct {
	ui.Block
	Label       string
	Data        string
	TextFgColor ui.Attribute
	TextBgColor ui.Attribute
	padding     Padding
}

func NewInput() *Input {
	i := &Input{
		Block:       *ui.NewBlock(),
		Label:       "input",
		TextFgColor: ui.ThemeAttr("par.text.fg"),
		TextBgColor: ui.ThemeAttr("par.text.bg"),
		padding:     Padding{4, 2},
	}
	i.Width, i.Height = 30, 3
	return i
}

func (i *Input) Buffer() ui.Buffer {
	var cell ui.Cell
	buf := i.Block.Buffer()

	x := i.padding[0]
	for _, ch := range i.Data {
		cell = ui.Cell{Ch: ch, Fg: i.TextFgColor, Bg: i.TextBgColor}
		buf.Set(x, 1, cell)
		x++
	}

	return buf
}

func (i *Input) KeyPress(e ui.Event) {
	ch := strings.Replace(e.Path, "/sys/kbd/", "", -1)
	if ch == "C-8" {
		idx := len(i.Data) - 1
		if idx > -1 {
			i.Data = i.Data[0:idx]
		}
		ui.Render(i)
	}
	if strings.Index(input_chars, ch) > -1 {
		i.Data += ch
		ui.Render(i)
	}
}

// Setup some default handlers for menu navigation
func (i *Input) InputHandlers() {
	ui.Handle("/sys/kbd/", i.KeyPress)
}
