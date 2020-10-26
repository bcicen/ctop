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

func NewNameCol() CompactCol {
	return &NameCol{NewTextCol("NAME")}
}

func (w *NameCol) SetMeta(m models.Meta) {
	w.Text = m.Get("name")
}

type CIDCol struct {
	*TextCol
}

func NewCIDCol() CompactCol {
	return &CIDCol{NewTextCol("CID")}
}

func (w *CIDCol) SetMeta(m models.Meta) {
	w.Text = m.Get("id")
}

type NetCol struct {
	*TextCol
}

func NewNetCol() CompactCol {
	return &NetCol{NewTextCol("NET RX/TX")}
}

func (w *NetCol) SetMetrics(m models.Metrics) {
	label := fmt.Sprintf("%s / %s", cwidgets.ByteFormat64Short(m.NetRx), cwidgets.ByteFormat64Short(m.NetTx))
	w.Text = label
}

type IOCol struct {
	*TextCol
}

func NewIOCol() CompactCol {
	return &IOCol{NewTextCol("IO R/W")}
}

func (w *IOCol) SetMetrics(m models.Metrics) {
	label := fmt.Sprintf("%s / %s", cwidgets.ByteFormat64Short(m.IOBytesRead), cwidgets.ByteFormat64Short(m.IOBytesWrite))
	w.Text = label
}

type PIDCol struct {
	*TextCol
}

func NewPIDCol() CompactCol {
	w := &PIDCol{NewTextCol("PIDS")}
	w.fWidth = 4
	return w
}

func (w *PIDCol) SetMetrics(m models.Metrics) {
	w.Text = fmt.Sprintf("%d", m.Pids)
}

type TextCol struct {
	*ui.Par
	header string
	fWidth int
}

func NewTextCol(header string) *TextCol {
	p := ui.NewPar("-")
	p.Border = false
	p.Height = 1
	p.Width = 20
	return &TextCol{p, header, 0}
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

func (w *TextCol) Reset()                    { w.Text = "-" }
func (w *TextCol) SetMeta(models.Meta)       {}
func (w *TextCol) SetMetrics(models.Metrics) {}
func (w *TextCol) Header() string            { return w.header }
func (w *TextCol) FixedWidth() int           { return w.fWidth }
