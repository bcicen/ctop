package compact

import (
	ui "github.com/gizak/termui"
)

var header *CompactHeader

type CompactGrid struct {
	ui.GridBufferer
	Rows   []ui.GridBufferer
	X, Y   int
	Width  int
	Height int
	Offset int // starting row offset
}

func NewCompactGrid() *CompactGrid {
	header = NewCompactHeader() // init column header
	return &CompactGrid{}
}

func (cg *CompactGrid) Align() {
	y := cg.Y
	if cg.Offset >= len(cg.Rows) {
		cg.Offset = 0
	}
	// update row ypos, width recursively
	for _, r := range cg.pageRows() {
		r.SetY(y)
		y += r.GetHeight()
		r.SetWidth(cg.Width)
	}
}

func (cg *CompactGrid) Clear()         { cg.Rows = []ui.GridBufferer{} }
func (cg *CompactGrid) GetHeight() int { return len(cg.Rows) + header.Height }
func (cg *CompactGrid) SetX(x int)     { cg.X = x }
func (cg *CompactGrid) SetY(y int)     { cg.Y = y }
func (cg *CompactGrid) SetWidth(w int) { cg.Width = w }
func (cg *CompactGrid) MaxRows() int   { return ui.TermHeight() - header.Height - cg.Y }

func (cg *CompactGrid) pageRows() (rows []ui.GridBufferer) {
	rows = append(rows, header)
	rows = append(rows, cg.Rows[cg.Offset:]...)
	return rows
}

func (cg *CompactGrid) Buffer() ui.Buffer {
	buf := ui.NewBuffer()
	for _, r := range cg.pageRows() {
		buf.Merge(r.Buffer())
	}
	return buf
}

func (cg *CompactGrid) AddRows(rows ...ui.GridBufferer) {
	for _, r := range rows {
		cg.Rows = append(cg.Rows, r)
	}
}
