package main

import (
	"fmt"
	"math"
	"strconv"

	"github.com/fsouza/go-dockerclient"
	ui "github.com/gizak/termui"
)

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
		cpucalc: &CpuCalc{},
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
			c.UpdateCPU(s.CPUStats.CPUUsage.TotalUsage, s.CPUStats.SystemCPUUsage, len(s.CPUStats.CPUUsage.PercpuUsage))
		}
	}()

}

func (c *Container) UpdateCPU(total uint64, system uint64, ncpus int) {
	util := c.widgets.cpucalc.Utilization(total, system, ncpus)
	c.widgets.cpu.Label = fmt.Sprintf("%s%%", strconv.Itoa(util))
	c.widgets.cpu.BarColor = colorScale(util)
	if util < 5 && util > 0 {
		util = 5
	}
	c.widgets.cpu.Percent = util
}

func (c *Container) UpdateMem(cur uint64, limit uint64) {
	percent := round((float64(cur) / float64(limit)) * 100)
	if percent < 5 {
		percent = 5
	}
	c.widgets.memory.Percent = percent
	c.widgets.memory.Label = fmt.Sprintf("%s / %s", byteFormat(cur), byteFormat(limit))
}
