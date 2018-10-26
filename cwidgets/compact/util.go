package compact

// Common helper functions

import (
	"fmt"

	ui "github.com/gizak/termui"
)

const colSpacing = 1

// per-column width. 0 == auto width
var colWidths = []int{
	3, // status
	0, // name
	0, // cid
	0, // cpu
	0, // memory
	0, // net
	0, // io
	4, // pids
}

// Calculate per-column width, given total width
func calcWidth(width int) int {
	spacing := colSpacing * len(colWidths)
	var staticCols int
	for _, w := range colWidths {
		width -= w
		if w == 0 {
			staticCols++
		}
	}
	return (width - spacing) / staticCols
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
