package main

import (
	"github.com/bcicen/ctop/metrics"
	"github.com/bcicen/ctop/widgets"
)

type Container struct {
	id      string
	name    string
	state   string
	metrics metrics.Metrics
	widgets widgets.ContainerWidgets
}

func (c *Container) Expand() {
	c.widgets = widgets.NewExpanded(c.id, c.name)
}

func (c *Container) Collapse() {
	c.widgets = widgets.NewCompact(c.id, c.name)
}

func (c *Container) SetState(s string) {
	c.state = s
	c.widgets.SetStatus(s)
}

// Read metric stream, updating widgets
func (c *Container) Read(stream chan metrics.Metrics) {
	log.Infof("starting reader for container: %s", c.id)
	go func() {
		for metrics := range stream {
			c.metrics = metrics
			c.widgets.SetCPU(metrics.CPUUtil)
			c.widgets.SetMem(metrics.MemUsage, metrics.MemLimit, metrics.MemPercent)
			c.widgets.SetNet(metrics.NetRx, metrics.NetTx)
		}
	}()
}
