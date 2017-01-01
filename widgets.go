package main

import (
	"fmt"
	"strconv"

	ui "github.com/gizak/termui"
)

type Widgets struct {
	cid    *ui.Par
	net    *ui.Par
	name   *ui.Par
	cpu    *ui.Gauge
	memory *ui.Gauge
}

func (w *Widgets) MakeRow() *ui.Row {
	return ui.NewRow(
		ui.NewCol(2, 0, w.name),
		ui.NewCol(2, 0, w.cid),
		ui.NewCol(2, 0, w.cpu),
		ui.NewCol(2, 0, w.memory),
		ui.NewCol(2, 0, w.net),
	)
}

func (w *Widgets) SetCPU(val int) {
	w.cpu.BarColor = colorScale(val)
	w.cpu.Label = fmt.Sprintf("%s%%", strconv.Itoa(val))
	if val < 5 {
		val = 5
		w.cpu.BarColor = ui.ColorBlack
	}
	w.cpu.Percent = val
}

func (w *Widgets) SetNet(rx int64, tx int64) {
	//w.net.Label = fmt.Sprintf("%s / %s", byteFormat(rx), byteFormat(tx))
	w.net.Text = fmt.Sprintf("%s / %s", byteFormat(rx), byteFormat(tx))
	//w.net2.Lines[0].Data = []int{0, 2, 5, 10, 20, 20, 2, 2, 0, 0}
}

func (w *Widgets) SetMem(val int64, limit int64) {
	percent := round((float64(val) / float64(limit)) * 100)
	w.memory.Label = fmt.Sprintf("%s / %s", byteFormat(val), byteFormat(limit))
	if percent < 5 {
		percent = 5
		w.memory.BarColor = ui.ColorBlack
	} else {
		w.memory.BarColor = ui.ColorGreen
	}
	w.memory.Percent = percent
}

func NewWidgets(id string, names string) *Widgets {

	cid := ui.NewPar(id)
	cid.Border = false
	cid.Height = 1
	cid.Width = 20
	cid.TextFgColor = ui.ColorWhite

	name := ui.NewPar(names)
	name.Border = false
	name.Height = 1
	name.Width = 20
	name.TextFgColor = ui.ColorWhite

	net := ui.NewPar("-")
	net.Border = false
	net.Height = 1
	net.Width = 20
	net.TextFgColor = ui.ColorWhite

	return &Widgets{cid, net, name, mkGauge(), mkGauge()}
}

func mkGauge() *ui.Gauge {
	g := ui.NewGauge()
	g.Height = 1
	g.Border = false
	g.Percent = 0
	g.PaddingBottom = 0
	g.BarColor = ui.ColorGreen
	g.Label = "-"
	return g
}
