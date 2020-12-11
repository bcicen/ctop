package compact

import (
	"fmt"

	"github.com/bcicen/ctop/cwidgets"
	"github.com/bcicen/ctop/models"

	ui "github.com/gizak/termui"
)

// Column that shows container's meta property i.e. name, id, image tc.
type MetaCol struct {
	*TextCol
	metaName string
}

func (w *MetaCol) SetMeta(m models.Meta) {
	w.setText(m.Get(w.metaName))
}

func NewNameCol() CompactCol {
	c := &MetaCol{NewTextCol("NAME"), "name"}
	c.fWidth = 30
	return c
}

func NewCIDCol() CompactCol {
	c := &MetaCol{NewTextCol("CID"), "id"}
	c.fWidth = 12
	return c
}

type NetCol struct {
	*TextCol
}

func NewNetCol() CompactCol {
	return &NetCol{NewTextCol("NET RX/TX")}
}

func (w *NetCol) SetMetrics(m models.Metrics) {
	label := fmt.Sprintf("%s / %s", cwidgets.ByteFormat64Short(m.NetRx), cwidgets.ByteFormat64Short(m.NetTx))
	w.setText(label)
}

type IOCol struct {
	*TextCol
}

func NewIOCol() CompactCol {
	return &IOCol{NewTextCol("IO R/W")}
}

func (w *IOCol) SetMetrics(m models.Metrics) {
	label := fmt.Sprintf("%s / %s", cwidgets.ByteFormat64Short(m.IOBytesRead), cwidgets.ByteFormat64Short(m.IOBytesWrite))
	w.setText(label)
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
	w.setText(fmt.Sprintf("%d", m.Pids))
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

	return &TextCol{
		Par:    p,
		header: header,
		fWidth: 0,
	}
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

// TextCol implements CompactCol
func (w *TextCol) Reset()                    { w.setText("-") }
func (w *TextCol) SetMeta(models.Meta)       {}
func (w *TextCol) SetMetrics(models.Metrics) {}
func (w *TextCol) Header() string            { return w.header }
func (w *TextCol) FixedWidth() int           { return w.fWidth }

func (w *TextCol) setText(s string) {
	if w.fWidth > 0 && len(s) > w.fWidth {
		s = s[0:w.fWidth]
	}
	w.Text = s
}
