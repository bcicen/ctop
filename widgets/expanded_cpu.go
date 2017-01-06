package widgets

import (
	ui "github.com/gizak/termui"
)

type ExpandedCpu struct {
	*ui.BarChart
	hist HistData
}

func NewExpandedCpu() *ExpandedCpu {
	cpu := &ExpandedCpu{ui.NewBarChart(), NewHistData(12)}
	cpu.BorderLabel = "CPU Util"
	cpu.Height = 10
	cpu.Width = 50
	cpu.BarColor = ui.ColorGreen
	cpu.BarWidth = 3
	cpu.BarGap = 1
	cpu.X = 0
	cpu.Y = 4
	cpu.Data = cpu.hist.data
	cpu.DataLabels = cpu.hist.labels
	return cpu
}

func (w *ExpandedCpu) Update(val int) {
	w.hist.Append(val)
	w.Data = w.hist.data
}
