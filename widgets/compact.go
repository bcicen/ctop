package widgets

import (
	"fmt"
	"strconv"

	"github.com/bcicen/ctop/logging"
	ui "github.com/gizak/termui"
)

var log = logging.Init()

const (
	mark = string('\u25C9')
	vBar = string('\u25AE')
)

type CompactGrid struct {
	ui.GridBufferer
	Rows   []ContainerWidgets
	X, Y   int
	Width  int
	Height int
}

func (c CompactGrid) SetX(x int) { c.X = x }
func (c CompactGrid) SetY(y int) {
	c.Y = y
	for n, r := range c.Rows {
		log.Infof("row %d: y=%d", n, c.Y+n)
		r.SetY(c.Y + n)
	}
}
func (c CompactGrid) GetHeight() int { return len(c.Rows) }
func (c CompactGrid) SetWidth(w int) {
	c.Width = w
	for _, r := range c.Rows {
		r.SetWidth(w)
	}
}

func (c CompactGrid) Buffer() ui.Buffer {
	buf := ui.NewBuffer()
	for _, r := range c.Rows {
		buf.Merge(r.Buffer())
	}
	return buf
}

type ContainerWidgets interface {
	Render(int, int)
	Reset()
	Buffer() ui.Buffer
	Highlight()
	UnHighlight()
	SetY(int)
	SetWidth(int)
	SetStatus(string)
	SetCPU(int)
	SetNet(int64, int64)
	SetMem(int64, int64, int)
}

type CompactHeader struct {
	pars []*ui.Par
}

func NewCompactHeader() *CompactHeader {
	fields := []string{"", "NAME", "CID", "CPU", "MEM", "NET RX/TX"}
	header := &CompactHeader{}
	for _, f := range fields {
		header.pars = append(header.pars, slimHeaderPar(f))
	}
	return header
}

func (c *CompactHeader) Buffer() ui.Buffer {
	buf := ui.NewBuffer()
	for _, p := range c.pars {
		buf.Merge(p.Buffer())
	}
	return buf
}

type Compact struct {
	Status *ui.Par
	Cid    *ui.Par
	Net    *ui.Par
	Name   *ui.Par
	Cpu    *ui.Gauge
	Memory *ui.Gauge
}

func NewCompact(id, name, status string) *Compact {
	w := &Compact{
		Status: slimPar(mark),
		Cid:    slimPar(id),
		Name:   slimPar(name),
	}
	w.Reset()
	w.SetStatus(status)
	return w
}

// Set gauges, counters to default unread values
func (w *Compact) Reset() {
	w.Net = slimPar("-")
	w.Cpu = slimGauge()
	w.Memory = slimGauge()
}

func (w *Compact) all() []ui.GridBufferer {
	return []ui.GridBufferer{
		w.Status,
		w.Name,
		w.Cid,
		w.Cpu,
		w.Memory,
		w.Net,
	}
}

func (w *Compact) SetY(y int) {
	for _, col := range w.all() {
		col.SetY(y)
	}
}

func (w *Compact) SetWidth(width int) {
	x := 1
	statusWidth := 3
	autoWidth := (width - statusWidth) / 5
	log.Infof("autowidth: %d", autoWidth)
	for n, col := range w.all() {
		if n == 0 {
			col.SetX(x)
			col.SetWidth(statusWidth)
			x += statusWidth
			continue
		}
		col.SetX(x)
		col.SetWidth(autoWidth)
		x += autoWidth
	}
}

func (w *Compact) Render(y, rowWidth int) {}

func (w *Compact) Buffer() ui.Buffer {
	buf := ui.NewBuffer()

	buf.Merge(w.Status.Buffer())
	buf.Merge(w.Name.Buffer())
	buf.Merge(w.Cid.Buffer())
	buf.Merge(w.Cpu.Buffer())
	buf.Merge(w.Memory.Buffer())
	buf.Merge(w.Net.Buffer())

	return buf
}

func (w *Compact) Highlight() {
	w.Name.TextFgColor = ui.ColorDefault
	w.Name.TextBgColor = ui.ColorWhite
}

func (w *Compact) UnHighlight() {
	w.Name.TextFgColor = ui.ColorWhite
	w.Name.TextBgColor = ui.ColorDefault
}

func (w *Compact) SetStatus(val string) {
	switch val {
	case "running":
		w.Status.Text = mark
		w.Status.TextFgColor = ui.ColorGreen
	case "exited":
		w.Status.Text = mark
		w.Status.TextFgColor = ui.ColorRed
	case "paused":
		w.Status.Text = fmt.Sprintf("%s%s", vBar, vBar)
		w.Status.TextFgColor = ui.ColorDefault
	default:
		w.Status.Text = mark
		w.Status.TextFgColor = ui.ColorRed
	}
}

func (w *Compact) SetCPU(val int) {
	w.Cpu.BarColor = colorScale(val)
	w.Cpu.Label = fmt.Sprintf("%s%%", strconv.Itoa(val))
	if val < 5 {
		val = 5
		w.Cpu.BarColor = ui.ColorBlack
	}
	w.Cpu.Percent = val
}

func (w *Compact) SetNet(rx int64, tx int64) {
	w.Net.Text = fmt.Sprintf("%s / %s", byteFormat(rx), byteFormat(tx))
}

func (w *Compact) SetMem(val int64, limit int64, percent int) {
	w.Memory.Label = fmt.Sprintf("%s / %s", byteFormat(val), byteFormat(limit))
	if percent < 5 {
		percent = 5
		w.Memory.BarColor = ui.ColorBlack
	} else {
		w.Memory.BarColor = ui.ColorGreen
	}
	w.Memory.Percent = percent
}

func centerParText(p *ui.Par) {
	var text string
	var padding string

	// strip existing left-padding
	for i, ch := range p.Text {
		if string(ch) != " " {
			text = p.Text[i:]
			break
		}
	}

	padlen := (p.InnerWidth() - len(text)) / 2
	for i := 0; i < padlen; i++ {
		padding += " "
	}
	p.Text = fmt.Sprintf("%s%s", padding, text)
}

func slimHeaderPar(s string) *ui.Par {
	p := slimPar(s)
	p.Y = 2
	p.Height = 2
	return p
}

func slimPar(s string) *ui.Par {
	p := ui.NewPar(s)
	p.Border = false
	p.Height = 1
	p.Width = 20
	p.TextFgColor = ui.ColorWhite
	return p
}

func slimGauge() *ui.Gauge {
	g := ui.NewGauge()
	g.Height = 1
	g.Border = false
	g.Percent = 0
	g.PaddingBottom = 0
	g.BarColor = ui.ColorGreen
	g.Label = "-"
	return g
}
