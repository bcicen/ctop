package expanded

import (
	ui "github.com/gizak/termui"
)

type Expanded struct {
	Info *ui.Table
	Net  *ExpandedNet
	Cpu  *ExpandedCpu
	Mem  *ExpandedMem
}

func NewExpanded(id, name string) *Expanded {
	return &Expanded{
		Info: NewInfo(id, name),
		Net:  NewExpandedNet(),
		Cpu:  NewExpandedCpu(),
		Mem:  NewExpandedMem(),
	}
}

func NewInfo(id, name string) *ui.Table {
	p := ui.NewTable()
	p.Rows = [][]string{
		[]string{"name", name},
		[]string{"id", id},
	}
	p.Height = 4
	p.Width = 50
	p.FgColor = ui.ColorWhite
	p.Seperator = false
	return p
}

func (w *Expanded) Buffer() ui.Buffer {
	buf := ui.NewBuffer()
	buf.Merge(w.Info.Buffer())
	buf.Merge(w.Cpu.Buffer())
	buf.Merge(w.Mem.Buffer())
	buf.Merge(w.Net.Buffer())
	return buf
}

func (w *Expanded) Reset()               {}
func (w *Expanded) SetY(_ int)           {}
func (w *Expanded) SetWidth(_ int)       {}
func (w *Expanded) Highlight()           {}
func (w *Expanded) UnHighlight()         {}
func (w *Expanded) SetStatus(val string) {}

func (w *Expanded) SetCPU(val int) {
	w.Cpu.Update(val)
}

func (w *Expanded) SetNet(rx int64, tx int64) {
	w.Net.Update(rx, tx)
}

func (w *Expanded) SetMem(val int64, limit int64, percent int) {
	w.Mem.Update(int(val), int(limit))
}
