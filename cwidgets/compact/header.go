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
	fields := []string{"", "NAME", "CID", "CPU", "MEM", "NET RX/TX"}
	header := &CompactHeader{}
	for _, f := range fields {
		header.pars = append(header.pars, headerPar(f))
	}
	return header
}

func (c *CompactHeader) SetWidth(w int) {
	x := 1
	autoWidth := calcWidth(w, 5)
	for n, col := range c.pars {
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
	c.Width = w
}

func (c *CompactHeader) SetY(y int) {
	for _, p := range c.pars {
		p.SetY(y)
	}
	c.Y = y
}

func (c *CompactHeader) Buffer() ui.Buffer {
	buf := ui.NewBuffer()
	for _, p := range c.pars {
		buf.Merge(p.Buffer())
	}
	return buf
}

func headerPar(s string) *ui.Par {
	p := ui.NewPar(s)
	p.Y = 2
	p.Height = 2
	p.Width = 20
	p.Border = false
	return p
}
