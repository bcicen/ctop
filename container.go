package main

import (
	"github.com/bcicen/ctop/collector"
	"github.com/bcicen/ctop/widgets"
)

type Container struct {
	id      string
	name    string
	state   string
	metrics collector.Metrics
	collect collector.Collector
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
	// start collector if necessary
	if s == "running" && !c.collect.Running() {
		c.Collect()
	}
}

func (c *Container) Collect() {
	log.Infof("starting collector for container: %s", c.id)
	c.collect.Start()
	go func() {
		for metrics := range c.collect.Stream() {
			c.metrics = metrics
			c.widgets.SetCPU(metrics.CPUUtil)
			c.widgets.SetMem(metrics.MemUsage, metrics.MemLimit, metrics.MemPercent)
			c.widgets.SetNet(metrics.NetRx, metrics.NetTx)
		}
	}()
}
