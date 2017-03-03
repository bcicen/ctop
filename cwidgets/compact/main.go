package compact

import (
	"fmt"
	"strconv"

	"github.com/bcicen/ctop/cwidgets"
	"github.com/bcicen/ctop/logging"
	"github.com/bcicen/ctop/metrics"
	ui "github.com/gizak/termui"
)

var log = logging.Init()

const (
	colSpacing = 1
)

type Compact struct {
	Status *Status
	Name   *TextCol
	Cid    *TextCol
	Cpu    *GaugeCol
	Memory *GaugeCol
	Net    *TextCol
	X, Y   int
	Width  int
	Height int
}

func NewCompact(id, name string) *Compact {
	row := &Compact{
		Status: NewStatus(),
		Name:   NewTextCol(name),
		Cid:    NewTextCol(id),
		Cpu:    NewGaugeCol(),
		Memory: NewGaugeCol(),
		Net:    NewTextCol("-"),
		Height: 1,
	}
	return row
}

func (row *Compact) SetMetrics(m metrics.Metrics) {
	row.SetCPU(m.CPUUtil)
	row.SetNet(m.NetRx, m.NetTx)
	row.SetMem(m.MemUsage, m.MemLimit, m.MemPercent)
}

// Set gauges, counters to default unread values
func (row *Compact) Reset() {
	row.Cpu.Reset()
	row.Memory.Reset()
	row.Net.Reset()
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

func (row *Compact) SetNet(rx int64, tx int64) {
	label := fmt.Sprintf("%s / %s", cwidgets.ByteFormat(rx), cwidgets.ByteFormat(tx))
	row.Net.Set(label)
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
