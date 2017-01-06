package widgets

import (
	"fmt"
	"strconv"

	ui "github.com/gizak/termui"
)

type CPU struct {
	*ui.Gauge
}

func NewCPU() *CPU {
	return &CPU{mkGauge()}
}

func (c *CPU) Set(val int) {
	c.BarColor = colorScale(val)
	c.Label = fmt.Sprintf("%s%%", strconv.Itoa(val))
	if val < 5 {
		val = 5
		c.BarColor = ui.ColorBlack
	}
	c.Percent = val
}
