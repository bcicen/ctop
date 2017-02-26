package cwidgets

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
func ByteFormatInt(n int) string {
	return ByteFormat(int64(n))
}

func ByteFormat(n int64) string {
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
	nf := float64(n) / gb
	return fmt.Sprintf("%sG", unpadFloat(nf))
}

func unpadFloat(f float64) string {
	return strconv.FormatFloat(f, 'f', getPrecision(f), 64)
}

func getPrecision(f float64) int {
	frac := int((f - float64(int(f))) * 100)
	if frac == 0 {
		return 0
	}
	if frac%10 == 0 {
		return 1
	}
	return 2 // default precision
}

func ColorScale(n int) ui.Attribute {
	if n > 70 {
		return ui.ColorRed
	}
	if n > 30 {
		return ui.ColorYellow
	}
	return ui.ColorGreen
}
