package entity

import (
	"github.com/bcicen/ctop/connector/collector"
	"github.com/bcicen/ctop/models"
)

// Metrics and metadata representing a container
type Container struct {
	models.Metrics
	Meta
	Id        string
	collector collector.Collector
}

func NewContainer(id string, collector collector.Collector) *Container {

	return &Container{
		Metrics:   models.NewMetrics(),
		Meta:      NewMeta(id),
		Id:        id,
		collector: collector,
	}
}

func (c *Container) SetState(val string) {
	c.Meta.SetMeta("state", val)
	// start collector, if needed
	if val == "running" && !c.collector.Running() {
		c.collector.Start()
		c.Read(c.collector.Stream())
	}
	// stop collector, if needed
	if val != "running" && c.collector.Running() {
		c.collector.Stop()
	}
}

// Return container log collector
func (c *Container) Logs() collector.LogCollector {
	return c.collector.Logs()
}

// Read metric stream, updating widgets
func (c *Container) Read(stream chan models.Metrics) {
	go func() {
		for metrics := range stream {
			c.Metrics = metrics
			c.Meta.updater.SetMetrics(metrics)
		}
		log.Infof("reader stopped for container: %s", c.Id)
		c.Metrics = models.NewMetrics()
		c.Meta.Widgets.Reset()
	}()
	log.Infof("reader started for container: %s", c.Id)
}
