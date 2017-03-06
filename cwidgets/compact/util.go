package compact

// Common helper functions

import (
	"fmt"
	ui "github.com/gizak/termui"
)

const colSpacing = 1

// Calculate per-column width, given total width and number of items
func calcWidth(width, items int) int {
	spacing := colSpacing * items
	return (width - statusWidth - spacing) / items
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
