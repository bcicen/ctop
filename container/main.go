package container

import (
	"github.com/bcicen/ctop/connector/collector"
	"github.com/bcicen/ctop/connector/manager"
	"github.com/bcicen/ctop/cwidgets"
	"github.com/bcicen/ctop/cwidgets/compact"
	"github.com/bcicen/ctop/logging"
	"github.com/bcicen/ctop/models"
)

var (
	log = logging.Init()
)

const (
	running = "running"
)

// docker compose project
type Project struct {
	Name    string
	WorkDir string
	Config  string
	Count   int // Containers Count
	Widgets *compact.CompactRow
	Metrics models.Metrics
}

// Metrics and metadata representing a container
type Container struct {
	models.Metrics
	Id        string
	Meta      models.Meta
	Project   *Project
	Widgets   *compact.CompactRow
	Display   bool // display this container in compact view
	updater   cwidgets.WidgetUpdater
	collector collector.Collector
	manager   manager.Manager
}

func New(id string, collector collector.Collector, manager manager.Manager) *Container {
	widgets := compact.NewCompactRow()
	return &Container{
		Metrics:   models.Metrics{},
		Id:        id,
		Meta:      models.NewMeta("id", id[:12]),
		Project:   nil,
		Widgets:   widgets,
		updater:   widgets,
		collector: collector,
		manager:   manager,
	}
}

func NewProject(name string) *Project {
	p := &Project{Name: name}
	// create a compact row for the project
	widgets := compact.NewCompactRow()
	meta := models.NewMeta("name", name)
	widgets.SetMeta(meta)
	p.Widgets = widgets
	return p
}

func (c *Container) RecreateWidgets() {
	c.SetUpdater(cwidgets.NullWidgetUpdater{})
	c.Widgets = compact.NewCompactRow()
	c.SetUpdater(c.Widgets)
}

func (c *Container) SetUpdater(u cwidgets.WidgetUpdater) {
	c.updater = u
	c.updater.SetMeta(c.Meta)
}

func (c *Container) SetMeta(k, v string) {
	c.Meta[k] = v
	c.updater.SetMeta(c.Meta)
}

func (c *Container) GetMeta(k string) string {
	return c.Meta.Get(k)
}

func (c *Container) SetState(s string) {
	c.SetMeta("state", s)
	// start collector, if needed
	if s == running && !c.collector.Running() {
		c.collector.Start()
		c.Read(c.collector.Stream())
	}
	// stop collector, if needed
	if s != running && c.collector.Running() {
		c.collector.Stop()
	}
}

// Logs returns container log collector
func (c *Container) Logs() collector.LogCollector {
	return c.collector.Logs()
}

// Read metric stream, updating widgets
func (c *Container) Read(stream chan models.Metrics) {
	go func() {
		for metrics := range stream {
			oldContainerMetrics := c.Metrics
			c.Project.Metrics.Subtract(oldContainerMetrics)
			c.Project.Metrics.Add(metrics)
			c.Project.Widgets.SetMetrics(c.Project.Metrics)
			c.Metrics = metrics
			c.updater.SetMetrics(metrics)
		}
		log.Infof("reader stopped for container: %s", c.Id)
		c.Project.Metrics.Subtract(c.Metrics)
		c.Metrics = models.Metrics{}
		c.Widgets.Reset()
	}()
	log.Infof("reader started for container: %s", c.Id)
}

func (c *Container) Start() {
	if c.Meta["state"] != running {
		if err := c.manager.Start(); err != nil {
			log.Warningf("container %s: %v", c.Id, err)
			log.StatusErr(err)
			return
		}
		c.SetState(running)
	}
}

func (c *Container) Stop() {
	if c.Meta["state"] == running {
		if err := c.manager.Stop(); err != nil {
			log.Warningf("container %s: %v", c.Id, err)
			log.StatusErr(err)
			return
		}
		c.SetState("exited")
	}
}

func (c *Container) Remove() {
	if err := c.manager.Remove(); err != nil {
		log.Warningf("container %s: %v", c.Id, err)
		log.StatusErr(err)
	}
}

func (c *Container) Pause() {
	if c.Meta["state"] == running {
		if err := c.manager.Pause(); err != nil {
			log.Warningf("container %s: %v", c.Id, err)
			log.StatusErr(err)
			return
		}
		c.SetState("paused")
	}
}

func (c *Container) Unpause() {
	if c.Meta["state"] == "paused" {
		if err := c.manager.Unpause(); err != nil {
			log.Warningf("container %s: %v", c.Id, err)
			log.StatusErr(err)
			return
		}
		c.SetState(running)
	}
}

func (c *Container) Restart() {
	if c.Meta["state"] == running {
		if err := c.manager.Restart(); err != nil {
			log.Warningf("container %s: %v", c.Id, err)
			log.StatusErr(err)
			return
		}
	}
}

func (c *Container) Exec(cmd []string) error {
	return c.manager.Exec(cmd)
}
