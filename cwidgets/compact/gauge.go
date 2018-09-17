package compact

import (
	ui "github.com/gizak/termui"
)

type GaugeCol struct {
	*ui.Gauge
}

func NewGaugeCol() *GaugeCol {
	g := ui.NewGauge()
	g.Height = 1
	g.Border = false
	g.Percent = 0
	g.PaddingBottom = 0
	g.Label = "-"
	return &GaugeCol{g}
}

func (w *GaugeCol) Reset() {
	w.Label = "-"
	w.Percent = 0
}

func (w *GaugeCol) Highlight() {
	w.Bg = ui.ThemeAttr("par.text.fg")
	w.PercentColor = ui.ThemeAttr("par.text.hi")
}

func (w *GaugeCol) UnHighlight() {
	w.Bg = ui.ThemeAttr("par.text.bg")
	w.PercentColor = ui.ThemeAttr("par.text.bg")
}

func colorScale(n int) ui.Attribute {
	if n > 70 {
		return ui.ThemeAttr("status.danger")
	}
	if n > 30 {
		return ui.ThemeAttr("status.warn")
	}
	return ui.ThemeAttr("status.ok")
}
