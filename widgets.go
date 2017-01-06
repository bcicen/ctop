package main

import (
	"fmt"

	"github.com/bcicen/ctop/widgets"
	ui "github.com/gizak/termui"
)

type Widgets struct {
	cid    *ui.Par
	net    *ui.Par
	name   *ui.Par
	cpu    *widgets.CPU
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

func (w *Widgets) SetNet(rx int64, tx int64) {
	w.net.Text = fmt.Sprintf("%s / %s", byteFormat(rx), byteFormat(tx))
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

	return &Widgets{cid, net, name, widgets.NewCPU(), mkGauge()}
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
