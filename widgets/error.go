package widgets

import (
	"fmt"
	"strings"
	"time"

	ui "github.com/gizak/termui"
)

type ErrorView struct {
	*ui.Par
	lines []string
}

func NewErrorView() *ErrorView {
	const yPad = 1
	const xPad = 2

	p := ui.NewPar("")
	p.X = xPad
	p.Y = yPad
	p.Border = true
	p.Height = 10
	p.Width = 20
	p.PaddingTop = yPad
	p.PaddingBottom = yPad
	p.PaddingLeft = xPad
	p.PaddingRight = xPad
	p.BorderLabel = " ctop - error "
	p.Bg = ui.ThemeAttr("bg")
	p.TextFgColor = ui.ThemeAttr("status.warn")
	p.TextBgColor = ui.ThemeAttr("menu.text.bg")
	p.BorderFg = ui.ThemeAttr("status.warn")
	p.BorderLabelFg = ui.ThemeAttr("status.warn")
	return &ErrorView{p, make([]string, 0, 50)}
}

func (w *ErrorView) Append(s string) {
	if len(w.lines)+2 >= cap(w.lines) {
		w.lines = append(w.lines[:0], w.lines[2:]...)
	}
	ts := time.Now().Local().Format("15:04:05 MST")
	w.lines = append(w.lines, fmt.Sprintf("[%s] %s", ts, s))
	w.lines = append(w.lines, "")
}

func (w *ErrorView) Buffer() ui.Buffer {
	offset := len(w.lines) - w.InnerHeight()
	if offset < 0 {
		offset = 0
	}
	w.Text = strings.Join(w.lines[offset:len(w.lines)], "\n")
	return w.Par.Buffer()
}

func (w *ErrorView) Resize() {
	w.Height = ui.TermHeight() - (w.PaddingTop + w.PaddingBottom)
	w.SetWidth(ui.TermWidth() - (w.PaddingLeft + w.PaddingRight))
}
