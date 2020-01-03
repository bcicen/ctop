package compact

import (
	"github.com/bcicen/ctop/config"
	"github.com/bcicen/ctop/logging"
	"github.com/bcicen/ctop/models"

	ui "github.com/gizak/termui"
)

const rowPadding = 1

var log = logging.Init()

type RowBufferer interface {
	SetY(int)
	SetWidths(int, []int)
	GetHeight() int
	Buffer() ui.Buffer
}

type CompactRow struct {
	Bg     *RowBg
	Cols   []CompactCol
	X, Y   int
	Height int
	widths []int // column widths
}

func NewCompactRow() *CompactRow {
	row := &CompactRow{
		Bg:     NewRowBg(),
		Cols:   newRowWidgets(),
		X:      rowPadding,
		Height: 1,
	}

	return row
}

func (row *CompactRow) SetMeta(m models.Meta) {
	for _, w := range row.Cols {
		w.SetMeta(m)
	}
}

func (row *CompactRow) SetMetrics(m models.Metrics) {
	for _, w := range row.Cols {
		w.SetMetrics(m)
	}
}

// Set gauges, counters, etc. to default unread values
func (row *CompactRow) Reset() {
	for _, w := range row.Cols {
		w.Reset()
	}
}

func (row *CompactRow) GetHeight() int { return row.Height }

//func (row *CompactRow) SetX(x int)     { row.X = x }

func (row *CompactRow) SetY(y int) {
	if y == row.Y {
		return
	}

	row.Bg.Y = y
	for _, w := range row.Cols {
		w.SetY(y)
	}
	row.Y = y
}

func (row *CompactRow) SetWidths(totalWidth int, widths []int) {
	x := row.X

	row.Bg.SetX(x)
	row.Bg.SetWidth(totalWidth)

	for n, w := range row.Cols {
		w.SetX(x)
		w.SetWidth(widths[n])
		x += widths[n] + colSpacing
	}
}

func (row *CompactRow) Buffer() ui.Buffer {
	buf := ui.NewBuffer()
	buf.Merge(row.Bg.Buffer())
	for _, w := range row.Cols {
		buf.Merge(w.Buffer())
	}
	return buf
}

func (row *CompactRow) Highlight() {
	row.Cols[1].Highlight()
	if config.GetSwitchVal("fullRowCursor") {
		for _, w := range row.Cols {
			w.Highlight()
		}
	}
}

func (row *CompactRow) UnHighlight() {
	row.Cols[1].UnHighlight()
	if config.GetSwitchVal("fullRowCursor") {
		for _, w := range row.Cols {
			w.UnHighlight()
		}
	}
}

type RowBg struct {
	*ui.Par
}

func NewRowBg() *RowBg {
	bg := ui.NewPar("")
	bg.Height = 1
	bg.Border = false
	bg.Bg = ui.ThemeAttr("par.text.bg")
	return &RowBg{bg}
}

func (w *RowBg) Highlight()   { w.Bg = ui.ThemeAttr("par.text.fg") }
func (w *RowBg) UnHighlight() { w.Bg = ui.ThemeAttr("par.text.bg") }
