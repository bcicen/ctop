package expanded

import (
	"github.com/bcicen/ctop/logging"
	"github.com/bcicen/ctop/metrics"
	ui "github.com/gizak/termui"
)

var (
	log       = logging.Init()
	sizeError = termSizeError()
	colWidth  = [2]int{65, 0} // left,right column width
)

type Expanded struct {
	Info  *Info
	Net   *Net
	Cpu   *Cpu
	Mem   *Mem
	IO    *IO
	X, Y  int
	Width int
}

func NewExpanded(id string) *Expanded {
	if len(id) > 12 {
		id = id[:12]
	}
	return &Expanded{
		Info:  NewInfo(id),
		Net:   NewNet(),
		Cpu:   NewCpu(),
		Mem:   NewMem(),
		IO:    NewIO(),
		Width: ui.TermWidth(),
	}
}

func (e *Expanded) Up() {
	if e.Y < 0 {
		e.Y++
		e.Align()
		ui.Render(e)
	}
}

func (e *Expanded) Down() {
	if e.Y > (ui.TermHeight() - e.GetHeight()) {
		e.Y--
		e.Align()
		ui.Render(e)
	}
}

func (e *Expanded) SetWidth(w int)      { e.Width = w }
func (e *Expanded) SetMeta(k, v string) { e.Info.Set(k, v) }

func (e *Expanded) SetMetrics(m metrics.Metrics) {
	e.Cpu.Update(m.CPUUtil)
	e.Net.Update(m.NetRx, m.NetTx)
	e.Mem.Update(int(m.MemUsage), int(m.MemLimit))
	e.IO.Update(m.IOBytesRead, m.IOBytesWrite)
}

// Return total column height
func (e *Expanded) GetHeight() (h int) {
	h += e.Info.Height
	h += e.Net.Height
	h += e.Cpu.Height
	h += e.Mem.Height
	h += e.IO.Height
	return h
}

func (e *Expanded) Align() {
	// reset offset if needed
	if e.GetHeight() <= ui.TermHeight() {
		e.Y = 0
	}

	y := e.Y
	for _, i := range e.all() {
		i.SetY(y)
		y += i.GetHeight()
	}

	if e.Width > colWidth[0] {
		colWidth[1] = e.Width - (colWidth[0] + 1)
	}
	e.Mem.Align()
	log.Debugf("align: width=%v left-col=%v right-col=%v", e.Width, colWidth[0], colWidth[1])
}

func calcWidth(w int) {
}

func (e *Expanded) Buffer() ui.Buffer {
	buf := ui.NewBuffer()
	if e.Width < (colWidth[0] + colWidth[1]) {
		ui.Clear()
		buf.Merge(sizeError.Buffer())
		return buf
	}
	buf.Merge(e.Info.Buffer())
	buf.Merge(e.Cpu.Buffer())
	buf.Merge(e.Mem.Buffer())
	buf.Merge(e.Net.Buffer())
	buf.Merge(e.IO.Buffer())
	return buf
}

func (e *Expanded) all() []ui.GridBufferer {
	return []ui.GridBufferer{
		e.Info,
		e.Cpu,
		e.Mem,
		e.Net,
		e.IO,
	}
}

func termSizeError() *ui.Par {
	p := ui.NewPar("screen too small!")
	p.Height = 1
	p.Width = 20
	p.Border = false
	return p
}
