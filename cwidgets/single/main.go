package single

import (
	"github.com/bcicen/ctop/logging"
	"github.com/bcicen/ctop/models"
	ui "github.com/gizak/termui"
)

var (
	log       = logging.Init()
	sizeError = termSizeError()
	colWidth  = [2]int{65, 0} // left,right column width
)

type Single struct {
	Info  *Info
	Net   *Net
	Cpu   *Cpu
	Mem   *Mem
	IO    *IO
	Env   *Env
	X, Y  int
	Width int
}

func NewSingle(id string) *Single {
	if len(id) > 12 {
		id = id[:12]
	}
	return &Single{
		Info:  NewInfo(id),
		Net:   NewNet(),
		Cpu:   NewCpu(),
		Mem:   NewMem(),
		IO:    NewIO(),
		Env:   NewEnv(),
		Width: ui.TermWidth(),
	}
}

func (e *Single) Up() {
	if e.Y < 0 {
		e.Y++
		e.Align()
		ui.Render(e)
	}
}

func (e *Single) Down() {
	if e.Y > (ui.TermHeight() - e.GetHeight()) {
		e.Y--
		e.Align()
		ui.Render(e)
	}
}

func (e *Single) SetWidth(w int) { e.Width = w }
func (e *Single) SetMeta(k, v string) {
	if k == "[ENV-VAR]" {
		e.Env.Set(k, v)
	} else {
		e.Info.Set(k, v)
	}
}

func (e *Single) SetMetrics(m models.Metrics) {
	e.Cpu.Update(m.CPUUtil)
	e.Net.Update(m.NetRx, m.NetTx)
	e.Mem.Update(int(m.MemUsage), int(m.MemLimit))
	e.IO.Update(m.IOBytesRead, m.IOBytesWrite)
}

// Return total column height
func (e *Single) GetHeight() (h int) {
	h += e.Info.Height
	h += e.Net.Height
	h += e.Cpu.Height
	h += e.Mem.Height
	h += e.IO.Height
	h += e.Env.Height
	return h
}

func (e *Single) Align() {
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

func (e *Single) Buffer() ui.Buffer {
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
	buf.Merge(e.Env.Buffer())
	return buf
}

func (e *Single) all() []ui.GridBufferer {
	return []ui.GridBufferer{
		e.Info,
		e.Cpu,
		e.Mem,
		e.Net,
		e.IO,
		e.Env,
	}
}

func termSizeError() *ui.Par {
	p := ui.NewPar("screen too small!")
	p.Height = 1
	p.Width = 20
	p.Border = false
	return p
}
