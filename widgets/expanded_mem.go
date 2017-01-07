package widgets

import (
	ui "github.com/gizak/termui"
)

type ExpandedMem struct {
	*ui.MBarChart
	valHist   IntHistData
	limitHist IntHistData
}

func NewExpandedMem() *ExpandedMem {
	mem := &ExpandedMem{
		ui.NewMBarChart(),
		NewIntHistData(8),
		NewIntHistData(8),
	}
	mem.BorderLabel = "MEM"
	mem.Height = 10
	mem.Width = 50
	mem.BarWidth = 5
	mem.BarGap = 1
	mem.X = 51
	mem.Y = 4
	mem.TextColor = ui.ColorDefault
	mem.Data[0] = mem.valHist.data
	mem.Data[0] = mem.valHist.data
	mem.Data[1] = mem.limitHist.data
	mem.BarColor[0] = ui.ColorGreen
	mem.BarColor[1] = ui.ColorBlack
	mem.DataLabels = mem.valHist.labels
	//mem.ShowScale = true
	return mem
}

func (w *ExpandedMem) Update(val int, limit int) {
	w.valHist.Append(val)
	w.limitHist.Append(limit - val)
	//w.Data[0] = w.hist.data
}
