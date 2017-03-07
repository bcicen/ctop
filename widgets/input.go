package widgets

import (
	"strings"

	ui "github.com/gizak/termui"
)

var (
	input_chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ-_."
)

type Padding [2]int // x,y padding

type Input struct {
	ui.Block
	Label       string
	Data        string
	MaxLen      int
	TextFgColor ui.Attribute
	TextBgColor ui.Attribute
	stream      chan string // stream text as it changes
	padding     Padding
}

func NewInput() *Input {
	i := &Input{
		Block:       *ui.NewBlock(),
		Label:       "input",
		MaxLen:      20,
		TextFgColor: ui.ThemeAttr("par.text.fg"),
		TextBgColor: ui.ThemeAttr("par.text.bg"),
		padding:     Padding{4, 2},
	}
	i.calcSize()
	return i
}

func (i *Input) calcSize() {
	i.Height = 3 // minimum height
	i.Width = i.MaxLen + (i.padding[0] * 2)
}

func (i *Input) Buffer() ui.Buffer {
	var cell ui.Cell
	buf := i.Block.Buffer()

	x := i.Block.X + i.padding[0]
	y := i.Block.Y + 1
	for _, ch := range i.Data {
		cell = ui.Cell{Ch: ch, Fg: i.TextFgColor, Bg: i.TextBgColor}
		buf.Set(x, y, cell)
		x++
	}

	return buf
}

func (i *Input) Stream() chan string {
	i.stream = make(chan string)
	return i.stream
}

func (i *Input) KeyPress(e ui.Event) {
	ch := strings.Replace(e.Path, "/sys/kbd/", "", -1)
	if ch == "C-8" {
		idx := len(i.Data) - 1
		if idx > -1 {
			i.Data = i.Data[0:idx]
			i.stream <- i.Data
		}
		ui.Render(i)
		return
	}
	if len(i.Data) >= i.MaxLen {
		return
	}
	if strings.Index(input_chars, ch) > -1 {
		i.Data += ch
		i.stream <- i.Data
		ui.Render(i)
	}
}

// Setup some default handlers for menu navigation
func (i *Input) InputHandlers() {
	ui.Handle("/sys/kbd/", i.KeyPress)
}
