package single

import (
	ui "github.com/gizak/termui"
)

type Cpu struct {
	*ui.LineChart
	hist FloatHist
}

func NewCpu() *Cpu {
	cpu := &Cpu{ui.NewLineChart(), NewFloatHist(55)}
	cpu.Mode = "dot"
	cpu.BorderLabel = "CPU"
	cpu.Height = 12
	cpu.Width = colWidth[0]
	cpu.X = 0
	cpu.DataLabels = cpu.hist.Labels

	// hack to force the default minY scale to 0
	tmpData := []float64{20}
	cpu.Data["CPU"] = tmpData
	_ = cpu.Buffer()

	cpu.Data["CPU"] = cpu.hist.Data
	return cpu
}

func (w *Cpu) Update(val int) {
	w.hist.Append(float64(val))
}
