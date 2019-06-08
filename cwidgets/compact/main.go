package compact

import (
	"github.com/bcicen/ctop/config"
	"github.com/bcicen/ctop/logging"
	"github.com/bcicen/ctop/models"
	ui "github.com/gizak/termui"
)

var log = logging.Init()

type CompactCol interface {
	ui.GridBufferer
	Reset()
	Highlight()
	UnHighlight()
	SetMeta(models.Meta)
	SetMetrics(models.Metrics)
}

type Compact struct {
	Bg     *RowBg
	Cols   []CompactCol
	X, Y   int
	Width  int
	Height int
}

func NewCompact(id string) *Compact {
	// truncate container id
	if len(id) > 12 {
		id = id[:12]
	}
	row := &Compact{
		Bg: NewRowBg(),
		Cols: []CompactCol{
			NewStatus(),
			&NameCol{NewTextCol("-")},
			&CIDCol{NewTextCol(id)},
			&CPUCol{NewGaugeCol()},
			&MemCol{NewGaugeCol()},
			&NetCol{NewTextCol("-")},
			&IOCol{NewTextCol("-")},
			&PIDCol{NewTextCol("-")},
		},
		X:      1,
		Height: 1,
	}
	return row
}

func (row *Compact) SetMeta(m models.Meta) {
	for _, w := range row.Cols {
		w.SetMeta(m)
	}
}

func (row *Compact) SetMetrics(m models.Metrics) {
	for _, w := range row.Cols {
		w.SetMetrics(m)
	}
}

// Set gauges, counters, etc. to default unread values
func (row *Compact) Reset() {
	for _, w := range row.Cols {
		w.Reset()
	}
}

func (row *Compact) GetHeight() int { return row.Height }
func (row *Compact) SetX(x int)     { row.X = x }

func (row *Compact) SetY(y int) {
	if y == row.Y {
		return
	}

	row.Bg.Y = y
	for _, w := range row.Cols {
		w.SetY(y)
	}
	row.Y = y
}

func (row *Compact) SetWidth(width int) {
	if width == row.Width {
		return
	}
	x := row.X

	row.Bg.SetX(x + colWidths[0] + 1)
	row.Bg.SetWidth(width)

	autoWidth := calcWidth(width)
	for n, w := range row.Cols {
		// set static width, if provided
		if colWidths[n] != 0 {
			w.SetX(x)
			w.SetWidth(colWidths[n])
			x += colWidths[n]
			continue
		}
		// else use auto width
		w.SetX(x)
		w.SetWidth(autoWidth)
		x += autoWidth + colSpacing
	}
	row.Width = width
}

func (row *Compact) Buffer() ui.Buffer {
	buf := ui.NewBuffer()
	buf.Merge(row.Bg.Buffer())
	for _, w := range row.Cols {
		buf.Merge(w.Buffer())
	}
	return buf
}

func (row *Compact) Highlight() {
	row.Cols[1].Highlight()
	if config.GetSwitchVal("fullRowCursor") {
		for _, w := range row.Cols {
			w.Highlight()
		}
	}
}

func (row *Compact) UnHighlight() {
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
