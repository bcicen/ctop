package main

import (
	"strings"

	"github.com/bcicen/ctop/widgets"
	"github.com/fsouza/go-dockerclient"
)

type Container struct {
	id      string
	name    string
	done    chan bool
	widgets widgets.ContainerWidgets
	reader  *MetricsReader
}

func NewContainer(c docker.APIContainers) *Container {
	id := c.ID[:12]
	name := strings.Replace(c.Names[0], "/", "", 1) // use primary container name
	return &Container{
		id:      id,
		name:    name,
		done:    make(chan bool),
		widgets: widgets.NewCompact(id, name),
		reader:  NewMetricsReader(),
	}
}

func (c *Container) Expand() {
	c.widgets = widgets.NewExpanded(c.id, c.name)
}

func (c *Container) Collapse() {
	c.widgets = widgets.NewCompact(c.id, c.name)
}

func (c *Container) Collect(client *docker.Client) {
	stats := make(chan *docker.Stats)

	go func() {
		opts := docker.StatsOptions{
			ID:     c.id,
			Stats:  stats,
			Stream: true,
			Done:   c.done,
		}
		client.Stats(opts)
	}()

	go func() {
		for metrics := range c.reader.Read(stats) {
			c.widgets.SetCPU(metrics.CPUUtil)
			c.widgets.SetMem(metrics.MemUsage, metrics.MemLimit, metrics.MemPercent)
			c.widgets.SetNet(metrics.NetRx, metrics.NetTx)
		}
	}()
}
