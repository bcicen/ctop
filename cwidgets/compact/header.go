package compact

import (
	ui "github.com/gizak/termui"
)

type CompactHeader struct {
	X, Y   int
	Width  int
	Height int
	pars   []*ui.Par
}

func NewCompactHeader() *CompactHeader {
	fields := []string{"", "NAME", "CID", "CPU", "MEM", "NET RX/TX", "IO R/W", "Pids"}
	ch := &CompactHeader{}
	ch.Height = 2
	for _, f := range fields {
		ch.addFieldPar(f)
	}
	return ch
}

func (ch *CompactHeader) GetHeight() int {
	return ch.Height
}

func (ch *CompactHeader) SetWidth(w int) {
	x := ch.X
	autoWidth := calcWidth(w, 7)
	for n, col := range ch.pars {
		// set status column to static width
		if n == 0 {
			col.SetX(x)
			col.SetWidth(statusWidth)
			x += statusWidth
			continue
		}
		col.SetX(x)
		col.SetWidth(autoWidth)
		x += autoWidth + colSpacing
	}
	ch.Width = w
}

func (ch *CompactHeader) SetX(x int) {
	ch.X = x
}

func (ch *CompactHeader) SetY(y int) {
	for _, p := range ch.pars {
		p.SetY(y)
	}
	ch.Y = y
}

func (ch *CompactHeader) Buffer() ui.Buffer {
	buf := ui.NewBuffer()
	for _, p := range ch.pars {
		buf.Merge(p.Buffer())
	}
	return buf
}

func (ch *CompactHeader) addFieldPar(s string) {
	p := ui.NewPar(s)
	p.Height = ch.Height
	p.Border = false
	ch.pars = append(ch.pars, p)
}
