package compact

import (
	"github.com/bcicen/ctop/logging"
	"github.com/bcicen/ctop/models"
	ui "github.com/gizak/termui"
)

var log = logging.Init()

type Compact struct {
	Status *Status
	Name   *TextCol
	Cid    *TextCol
	Cpu    *GaugeCol
	Mem    *GaugeCol
	Net    *TextCol
	IO     *TextCol
	Pids   *TextCol
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
		Mem:    NewGaugeCol(),
		Net:    NewTextCol("-"),
		IO:     NewTextCol("-"),
		Pids:   NewTextCol("-"),
		X:      1,
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

func (row *Compact) SetMeta(k, v string) {
	switch k {
	case "name":
		row.Name.Set(v)
	case "state":
		row.Status.Set(v)
	case "health":
		row.Status.SetHealth(v)
	}
}

func (row *Compact) SetMetrics(m models.Metrics) {
	row.SetCPU(m.CPUUtil)
	row.SetNet(m.NetRx, m.NetTx)
	row.SetMem(m.MemUsage, m.MemLimit, m.MemPercent)
	row.SetIO(m.IOBytesRead, m.IOBytesWrite)
	row.SetPids(m.Pids)
}

// Set gauges, counters to default unread values
func (row *Compact) Reset() {
	row.Cpu.Reset()
	row.Mem.Reset()
	row.Net.Reset()
	row.IO.Reset()
	row.Pids.Reset()
}

func (row *Compact) GetHeight() int {
	return row.Height
}

func (row *Compact) SetX(x int) {
	row.X = x
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
	if width == row.Width {
		return
	}
	x := row.X
	autoWidth := calcWidth(width)
	for n, col := range row.all() {
		if colWidths[n] != 0 {
			col.SetX(x)
			col.SetWidth(colWidths[n])
			x += colWidths[n]
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
	buf.Merge(row.Mem.Buffer())
	buf.Merge(row.Net.Buffer())
	buf.Merge(row.IO.Buffer())
	buf.Merge(row.Pids.Buffer())
	return buf
}

func (row *Compact) all() []ui.GridBufferer {
	return []ui.GridBufferer{
		row.Status,
		row.Name,
		row.Cid,
		row.Cpu,
		row.Mem,
		row.Net,
		row.IO,
		row.Pids,
	}
}
