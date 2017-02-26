package compact

import (
	ui "github.com/gizak/termui"
)

type CompactHeader struct {
	pars []*ui.Par
}

func NewCompactHeader() *CompactHeader {
	fields := []string{"", "NAME", "CID", "CPU", "MEM", "NET RX/TX"}
	header := &CompactHeader{}
	for _, f := range fields {
		header.pars = append(header.pars, slimHeaderPar(f))
	}
	return header
}

func (c *CompactHeader) SetWidth(w int) {
	x := 1
	autoWidth := calcWidth(w, 5)
	for n, col := range c.pars {
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
}

func (c *CompactHeader) SetY(y int) {
	for _, p := range c.pars {
		p.SetY(y)
	}
}

func (c *CompactHeader) Buffer() ui.Buffer {
	buf := ui.NewBuffer()
	for _, p := range c.pars {
		buf.Merge(p.Buffer())
	}
	return buf
}
