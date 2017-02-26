package compact

// Common helper functions

import (
	"fmt"
	ui "github.com/gizak/termui"
)

// Calculate per-column width, given total width and number of items
func calcWidth(width, items int) int {
	spacing := colSpacing * items
	return (width - statusWidth - spacing) / items
}

func slimHeaderPar(s string) *ui.Par {
	p := slimPar(s)
	p.Y = 2
	p.Height = 2
	return p
}

func slimPar(s string) *ui.Par {
	p := ui.NewPar(s)
	p.Border = false
	p.Height = 1
	p.Width = 20
	p.TextFgColor = ui.ColorWhite
	return p
}

func slimGauge() *ui.Gauge {
	g := ui.NewGauge()
	g.Height = 1
	g.Border = false
	g.Percent = 0
	g.PaddingBottom = 0
	g.BarColor = ui.ColorGreen
	g.Label = "-"
	return g
}

func centerParText(p *ui.Par) {
	var text string
	var padding string

	// strip existing left-padding
	for i, ch := range p.Text {
		if string(ch) != " " {
			text = p.Text[i:]
			break
		}
	}

	padlen := (p.InnerWidth() - len(text)) / 2
	for i := 0; i < padlen; i++ {
		padding += " "
	}
	p.Text = fmt.Sprintf("%s%s", padding, text)
}
