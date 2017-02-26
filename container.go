package main

import (
	"strings"

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

func NewContainer(id, name string) *Container {
	c := &Container{
		id:   id,
		name: name,
	}
	c.widgets = widgets.NewCompact(c.ShortID(), c.ShortName(), c.state)
	return c
}

func (c *Container) ShortID() string {
	return c.id[:12]
}

func (c *Container) ShortName() string {
	return strings.Replace(c.name, "/", "", 1) // use primary container name
}

func (c *Container) Expand() {
	var curWidgets widgets.ContainerWidgets

	curWidgets = c.widgets
	c.widgets = widgets.NewExpanded(c.ShortID(), c.name)
	c.widgets.Render(0, 0)
	c.widgets = curWidgets
}

func (c *Container) SetState(s string) {
	c.state = s
	c.widgets.SetStatus(s)
}

// Set metrics to zero state, clear widget gauges
func (c *Container) reset() {
	c.metrics = metrics.Metrics{}
	c.widgets.Reset()
}

// Read metric stream, updating widgets
func (c *Container) Read(stream chan metrics.Metrics) {
	go func() {
		for metrics := range stream {
			c.metrics = metrics
			c.widgets.SetCPU(metrics.CPUUtil)
			c.widgets.SetMem(metrics.MemUsage, metrics.MemLimit, metrics.MemPercent)
			c.widgets.SetNet(metrics.NetRx, metrics.NetTx)
		}
		log.Infof("reader stopped for container: %s", c.id)
		c.reset()
	}()
	log.Infof("reader started for container: %s", c.id)
}
