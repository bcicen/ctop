package widgets

import (
	"fmt"
	"strconv"

	ui "github.com/gizak/termui"
)

const (
	kb = 1024
	mb = kb * 1024
	gb = mb * 1024
)

// convenience method
func byteFormatInt(n int) string {
	return byteFormat(int64(n))
}

func byteFormat(n int64) string {
	if n < kb {
		return fmt.Sprintf("%sB", strconv.FormatInt(n, 10))
	}
	if n < mb {
		n = n / kb
		return fmt.Sprintf("%sK", strconv.FormatInt(n, 10))
	}
	if n < gb {
		n = n / mb
		return fmt.Sprintf("%sM", strconv.FormatInt(n, 10))
	}
	n = n / gb
	return fmt.Sprintf("%sG", strconv.FormatInt(n, 10))
}

func compactPar(s string) *ui.Par {
	p := ui.NewPar(s)
	p.Border = false
	p.Height = 1
	p.Width = 20
	p.TextFgColor = ui.ColorWhite
	return p
}

func colorScale(n int) ui.Attribute {
	if n > 70 {
		return ui.ColorRed
	}
	if n > 30 {
		return ui.ColorYellow
	}
	return ui.ColorGreen
}
