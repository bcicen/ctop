package main

import (
	"math"

	ui "github.com/gizak/termui"
)

func byteFormat(n uint64) string {
	if n < 1024 {
		return fmt.Sprintf("%sB", strconv.FormatUint(n, 10))
	}
	if n < 1048576 {
		n = n / 1024
		return fmt.Sprintf("%sK", strconv.FormatUint(n, 10))
	}
	if n < 1073741824 {
		n = n / 1048576
		return fmt.Sprintf("%sM", strconv.FormatUint(n, 10))
	}
	n = n / 1024000000
	return fmt.Sprintf("%sG", strconv.FormatUint(n, 10))
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
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
