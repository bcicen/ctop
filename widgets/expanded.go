package widgets

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

func (w *Expanded) Reset() {
}

func (w *Expanded) Render() {
	ui.Render(w.Info, w.Cpu, w.Mem, w.Net)
	ui.Handle("/timer/1s", func(ui.Event) {
		ui.Render(w.Info, w.Cpu, w.Mem, w.Net)
	})
	ui.Handle("/sys/kbd/", func(ui.Event) {
		ui.StopLoop()
	})
	ui.Loop()
}

func (w *Expanded) Row() *ui.Row {
	return ui.NewRow(
		ui.NewCol(2, 0, w.Cpu),
		ui.NewCol(2, 0, w.Mem),
		ui.NewCol(2, 0, w.Net),
	)
}

func (w *Expanded) Highlight() {
}

func (w *Expanded) UnHighlight() {
}

func (w *Expanded) SetStatus(val string) {
}

func (w *Expanded) SetCPU(val int) {
	w.Cpu.Update(val)
}

func (w *Expanded) SetNet(rx int64, tx int64) {
	w.Net.Update(rx, tx)
}

func (w *Expanded) SetMem(val int64, limit int64, percent int) {
	w.Mem.Update(int(val), int(limit))
}
