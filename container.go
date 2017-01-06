package main

import (
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
			c.widgets.cpu.Set(c.reader.CPUUtil)
			c.widgets.SetMem(c.reader.MemUsage, c.reader.MemLimit)
			c.widgets.SetNet(c.reader.NetRx, c.reader.NetTx)
		}
	}()

}
