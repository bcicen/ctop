package main

import (
	"github.com/bcicen/ctop/cwidgets/compact"
	"github.com/bcicen/ctop/metrics"
)

// Metrics and metadata representing a container
type Container struct {
	metrics.Metrics
	Id        string
	Name      string
	State     string
	Meta      map[string]string
	Updates   chan [2]string
	Widgets   *compact.Compact
	collector metrics.Collector
}

func NewContainer(id string, collector metrics.Collector) *Container {
	return &Container{
		Metrics:   metrics.NewMetrics(),
		Id:        id,
		Meta:      make(map[string]string),
		Updates:   make(chan [2]string),
		Widgets:   compact.NewCompact(id),
		collector: collector,
	}
}

func (c *Container) GetMeta(k string) string {
	if v, ok := c.Meta[k]; ok {
		return v
	}
	return ""
}

func (c *Container) SetMeta(k, v string) {
	c.Meta[k] = v
	c.Updates <- [2]string{k, v}
}

func (c *Container) SetName(n string) {
	c.Name = n
	c.Widgets.Name.Set(n)
}

func (c *Container) SetState(s string) {
	c.State = s
	c.Widgets.Status.Set(s)
	// start collector, if needed
	if c.State == "running" && !c.collector.Running() {
		c.collector.Start()
		c.Read(c.collector.Stream())
	}
	// stop collector, if needed
	if c.State != "running" && c.collector.Running() {
		c.collector.Stop()
	}
}

// Read metric stream, updating widgets
func (c *Container) Read(stream chan metrics.Metrics) {
	go func() {
		for metrics := range stream {
			c.Metrics = metrics
			c.Widgets.SetMetrics(metrics)
		}
		log.Infof("reader stopped for container: %s", c.Id)
		c.Widgets.Reset()
	}()
	log.Infof("reader started for container: %s", c.Id)
}
