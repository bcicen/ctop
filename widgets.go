package main

import (
	"fmt"
	"math"
	"strconv"

	"github.com/fsouza/go-dockerclient"
	ui "github.com/gizak/termui"
)

type CpuCalc struct {
	lastCpu    uint64
	lastSysCpu uint64
}

func (c *CpuCalc) Utilization(cpu uint64, syscpu uint64, ncpus int) int {
	cpudiff := float64(cpu) - float64(c.lastCpu)
	syscpudiff := float64(syscpu) - float64(c.lastSysCpu)
	util := round((cpudiff / syscpudiff * 100) * float64(ncpus))
	c.lastCpu = cpu
	c.lastSysCpu = syscpu
	return util
}

type Widgets struct {
	cid     *ui.Par
	cpu     *ui.Gauge
	memory  *ui.Gauge
	cpucalc *CpuCalc
}

func NewWidgets(id string) *Widgets {
	cid := ui.NewPar(id)
	cid.Border = false
	cid.Height = 1
	cid.Width = 10
	cid.TextFgColor = ui.ColorWhite
	return &Widgets{cid, mkGauge(), mkGauge()}
}

func mkGauge() *ui.Gauge {
	g := ui.NewGauge()
	g.Height = 1
	g.Border = false
	g.Percent = 0
	g.PaddingBottom = 0
	g.BarColor = ui.ColorGreen
	g.Label = "-"
	return g
}
