package widgets

import (
	"fmt"
	ui "github.com/gizak/termui"
)

type ErrorView struct {
	*ui.Par
}

func NewErrorView() *ErrorView {
	p := ui.NewPar("")
	p.Border = true
	p.Height = 10
	p.Width = 20
	p.PaddingTop = 1
	p.PaddingBottom = 1
	p.PaddingLeft = 2
	p.PaddingRight = 2
	p.Bg = ui.ThemeAttr("bg")
	p.TextFgColor = ui.ThemeAttr("status.warn")
	p.TextBgColor = ui.ThemeAttr("menu.text.bg")
	p.BorderFg = ui.ThemeAttr("status.warn")
	p.BorderLabelFg = ui.ThemeAttr("status.warn")
	return &ErrorView{p}
}

func (w *ErrorView) Buffer() ui.Buffer {
	w.BorderLabel = fmt.Sprintf(" %s ", timeStr())
	return w.Par.Buffer()
}

func (w *ErrorView) Resize() {
	w.SetX(ui.TermWidth() / 12)
	w.SetY(ui.TermHeight() / 3)
	w.SetWidth(w.X * 10)
}
