package widgets

import (
	"fmt"

	ui "github.com/gizak/termui"
)

type ExpandedNet struct {
	*ui.Sparklines
	rxHist IntHistData
	txHist IntHistData
}

func NewExpandedNet() *ExpandedNet {
	net := &ExpandedNet{ui.NewSparklines(), NewIntHistData(60), NewIntHistData(60)}

	rx := ui.NewSparkline()
	rx.Title = "RX"
	rx.Height = 3
	rx.Data = net.rxHist.data
	rx.TitleColor = ui.ColorDefault
	rx.LineColor = ui.ColorGreen

	tx := ui.NewSparkline()
	tx.Title = "TX"
	tx.Height = 3
	tx.Data = net.txHist.data
	tx.TitleColor = ui.ColorDefault
	tx.LineColor = ui.ColorGreen

	net.Lines = []ui.Sparkline{rx, tx}
	net.Height = 12
	net.Width = 50
	net.X = 0
	net.Y = 15
	return net
}

func (w *ExpandedNet) Update(rx int64, tx int64) {
	w.rxHist.Append(int(rx))
	w.txHist.Append(int(tx))
	w.Lines[0].Title = fmt.Sprintf("RX (%s)", byteFormat(rx))
	w.Lines[1].Title = fmt.Sprintf("TX (%s)", byteFormat(tx))
}
