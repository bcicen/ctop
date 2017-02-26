package compact

import (
	"github.com/bcicen/ctop/cwidgets"
	ui "github.com/gizak/termui"
)

type CompactGrid struct {
	ui.GridBufferer
	Rows     []cwidgets.ContainerWidgets
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

func (c *CompactGrid) Align() {
	// Update y recursively
	c.header.SetY(c.Y)
	y := c.Y + 1
	for n, r := range c.Rows {
		r.SetY(y + n)
	}
	// Update width recursively
	c.header.SetWidth(c.Width)
	for _, r := range c.Rows {
		r.SetWidth(c.Width)
	}
}

func (c *CompactGrid) Clear()         { c.Rows = []cwidgets.ContainerWidgets{} }
func (c *CompactGrid) GetHeight() int { return len(c.Rows) }
func (c *CompactGrid) SetX(x int)     { c.X = x }
func (c *CompactGrid) SetY(y int)     { c.Y = y }
func (c *CompactGrid) SetWidth(w int) { c.Width = w }

func (c *CompactGrid) Buffer() ui.Buffer {
	buf := ui.NewBuffer()
	buf.Merge(c.header.Buffer())
	for _, r := range c.Rows {
		buf.Merge(r.Buffer())
	}
	return buf
}
