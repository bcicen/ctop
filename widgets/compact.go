package widgets

import (
	"fmt"
	"strconv"

	ui "github.com/gizak/termui"
)

const (
	mark = string('\u25C9')
	vBar = string('\u25AE')
)

type ContainerWidgets interface {
	Row() *ui.Row
	Render()
	Highlight()
	UnHighlight()
	SetStatus(string)
	SetCPU(int)
	SetNet(int64, int64)
	SetMem(int64, int64, int)
}

type Compact struct {
	Status *ui.Par
	Cid    *ui.Par
	Net    *ui.Par
	Name   *ui.Par
	Cpu    *ui.Gauge
	Memory *ui.Gauge
}

func NewCompact(id string, name string) *Compact {
	return &Compact{
		Status: slimPar(""),
		Cid:    slimPar(id),
		Net:    slimPar("-"),
		Name:   slimPar(name),
		Cpu:    slimGauge(),
		Memory: slimGauge(),
	}
}

func (w *Compact) Render() {
}

func (w *Compact) Row() *ui.Row {
	return ui.NewRow(
		ui.NewCol(1, 0, w.Status),
		ui.NewCol(2, 0, w.Name),
		ui.NewCol(2, 0, w.Cid),
		ui.NewCol(2, 0, w.Cpu),
		ui.NewCol(2, 0, w.Memory),
		ui.NewCol(2, 0, w.Net),
	)
}

func (w *Compact) Highlight() {
	w.Name.TextFgColor = ui.ColorDefault
	w.Name.TextBgColor = ui.ColorWhite
}

func (w *Compact) UnHighlight() {
	w.Name.TextFgColor = ui.ColorWhite
	w.Name.TextBgColor = ui.ColorDefault
}

func (w *Compact) SetStatus(val string) {
	switch val {
	case "running":
		w.Status.Text = mark
		w.Status.TextFgColor = ui.ColorGreen
	case "exited":
		w.Status.Text = mark
		w.Status.TextFgColor = ui.ColorRed
	case "paused":
		w.Status.Text = fmt.Sprintf("%s%s", vBar, vBar)
		w.Status.TextFgColor = ui.ColorDefault
	default:
		w.Status.Text = mark
		w.Status.TextFgColor = ui.ColorRed
	}
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

func centerParText(p *ui.Par) {
	var text string
	var padding string

	// strip existing left-padding
	for i, ch := range p.Text {
		if string(ch) != " " {
			text = p.Text[i:]
			break
		}
	}

	padlen := (p.InnerWidth() - len(text)) / 2
	for i := 0; i < padlen; i++ {
		padding += " "
	}
	p.Text = fmt.Sprintf("%s%s", padding, text)
}

func slimPar(s string) *ui.Par {
	p := ui.NewPar(s)
	p.Border = false
	p.Height = 1
	p.Width = 20
	p.TextFgColor = ui.ColorWhite
	return p
}

func slimGauge() *ui.Gauge {
	g := ui.NewGauge()
	g.Height = 1
	g.Border = false
	g.Percent = 0
	g.PaddingBottom = 0
	g.BarColor = ui.ColorGreen
	g.Label = "-"
	return g
}
