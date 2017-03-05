package compact

import (
	ui "github.com/gizak/termui"
)

type CompactGrid struct {
	ui.GridBufferer
	Rows     []*Compact // rows to render
	X, Y     int
	Width    int
	Height   int
	header   *CompactHeader
	cursorID string
}

func NewCompactGrid() *CompactGrid {
	return &CompactGrid{
		header: NewCompactHeader(),
	}
}

func (cg *CompactGrid) Align() {
	// Update y recursively
	cg.header.SetY(cg.Y)
	y := cg.Y + 1
	for _, r := range cg.Rows {
		r.SetY(y)
		y += r.Height
	}
	// Update width recursively
	cg.header.SetWidth(cg.Width)
	for _, r := range cg.Rows {
		r.SetWidth(cg.Width)
	}
}

func (cg *CompactGrid) Clear()         { cg.Rows = []*Compact{} }
func (cg *CompactGrid) GetHeight() int { return len(cg.Rows) }
func (cg *CompactGrid) SetX(x int)     { cg.X = x }
func (cg *CompactGrid) SetY(y int)     { cg.Y = y }
func (cg *CompactGrid) SetWidth(w int) { cg.Width = w }

func (cg *CompactGrid) Buffer() ui.Buffer {
	buf := ui.NewBuffer()
	buf.Merge(cg.header.Buffer())
	for _, r := range cg.Rows {
		buf.Merge(r.Buffer())
	}
	return buf
}

func (cg *CompactGrid) AddRows(rows ...*Compact) {
	for _, r := range rows {
		cg.Rows = append(cg.Rows, r)
	}
}
