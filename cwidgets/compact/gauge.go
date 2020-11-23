package compact

import (
	"fmt"

	"github.com/bcicen/ctop/cwidgets"
	"github.com/bcicen/ctop/models"

	ui "github.com/gizak/termui"
)

type CPUCol struct {
	*GaugeCol
	scaleCpu bool
}

func NewCPUCol() CompactCol {
	return &CPUCol{NewGaugeCol("CPU"), false}
}

func NewCpuScaledCol() CompactCol {
	return &CPUCol{NewGaugeCol("CPUS"), true}
}

func (w *CPUCol) SetMetrics(m models.Metrics) {
	val := m.CPUUtil
	w.BarColor = colorScale(val)
	if !w.scaleCpu {
		val = val * int(m.NCpus)
	}
	w.Label = fmt.Sprintf("%d%%", val)

	if val > 100 {
		val = 100
	}
	w.Percent = val
}

type MemCol struct {
	*GaugeCol
}

func NewMemCol() CompactCol {
	return &MemCol{NewGaugeCol("MEM")}
}

func (w *MemCol) SetMetrics(m models.Metrics) {
	w.BarColor = ui.ThemeAttr("gauge.bar.bg")
	w.Label = fmt.Sprintf("%s / %s", cwidgets.ByteFormat64Short(m.MemUsage), cwidgets.ByteFormat64Short(m.MemLimit))
	w.Percent = m.MemPercent
}

type GaugeCol struct {
	*ui.Gauge
	header string
	fWidth int
}

func NewGaugeCol(header string) *GaugeCol {
	g := &GaugeCol{ui.NewGauge(), header, 0}
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
func (w *GaugeCol) SetMeta(models.Meta)       {}
func (w *GaugeCol) SetMetrics(models.Metrics) {}
func (w *GaugeCol) Header() string            { return w.header }
func (w *GaugeCol) FixedWidth() int           { return w.fWidth }

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
	if n <= 70 {
		return ui.ThemeAttr("status.ok")
	}
	if n <= 90 {
		return ui.ThemeAttr("status.warn")
	}
	return ui.ThemeAttr("status.danger")
}
