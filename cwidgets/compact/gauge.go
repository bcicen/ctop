package compact

import (
	"fmt"

	"github.com/bcicen/ctop/cwidgets"
	"github.com/bcicen/ctop/models"
	ui "github.com/gizak/termui"
)

type CPUCol struct {
	*GaugeCol
}

func (w *CPUCol) SetMetrics(m models.Metrics) {
	val := m.CPUUtil
	w.BarColor = colorScale(val)
	w.Label = fmt.Sprintf("%d%%", val)

	if val > 100 {
		val = 100
	}
	w.Percent = val
}

type MemCol struct {
	*GaugeCol
}

func (w *MemCol) SetMetrics(m models.Metrics) {
	w.BarColor = ui.ThemeAttr("gauge.bar.bg")
	w.Label = fmt.Sprintf("%s / %s", cwidgets.ByteFormat(m.MemUsage), cwidgets.ByteFormat(m.MemLimit))
	w.Percent = m.MemPercent
}

type GaugeCol struct {
	*ui.Gauge
}

func NewGaugeCol() *GaugeCol {
	g := &GaugeCol{ui.NewGauge()}
	g.Height = 1
	g.Border = false
	g.PaddingBottom = 0
	g.Reset()
	return g
}

func (w *GaugeCol) Reset() {
	w.Label = "-"
	w.Percent = 0
}

func (w *GaugeCol) Buffer() ui.Buffer {
	// if bar would not otherwise be visible, set a minimum
	// percentage value and low-contrast color for structure
	if w.Percent < 5 {
		w.Percent = 5
		w.BarColor = ui.ColorBlack
	}

	return w.Gauge.Buffer()
}

// GaugeCol implements CompactCol
func (w *GaugeCol) SetMeta(models.Meta) {}

// GaugeCol implements CompactCol
func (w *GaugeCol) SetMetrics(models.Metrics) {}

// GaugeCol implements CompactCol
func (w *GaugeCol) Highlight() {
	w.Bg = ui.ThemeAttr("par.text.fg")
	w.PercentColor = ui.ThemeAttr("par.text.hi")
}

// GaugeCol implements CompactCol
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
