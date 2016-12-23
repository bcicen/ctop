package main

import (
	"fmt"
	"math"
	"strconv"

	"github.com/fsouza/go-dockerclient"
	ui "github.com/gizak/termui"
)

type Widgets struct {
	cid    *ui.Par
	cpu    *ui.Gauge
	memory *ui.Gauge
}

func NewWidgets(id string) *Widgets {
	cid := ui.NewPar(id)
	cid.Border = false
	cid.Height = 2
	cid.Width = 10
	cid.TextFgColor = ui.ColorWhite
	return &Widgets{cid, mkGauge(), mkGauge()}
}

type Container struct {
	id      string
	widgets *Widgets
	stats   chan *docker.Stats
	done    chan bool
}

func NewContainer(cid string) *Container {
	return &Container{
		id:      cid,
		widgets: NewWidgets(cid),
		stats:   make(chan *docker.Stats),
		done:    make(chan bool),
	}
}

func (c *Container) Collect(client *docker.Client) {

	go func() {
		fmt.Sprintf("starting collector for container: %s\n", c.id)
		opts := docker.StatsOptions{
			ID:     c.id,
			Stats:  c.stats,
			Stream: true,
			Done:   c.done,
		}
		client.Stats(opts)
		fmt.Sprintf("stopping collector for container: %s\n", c.id)
	}()

	go func() {
		for s := range c.stats {
			c.UpdateMem(s.MemoryStats.Usage, s.MemoryStats.Limit)
		}
	}()

}

func (c *Container) UpdateCPU(total uint64, system uint64) {
	c.widgets.cpu.BarColor = colorScale(n)
	c.widgets.cpu.Percent = n
}

func (c *Container) UpdateMem(cur uint64, limit uint64) {
	percent := round((float64(cur) / float64(limit)) * 100)
	if percent < 5 {
		percent = 5
	}
	c.widgets.memory.Percent = percent
	c.widgets.memory.Label = fmt.Sprintf("%s / %s", byteFormat(cur), byteFormat(limit))
}

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
