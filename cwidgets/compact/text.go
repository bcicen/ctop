package compact

import (
	"fmt"

	ui "github.com/gizak/termui"
)

const (
	mark        = string('\u25C9')
	vBar        = string('\u25AE')
	statusWidth = 3
)

type TextCol struct {
	*ui.Par
}

func NewTextCol(s string) *TextCol {
	p := ui.NewPar(s)
	p.Border = false
	p.Height = 1
	p.Width = 20
	return &TextCol{p}
}

func (w *TextCol) Highlight() {
	w.TextFgColor = ui.ThemeAttr("par.text.bg")
	w.TextBgColor = ui.ThemeAttr("par.text.fg")
}

func (w *TextCol) UnHighlight() {
	w.TextFgColor = ui.ThemeAttr("par.text.fg")
	w.TextBgColor = ui.ThemeAttr("par.text.bg")
}

func (w *TextCol) Reset() {
	w.Text = "-"
}

func (w *TextCol) Set(s string) {
	w.Text = s
}

type Status struct {
	*ui.Par
}

func NewStatus() *Status {
	p := ui.NewPar(mark)
	p.Border = false
	p.Height = 1
	p.Width = statusWidth
	return &Status{p}
}

func (s *Status) Set(val string) {
	// defaults
	text := mark
	color := ui.ColorDefault

	switch val {
	case "running":
		color = ui.ColorGreen
	case "exited":
		color = ui.ColorRed
	case "paused":
		text = fmt.Sprintf("%s%s", vBar, vBar)
	}

	s.Text = text
	s.TextFgColor = color
}
