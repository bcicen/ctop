package single

import (
	"fmt"
	"strings"

	"github.com/bcicen/ctop/cwidgets"
	ui "github.com/gizak/termui"
)

type Net struct {
	*ui.Sparklines
	rxHist *DiffHist
	txHist *DiffHist
}

func NewNet() *Net {
	net := &Net{ui.NewSparklines(), NewDiffHist(60), NewDiffHist(60)}
	net.BorderLabel = "NET"
	net.Height = 6
	net.Width = colWidth[0]
	net.X = 0
	net.Y = 24

	rx := ui.NewSparkline()
	rx.Title = "RX"
	rx.Height = 1
	rx.Data = net.rxHist.Data
	rx.LineColor = ui.ColorGreen

	tx := ui.NewSparkline()
	tx.Title = "TX"
	tx.Height = 1
	tx.Data = net.txHist.Data
	tx.LineColor = ui.ColorYellow

	net.Lines = []ui.Sparkline{rx, tx}
	return net
}

func (w *Net) Update(rx int64, tx int64) {
	var rate string

	w.rxHist.Append(int(rx))
	rate = strings.ToLower(cwidgets.ByteFormatInt(w.rxHist.Val))
	w.Lines[0].Title = fmt.Sprintf("RX [%s/s]", rate)

	w.txHist.Append(int(tx))
	rate = strings.ToLower(cwidgets.ByteFormatInt(w.txHist.Val))
	w.Lines[1].Title = fmt.Sprintf("TX [%s/s]", rate)
}
