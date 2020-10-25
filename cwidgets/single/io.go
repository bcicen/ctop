package single

import (
	"fmt"
	"strings"

	"github.com/bcicen/ctop/cwidgets"
	ui "github.com/gizak/termui"
)

type IO struct {
	*ui.Sparklines
	readHist  *DiffHist
	writeHist *DiffHist
}

func NewIO() *IO {
	io := &IO{ui.NewSparklines(), NewDiffHist(60), NewDiffHist(60)}
	io.BorderLabel = "IO"
	io.Height = 6
	io.Width = colWidth[0]
	io.X = 0
	io.Y = 24

	read := ui.NewSparkline()
	read.Title = "READ"
	read.Height = 1
	read.Data = io.readHist.Data
	read.LineColor = ui.ColorGreen

	write := ui.NewSparkline()
	write.Title = "WRITE"
	write.Height = 1
	write.Data = io.writeHist.Data
	write.LineColor = ui.ColorYellow

	io.Lines = []ui.Sparkline{read, write}
	return io
}

func (w *IO) Update(read int64, write int64) {
	var rate string

	w.readHist.Append(int(read))
	rate = strings.ToLower(cwidgets.ByteFormatShort(w.readHist.Val))
	w.Lines[0].Title = fmt.Sprintf("read [%s/s]", rate)

	w.writeHist.Append(int(write))
	rate = strings.ToLower(cwidgets.ByteFormatShort(w.writeHist.Val))
	w.Lines[1].Title = fmt.Sprintf("write [%s/s]", rate)
}
