package compact

import (
	ui "github.com/gizak/termui"
)

type TextCol struct {
	*ui.Par
	isHighlight bool
}

func NewTextCol(s string) *TextCol {
	p := ui.NewPar(s)
	p.Border = false
	p.Height = 1
	p.Width = 20
	return &TextCol{p, false}
}

func (w *TextCol) Highlight() {
	if w.TextFgColor ==ui.ThemeAttr("par.text.fg"){
		w.TextFgColor = ui.ThemeAttr("par.text.hi")
	}
	w.TextBgColor = ui.ThemeAttr("par.text.fg")
	w.isHighlight = true
}

func (w *TextCol) UnHighlight() {
	if w.TextFgColor == ui.ThemeAttr("par.text.hi"){
		w.TextFgColor = ui.ThemeAttr("par.text.fg")
	}
	w.TextBgColor = ui.ThemeAttr("par.text.bg")
	w.isHighlight = false
}

func (w *TextCol) Reset() {
	w.Text = "-"
}

func (w *TextCol) Set(s string) {
	w.Text = s
}

func (w *TextCol) Color(s string){
	color := ui.ThemeAttr("par.text.fg")
	if w.isHighlight{
		color = ui.ThemeAttr("par.text.hi")
	}
	switch s {
	case "healthy":
		color = ui.ColorGreen
	case "unhealthy":
		color = ui.ColorMagenta
	case "starting":
		color = ui.ColorYellow
	}
	w.TextFgColor = color
}
