package container

import (
	"github.com/bcicen/ctop/connector/collector"
	"github.com/bcicen/ctop/cwidgets"
	"github.com/bcicen/ctop/cwidgets/compact"
	"github.com/bcicen/ctop/logging"
	"github.com/bcicen/ctop/models"
)

var (
	log = logging.Init()
)

// Metrics and metadata representing a container
type Container struct {
	models.Metrics
	Id        string
	Meta      map[string]string
	Widgets   *compact.Compact
	Display   bool // display this container in compact view
	updater   cwidgets.WidgetUpdater
	collector collector.Collector
}

func New(id string, collector collector.Collector) *Container {
	widgets := compact.NewCompact(id)
	return &Container{
		Metrics:   models.NewMetrics(),
		Id:        id,
		Meta:      make(map[string]string),
		Widgets:   widgets,
		updater:   widgets,
		collector: collector,
	}
}

func (c *Container) SetUpdater(u cwidgets.WidgetUpdater) {
	c.updater = u
	for k, v := range c.Meta {
		c.updater.SetMeta(k, v)
	}
}

func (c *Container) SetMeta(k, v string) {
	c.Meta[k] = v
	c.updater.SetMeta(k, v)
}

func (c *Container) GetMeta(k string) string {
	if v, ok := c.Meta[k]; ok {
		return v
	}
	return ""
}

func (c *Container) SetState(s string) {
	c.SetMeta("state", s)
	// start collector, if needed
	if s == "running" && !c.collector.Running() {
		c.collector.Start()
		c.Read(c.collector.Stream())
	}
	// stop collector, if needed
	if s != "running" && c.collector.Running() {
		c.collector.Stop()
	}
}

// Read metric stream, updating widgets
func (c *Container) Read(stream chan models.Metrics) {
	go func() {
		for metrics := range stream {
			c.Metrics = metrics
			c.updater.SetMetrics(metrics)
		}
		log.Infof("reader stopped for container: %s", c.Id)
		c.Metrics = models.NewMetrics()
		c.Widgets.Reset()
	}()
	log.Infof("reader started for container: %s", c.Id)
}
