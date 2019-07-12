package widgets

import (
	"fmt"
	"time"

	ui "github.com/gizak/termui"
)

type CTopHeader struct {
	Time   *ui.Par
	Count  *ui.Par
	Filter *ui.Par
	bg     *ui.Par
}

func NewCTopHeader() *CTopHeader {
	return &CTopHeader{
		Time:   headerPar(2, ""),
		Count:  headerPar(24, "-"),
		Filter: headerPar(40, ""),
		bg:     headerBg(),
	}
}

func (c *CTopHeader) Buffer() ui.Buffer {
	buf := ui.NewBuffer()
	c.Time.Text = timeStr()
	buf.Merge(c.bg.Buffer())
	buf.Merge(c.Time.Buffer())
	buf.Merge(c.Count.Buffer())
	buf.Merge(c.Filter.Buffer())
	return buf
}

func (c *CTopHeader) Align() {
	c.bg.SetWidth(ui.TermWidth() - 1)
}

func (c *CTopHeader) Height() int {
	return c.bg.Height
}

func headerBgBordered() *ui.Par {
	bg := ui.NewPar("")
	bg.X = 1
	bg.Height = 3
	bg.Bg = ui.ThemeAttr("header.bg")
	return bg
}

func headerBg() *ui.Par {
	bg := ui.NewPar("")
	bg.X = 1
	bg.Height = 1
	bg.Border = false
	bg.Bg = ui.ThemeAttr("header.bg")
	return bg
}

func (c *CTopHeader) SetCount(val int) {
	c.Count.Text = fmt.Sprintf("%d containers", val)
}

func (c *CTopHeader) SetFilter(val string) {
	if val == "" {
		c.Filter.Text = ""
	} else {
		c.Filter.Text = fmt.Sprintf("filter: %s", val)
	}
}

func timeStr() string {
	ts := time.Now().Local().Format("15:04:05 MST")
	return fmt.Sprintf("ctop - %s", ts)
}

func headerPar(x int, s string) *ui.Par {
	p := ui.NewPar(fmt.Sprintf(" %s", s))
	p.X = x
	p.Border = false
	p.Height = 1
	p.Width = 20
	p.Bg = ui.ThemeAttr("header.bg")
	p.TextFgColor = ui.ThemeAttr("header.fg")
	p.TextBgColor = ui.ThemeAttr("header.bg")
	return p
}
