package expanded

import (
	"fmt"
	"strings"

	"github.com/bcicen/ctop/cwidgets"
	ui "github.com/gizak/termui"
)

type ExpandedNet struct {
	*ui.Sparklines
	rxHist *DiffHist
	txHist *DiffHist
}

func NewExpandedNet() *ExpandedNet {
	net := &ExpandedNet{ui.NewSparklines(), NewDiffHist(50), NewDiffHist(50)}
	net.BorderLabel = "NET"
	net.Height = 6
	net.Width = colWidth[0]
	net.X = 0
	net.Y = 24

	rx := ui.NewSparkline()
	rx.Title = "RX"
	rx.Height = 1
	rx.Data = net.rxHist.Data
	rx.TitleColor = ui.ColorDefault
	rx.LineColor = ui.ColorGreen

	tx := ui.NewSparkline()
	tx.Title = "TX"
	tx.Height = 1
	tx.Data = net.txHist.Data
	tx.TitleColor = ui.ColorDefault
	tx.LineColor = ui.ColorYellow

	net.Lines = []ui.Sparkline{rx, tx}
	return net
}

func (w *ExpandedNet) Update(rx int64, tx int64) {
	var rate string

	w.rxHist.Append(int(rx))
	rate = strings.ToLower(cwidgets.ByteFormatInt(w.rxHist.Val))
	w.Lines[0].Title = fmt.Sprintf("RX [%s/s]", rate)

	w.txHist.Append(int(tx))
	rate = strings.ToLower(cwidgets.ByteFormatInt(w.txHist.Val))
	w.Lines[1].Title = fmt.Sprintf("TX [%s/s]", rate)
}
