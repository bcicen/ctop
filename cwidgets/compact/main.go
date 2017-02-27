package compact

import (
	"fmt"
	"strconv"

	"github.com/bcicen/ctop/cwidgets"
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
	Name   *ui.Par
	Cid    *ui.Par
	Cpu    *ui.Gauge
	Memory *ui.Gauge
	Net    *ui.Par
	X, Y   int
	Width  int
	Height int
}

func NewCompact(id, name, status string) *Compact {
	row := &Compact{
		Status: slimPar(mark),
		Name:   slimPar(name),
		Cid:    slimPar(id),
		Cpu:    slimGauge(),
		Memory: slimGauge(),
		Net:    slimPar("-"),
		Height: 1,
	}
	row.Reset()
	row.SetStatus(status)
	return row
}

// Set gauges, counters to default unread values
func (row *Compact) Reset() {
	row.Cpu.Percent = 0
	row.Cpu.Label = "-"
	row.Memory.Percent = 0
	row.Memory.Label = "-"
	row.Net.Text = "-"
}

func (row *Compact) all() []ui.GridBufferer {
	return []ui.GridBufferer{
		row.Status,
		row.Name,
		row.Cid,
		row.Cpu,
		row.Memory,
		row.Net,
	}
}

func (row *Compact) SetY(y int) {
	if y == row.Y {
		return
	}
	for _, col := range row.all() {
		col.SetY(y)
	}
	row.Y = y
}

func (row *Compact) SetWidth(width int) {
	if row.Width == width {
		return
	}
	x := 1
	autoWidth := calcWidth(width, 5)
	for n, col := range row.all() {
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
	row.Width = width
}

func (row *Compact) Buffer() ui.Buffer {
	buf := ui.NewBuffer()

	buf.Merge(row.Status.Buffer())
	buf.Merge(row.Name.Buffer())
	buf.Merge(row.Cid.Buffer())
	buf.Merge(row.Cpu.Buffer())
	buf.Merge(row.Memory.Buffer())
	buf.Merge(row.Net.Buffer())

	return buf
}

func (row *Compact) Highlight() {
	row.Name.TextFgColor = ui.ColorDefault
	row.Name.TextBgColor = ui.ColorWhite
}

func (row *Compact) UnHighlight() {
	row.Name.TextFgColor = ui.ColorWhite
	row.Name.TextBgColor = ui.ColorDefault
}

func (row *Compact) SetStatus(val string) {
	switch val {
	case "running":
		row.Status.Text = mark
		row.Status.TextFgColor = ui.ColorGreen
	case "exited":
		row.Status.Text = mark
		row.Status.TextFgColor = ui.ColorRed
	case "paused":
		row.Status.Text = fmt.Sprintf("%s%s", vBar, vBar)
		row.Status.TextFgColor = ui.ColorDefault
	default:
		row.Status.Text = mark
		row.Status.TextFgColor = ui.ColorRed
	}
}

func (row *Compact) SetCPU(val int) {
	row.Cpu.BarColor = cwidgets.ColorScale(val)
	row.Cpu.Label = fmt.Sprintf("%s%%", strconv.Itoa(val))
	if val < 5 {
		val = 5
		row.Cpu.BarColor = ui.ColorBlack
	}
	row.Cpu.Percent = val
}

func (row *Compact) SetNet(rx int64, tx int64) {
	row.Net.Text = fmt.Sprintf("%s / %s", cwidgets.ByteFormat(rx), cwidgets.ByteFormat(tx))
}

func (row *Compact) SetMem(val int64, limit int64, percent int) {
	row.Memory.Label = fmt.Sprintf("%s / %s", cwidgets.ByteFormat(val), cwidgets.ByteFormat(limit))
	if percent < 5 {
		percent = 5
		row.Memory.BarColor = ui.ColorBlack
	} else {
		row.Memory.BarColor = ui.ColorGreen
	}
	row.Memory.Percent = percent
}
