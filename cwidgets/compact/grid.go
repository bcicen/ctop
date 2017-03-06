package compact

import (
	ui "github.com/gizak/termui"
)

var header = NewCompactHeader()

type CompactGrid struct {
	ui.GridBufferer
	Rows     []ui.GridBufferer
	X, Y     int
	Width    int
	Height   int
	cursorID string
}

func NewCompactGrid() *CompactGrid {
	return &CompactGrid{}
}

func (cg *CompactGrid) Align() {
	// update row y pos recursively
	y := cg.Y
	for _, r := range cg.Rows {
		r.SetY(y)
		y += r.GetHeight()
	}

	// update row width recursively
	for _, r := range cg.Rows {
		r.SetWidth(cg.Width)
	}
}

func (cg *CompactGrid) Clear()         { cg.Rows = []ui.GridBufferer{header} }
func (cg *CompactGrid) GetHeight() int { return len(cg.Rows) }
func (cg *CompactGrid) SetX(x int)     { cg.X = x }
func (cg *CompactGrid) SetY(y int)     { cg.Y = y }
func (cg *CompactGrid) SetWidth(w int) { cg.Width = w }

func (cg *CompactGrid) Buffer() ui.Buffer {
	buf := ui.NewBuffer()
	for _, r := range cg.Rows {
		buf.Merge(r.Buffer())
	}
	return buf
}

func (cg *CompactGrid) AddRows(rows ...ui.GridBufferer) {
	for _, r := range rows {
		cg.Rows = append(cg.Rows, r)
	}
}
