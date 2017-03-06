package expanded

import (
	ui "github.com/gizak/termui"
)

type ExpandedCpu struct {
	*ui.LineChart
	hist FloatHist
}

func NewExpandedCpu() *ExpandedCpu {
	cpu := &ExpandedCpu{ui.NewLineChart(), NewFloatHist(60)}
	cpu.BorderLabel = "CPU"
	cpu.Height = 10
	cpu.Width = 50
	cpu.X = 0
	cpu.Y = 4
	cpu.Data = cpu.hist.Data
	cpu.DataLabels = cpu.hist.Labels
	cpu.AxesColor = ui.ColorDefault
	cpu.LineColor = ui.ColorGreen
	return cpu
}

func (w *ExpandedCpu) Update(val int) {
	w.hist.Append(float64(val))
}
