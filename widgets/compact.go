package widgets

import (
	"fmt"
	"strconv"

	ui "github.com/gizak/termui"
)

type Compact struct {
	Cid    *ui.Par
	Net    *ui.Par
	Name   *ui.Par
	Cpu    *ui.Gauge
	Memory *ui.Gauge
}

func NewCompact(id string, name string) *Compact {
	return &Compact{
		Cid:    compactPar(id),
		Net:    compactPar("-"),
		Name:   compactPar(name),
		Cpu:    mkGauge(),
		Memory: mkGauge(),
	}
}

func (w *Compact) Row() *ui.Row {
	return ui.NewRow(
		ui.NewCol(2, 0, w.Name),
		ui.NewCol(2, 0, w.Cid),
		ui.NewCol(2, 0, w.Cpu),
		ui.NewCol(2, 0, w.Memory),
		ui.NewCol(2, 0, w.Net),
	)
}

func (w *Compact) SetCPU(val int) {
	w.Cpu.BarColor = colorScale(val)
	w.Cpu.Label = fmt.Sprintf("%s%%", strconv.Itoa(val))
	if val < 5 {
		val = 5
		w.Cpu.BarColor = ui.ColorBlack
	}
	w.Cpu.Percent = val
}

func (w *Compact) SetNet(rx int64, tx int64) {
	w.Net.Text = fmt.Sprintf("%s / %s", byteFormat(rx), byteFormat(tx))
}

func (w *Compact) SetMem(val int64, limit int64, percent int) {
	w.Memory.Label = fmt.Sprintf("%s / %s", byteFormat(val), byteFormat(limit))
	if percent < 5 {
		percent = 5
		w.Memory.BarColor = ui.ColorBlack
	} else {
		w.Memory.BarColor = ui.ColorGreen
	}
	w.Memory.Percent = percent
}
