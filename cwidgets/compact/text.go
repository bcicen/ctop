package compact

import (
	"fmt"

	"github.com/bcicen/ctop/cwidgets"
	"github.com/bcicen/ctop/models"
	ui "github.com/gizak/termui"
)

type NameCol struct {
	*TextCol
}

func (w *NameCol) SetMeta(m models.Meta) {
	if s, ok := m["name"]; ok {
		w.Text = s
	}
}

func (w *NameCol) SetMetrics(m models.Metrics) {
}

type CIDCol struct {
	*TextCol
}

type NetCol struct {
	*TextCol
}

func (w *NetCol) SetMetrics(m models.Metrics) {
	label := fmt.Sprintf("%s / %s", cwidgets.ByteFormat(m.NetRx), cwidgets.ByteFormat(m.NetTx))
	w.Text = label
}

type IOCol struct {
	*TextCol
}

func (w *IOCol) SetMetrics(m models.Metrics) {
	label := fmt.Sprintf("%s / %s", cwidgets.ByteFormat(m.IOBytesRead), cwidgets.ByteFormat(m.IOBytesWrite))
	w.Text = label
}

type PIDCol struct {
	*TextCol
}

func (w *PIDCol) SetMetrics(m models.Metrics) {
	w.Text = fmt.Sprintf("%d", m.Pids)
}

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
	w.Bg = ui.ThemeAttr("par.text.fg")
	w.TextFgColor = ui.ThemeAttr("par.text.hi")
	w.TextBgColor = ui.ThemeAttr("par.text.fg")
}

func (w *TextCol) UnHighlight() {
	w.Bg = ui.ThemeAttr("par.text.bg")
	w.TextFgColor = ui.ThemeAttr("par.text.fg")
	w.TextBgColor = ui.ThemeAttr("par.text.bg")
}

//func (w *TextCol) Set(s string) { w.Text = s }
func (w *TextCol) Reset()                    { w.Text = "-" }
func (w *TextCol) SetMeta(models.Meta)       {}
func (w *TextCol) SetMetrics(models.Metrics) {}
