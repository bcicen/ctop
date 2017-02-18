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
		Time:   headerPar(2, timeStr()),
		Count:  headerPar(22, "-"),
		Filter: headerPar(42, ""),
		bg:     headerBg(),
	}
}

func (c *CTopHeader) Render() {
	c.Time.Text = timeStr()
	ui.Render(c.bg)
	ui.Render(c.Time, c.Count, c.Filter)
}

func (c *CTopHeader) Height() int {
	return c.bg.Height
}

func headerBgBordered() *ui.Par {
	bg := ui.NewPar("")
	bg.X = 1
	bg.Width = ui.TermWidth() - 1
	bg.Height = 3
	bg.Bg = ui.ColorWhite
	return bg
}

func headerBg() *ui.Par {
	bg := ui.NewPar("")
	bg.X = 1
	bg.Width = ui.TermWidth() - 1
	bg.Height = 1
	bg.Border = false
	bg.Bg = ui.ColorWhite
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
	return time.Now().Local().Format("15:04:05 MST")
}

func headerPar(x int, s string) *ui.Par {
	p := ui.NewPar(fmt.Sprintf(" %s", s))
	p.X = x
	p.Border = false
	p.Height = 1
	p.Width = 20
	p.TextFgColor = ui.ColorDefault
	p.TextBgColor = ui.ColorWhite
	p.Bg = ui.ColorWhite
	return p
}
