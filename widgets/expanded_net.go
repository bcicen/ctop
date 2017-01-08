package widgets

import (
	"fmt"

	ui "github.com/gizak/termui"
)

type ExpandedNet struct {
	*ui.Sparklines
	rxHist DiffHistData
	txHist DiffHistData
}

func NewExpandedNet() *ExpandedNet {
	net := &ExpandedNet{ui.NewSparklines(), NewDiffHistData(30), NewDiffHistData(30)}

	rx := ui.NewSparkline()
	rx.Title = "RX"
	rx.Height = 2
	rx.Data = net.rxHist.data
	rx.TitleColor = ui.ColorDefault
	rx.LineColor = ui.ColorGreen

	tx := ui.NewSparkline()
	tx.Title = "TX"
	tx.Height = 2
	tx.Data = net.txHist.data
	tx.TitleColor = ui.ColorDefault
	tx.LineColor = ui.ColorYellow

	net.Lines = []ui.Sparkline{rx, tx}
	net.Height = 8
	net.Width = 35
	net.X = 0
	net.Y = 15
	return net
}

func (w *ExpandedNet) Update(rx int64, tx int64) {
	w.rxHist.Append(int(rx))
	w.Lines[0].Title = fmt.Sprintf("RX [%s/s]", byteFormat(int64(w.rxHist.Last())))

	w.txHist.Append(int(tx))
	w.Lines[1].Title = fmt.Sprintf("TX [%s/s]", byteFormat(int64(w.txHist.Last())))
}
