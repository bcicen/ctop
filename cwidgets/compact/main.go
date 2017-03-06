package compact

import (
	"github.com/bcicen/ctop/logging"
	"github.com/bcicen/ctop/metrics"
	ui "github.com/gizak/termui"
)

var log = logging.Init()

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

func NewCompact(id string) *Compact {
	// truncate container id
	if len(id) > 12 {
		id = id[:12]
	}
	row := &Compact{
		Status: NewStatus(),
		Name:   NewTextCol("-"),
		Cid:    NewTextCol(id),
		Cpu:    NewGaugeCol(),
		Memory: NewGaugeCol(),
		Net:    NewTextCol("-"),
		Height: 1,
	}
	return row
}

//func (row *Compact) ToggleExpand() {
//if row.Height == 1 {
//row.Height = 4
//} else {
//row.Height = 1
//}
//}

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

func (row *Compact) SetY(y int) {
	for _, col := range row.all() {
		col.SetY(y)
	}
	row.Y = y
}

func (row *Compact) SetWidth(width int) {
	x := 1
	autoWidth := calcWidth(width, 5)
	for n, col := range row.all() {
		// set status column to static width
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
