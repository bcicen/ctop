package compact

import (
	"fmt"
	"github.com/bcicen/ctop/cwidgets"
	"strconv"

	ui "github.com/gizak/termui"
)

const (
	mark        = string('\u25C9')
	vBar        = string('\u25AE')
	colSpacing  = 1
	statusWidth = 3
)

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
	autoWidth := calcWidth(width, 5)
	for n, col := range w.all() {
		if n == 0 {
			col.SetX(x)
			col.SetWidth(statusWidth)
			x += statusWidth
			continue
		}
		col.SetX(x)
		col.SetWidth(autoWidth)
		x += autoWidth + colSpacing
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
	w.Cpu.BarColor = cwidgets.ColorScale(val)
	w.Cpu.Label = fmt.Sprintf("%s%%", strconv.Itoa(val))
	if val < 5 {
		val = 5
		w.Cpu.BarColor = ui.ColorBlack
	}
	w.Cpu.Percent = val
}

func (w *Compact) SetNet(rx int64, tx int64) {
	w.Net.Text = fmt.Sprintf("%s / %s", cwidgets.ByteFormat(rx), cwidgets.ByteFormat(tx))
}

func (w *Compact) SetMem(val int64, limit int64, percent int) {
	w.Memory.Label = fmt.Sprintf("%s / %s", cwidgets.ByteFormat(val), cwidgets.ByteFormat(limit))
	if percent < 5 {
		percent = 5
		w.Memory.BarColor = ui.ColorBlack
	} else {
		w.Memory.BarColor = ui.ColorGreen
	}
	w.Memory.Percent = percent
}
