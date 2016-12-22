package main

import (
	ui "github.com/gizak/termui"
)

type Container struct {
	cid    *ui.Par
	cpu    *ui.Gauge
	memory *ui.Gauge
}

func (c *Container) UpdateCPU(n int) {
	c.cpu.BarColor = colorScale(n)
	c.cpu.Percent = n
}

func (c *Container) UpdateMem(n int) {
	c.memory.Percent = n
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
