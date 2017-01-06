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
	stats   chan *docker.Stats
	widgets widgets.ContainerWidgets
	reader  *StatReader
}

func NewContainer(c docker.APIContainers) *Container {
	id := c.ID[:12]
	name := strings.Replace(c.Names[0], "/", "", 1) // use primary container name
	return &Container{
		id:      id,
		name:    name,
		done:    make(chan bool),
		stats:   make(chan *docker.Stats),
		widgets: widgets.NewCompact(id, name),
		reader:  &StatReader{},
	}
}

func (c *Container) Expand() {
	c.widgets = widgets.NewExpanded(c.id, c.name)
}

func (c *Container) Collapse() {
	c.widgets = widgets.NewCompact(c.id, c.name)
}

func (c *Container) Collect(client *docker.Client) {
	go func() {
		opts := docker.StatsOptions{
			ID:     c.id,
			Stats:  c.stats,
			Stream: true,
			Done:   c.done,
		}
		client.Stats(opts)
	}()

	go func() {
		for s := range c.stats {
			c.reader.Read(s)
			c.widgets.SetCPU(c.reader.CPUUtil)
			c.widgets.SetMem(c.reader.MemUsage, c.reader.MemLimit, c.reader.MemPercent)
			c.widgets.SetNet(c.reader.NetRx, c.reader.NetTx)
		}
	}()
}
