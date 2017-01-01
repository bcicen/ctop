package main

import (
	"fmt"

	"github.com/fsouza/go-dockerclient"
)

type Container struct {
	id      string
	name    string
	done    chan bool
	stats   chan *docker.Stats
	widgets *Widgets
	reader  *StatReader
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
			c.reader.Read(s)
			c.widgets.SetCPU(c.reader.CPUUtil)
			c.widgets.SetMem(c.reader.MemUsage, c.reader.MemLimit)
			c.widgets.SetNet(c.reader.NetRx, c.reader.NetTx)
		}
	}()

}
