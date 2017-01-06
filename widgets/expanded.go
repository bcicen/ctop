package widgets

import (
	"fmt"

	ui "github.com/gizak/termui"
)

type Expanded struct {
	Info   *ui.Table
	Net    *ui.Par
	Cpu    *ExpandedCpu
	Memory *ui.Gauge
}

func NewExpanded(id, name string) *Expanded {
	return &Expanded{
		Info:   NewInfo(id, name),
		Net:    ui.NewPar("-"),
		Cpu:    NewExpandedCpu(),
		Memory: mkGauge(),
	}
}

func NewInfo(id, name string) *ui.Table {
	p := ui.NewTable()
	p.Rows = [][]string{
		[]string{"name", name},
		[]string{"id", id},
	}
	p.Height = 4
	p.Width = 40
	p.FgColor = ui.ColorWhite
	p.Seperator = false
	return p
}

func (w *Expanded) Render() {
	ui.Render(w.Info, w.Cpu)
	ui.Handle("/timer/1s", func(ui.Event) {
		ui.Render(w.Info, w.Cpu)
	})
	ui.Handle("/sys/kbd/", func(ui.Event) {
		ui.StopLoop()
	})
	ui.Loop()
}

func (w *Expanded) Row() *ui.Row {
	return ui.NewRow(
		ui.NewCol(2, 0, w.Cpu),
		ui.NewCol(2, 0, w.Memory),
		ui.NewCol(2, 0, w.Net),
	)
}

func (w *Expanded) Highlight() {
}

func (w *Expanded) UnHighlight() {
}

func (w *Expanded) SetCPU(val int) {
	w.Cpu.Update(val)
}

func (w *Expanded) SetNet(rx int64, tx int64) {
	w.Net.Text = fmt.Sprintf("%s / %s", byteFormat(rx), byteFormat(tx))
}

func (w *Expanded) SetMem(val int64, limit int64, percent int) {
	w.Memory.Label = fmt.Sprintf("%s / %s", byteFormat(val), byteFormat(limit))
	if percent < 5 {
		percent = 5
		w.Memory.BarColor = ui.ColorBlack
	} else {
		w.Memory.BarColor = ui.ColorGreen
	}
	w.Memory.Percent = percent
}
