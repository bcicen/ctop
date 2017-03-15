package compact

import (
	ui "github.com/gizak/termui"
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
	w.TextFgColor = ui.ThemeAttr("par.text.hi")
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
