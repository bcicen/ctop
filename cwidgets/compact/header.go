package compact

import (
	ui "github.com/gizak/termui"
)

type CompactHeader struct {
	X, Y   int
	Width  int
	Height int
	cols   []CompactCol
	widths []int
	pars   []*ui.Par
}

func NewCompactHeader() *CompactHeader {
	return &CompactHeader{Height: 2}
}

func (row *CompactHeader) GetHeight() int {
	return row.Height
}

func (row *CompactHeader) SetWidths(totalWidth int, widths []int) {
	x := row.X

	for n, w := range row.pars {
		w.SetX(x)
		w.SetWidth(widths[n])
		x += widths[n] + colSpacing
	}
	row.Width = totalWidth
}

func (row *CompactHeader) SetX(x int) {
	row.X = x
}

func (row *CompactHeader) SetY(y int) {
	for _, p := range row.pars {
		p.SetY(y)
	}
	row.Y = y
}

func (row *CompactHeader) Buffer() ui.Buffer {
	buf := ui.NewBuffer()
	for _, p := range row.pars {
		buf.Merge(p.Buffer())
	}
	return buf
}

func (row *CompactHeader) addFieldPar(s string) {
	p := ui.NewPar(s)
	p.Height = row.Height
	p.Border = false
	row.pars = append(row.pars, p)
}
