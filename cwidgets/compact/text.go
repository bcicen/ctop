package compact

import (
	"fmt"

	"github.com/lordoverlord/ctop/cwidgets"
	"github.com/lordoverlord/ctop/models"

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
	c.minWidth = 25
	c.maxWidth = 40
	c.growPriority = 1 // Highest priority
	return c
}

func NewCIDCol() CompactCol {
	c := &MetaCol{NewTextCol("CID"), "id"}
	c.minWidth = 13
	c.maxWidth = 15
	c.growPriority = 2 // Lower priority
	return c
}

func NewImageCol() CompactCol {
	c := &MetaCol{NewTextCol("IMAGE"), "image"}
	c.minWidth = 15
	c.maxWidth = 35
	c.growPriority = 1 // Highest priority
	return c
}

func NewPortsCol() CompactCol {
	c := &MetaCol{NewTextCol("PORTS"), "ports"}
	c.minWidth = 10
	c.maxWidth = 25
	c.growPriority = 2 // Lower priority
	return c
}

func NewIpsCol() CompactCol {
	c := &MetaCol{NewTextCol("IPs"), "IPs"}
	c.minWidth = 10
	c.maxWidth = 20
	c.growPriority = 2 // Lower priority
	return c
}

func NewCreatedCol() CompactCol {
	c := &MetaCol{NewTextCol("CREATED"), "created"}
	c.fWidth = 19 // Fixed width - e.g. "Thu Nov 26 07:44:03" without year
	return c
}

type NetCol struct {
	*TextCol
}

func NewNetCol() CompactCol {
	c := &NetCol{NewTextCol("NET RX/TX")}
	c.minWidth = 12
	c.maxWidth = 14
	c.growPriority = 2 // Lower priority
	return c
}

func (w *NetCol) SetMetrics(m models.Metrics) {
	label := fmt.Sprintf("%s / %s", cwidgets.ByteFormat64Short(m.NetRx), cwidgets.ByteFormat64Short(m.NetTx))
	w.setText(label)
}

type IOCol struct {
	*TextCol
}

func NewIOCol() CompactCol {
	c := &IOCol{NewTextCol("IO R/W")}
	c.minWidth = 12
	c.maxWidth = 14
	c.growPriority = 2 // Lower priority
	return c
}

func (w *IOCol) SetMetrics(m models.Metrics) {
	// Show dash if BlockIO stats unavailable (requires kernel blkio cgroup support)
	if m.IOBytesRead == 0 && m.IOBytesWrite == 0 {
		w.setText("-")
	} else {
		label := fmt.Sprintf("%s / %s", cwidgets.ByteFormat64Short(m.IOBytesRead), cwidgets.ByteFormat64Short(m.IOBytesWrite))
		w.setText(label)
	}
}

type PIDCol struct {
	*TextCol
}

func NewPIDCol() CompactCol {
	w := &PIDCol{NewTextCol("PIDS")}
	w.minWidth = 4
	w.maxWidth = 6
	w.growPriority = 2 // Lower priority
	return w
}

func (w *PIDCol) SetMetrics(m models.Metrics) {
	// Show dash if PIDs stats unavailable (requires kernel cgroup support)
	if m.Pids == 0 {
		w.setText("-")
	} else {
		w.setText(fmt.Sprintf("%d", m.Pids))
	}
}

type UptimeCol struct {
	*TextCol
}

func NewUptimeCol() CompactCol {
	c := &UptimeCol{NewTextCol("UPTIME")}
	c.minWidth = 10
	c.maxWidth = 15
	c.growPriority = 2 // Lower priority
	return c
}

func (w *UptimeCol) SetMeta(m models.Meta) {
	w.Text = m.Get("uptime")
}

type TextCol struct {
	*ui.Par
	header       string
	fWidth       int // fixed width (0 = growable)
	minWidth     int // minimum width
	maxWidth     int // maximum width
	growPriority int // 0=fixed, 1=highest priority, 2=medium, 3=lowest
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
func (w *TextCol) MinWidth() int             { return w.minWidth }
func (w *TextCol) MaxWidth() int             { return w.maxWidth }
func (w *TextCol) GrowPriority() int         { return w.growPriority }

func (w *TextCol) setText(s string) {
	if w.fWidth > 0 && len(s) > w.fWidth {
		s = s[0:w.fWidth]
	}
	w.Text = s
}
