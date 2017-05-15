package metrics

import (
	api "github.com/fsouza/go-dockerclient"
)

// Docker collector
type Docker struct {
	Metrics
	id         string
	client     *api.Client
	running    bool
	stream     chan Metrics
	done       chan bool
	lastCpu    float64
	lastSysCpu float64
}

func NewDocker(client *api.Client, id string) *Docker {
	return &Docker{
		Metrics: Metrics{},
		id:      id,
		client:  client,
	}
}

func (c *Docker) Start() {
	c.done = make(chan bool)
	c.stream = make(chan Metrics)
	stats := make(chan *api.Stats)

	go func() {
		opts := api.StatsOptions{
			ID:     c.id,
			Stats:  stats,
			Stream: true,
			Done:   c.done,
		}
		c.client.Stats(opts)
		c.running = false
	}()

	go func() {
		defer close(c.stream)
		for s := range stats {
			c.ReadCPU(s)
			c.ReadMem(s)
			c.ReadNet(s)
			c.ReadIO(s)
			c.stream <- c.Metrics
		}
		log.Infof("collector stopped for container: %s", c.id)
	}()

	c.running = true
	log.Infof("collector started for container: %s", c.id)
}

func (c *Docker) Running() bool {
	return c.running
}

func (c *Docker) Stream() chan Metrics {
	return c.stream
}

// Stop collector
func (c *Docker) Stop() {
	c.done <- true
}

func (c *Docker) ReadCPU(stats *api.Stats) {
	ncpus := float64(len(stats.CPUStats.CPUUsage.PercpuUsage))
	total := float64(stats.CPUStats.CPUUsage.TotalUsage)
	system := float64(stats.CPUStats.SystemCPUUsage)

	cpudiff := total - c.lastCpu
	syscpudiff := system - c.lastSysCpu

	c.CPUUtil = round((cpudiff / syscpudiff * 100) * ncpus)
	c.lastCpu = total
	c.lastSysCpu = system
	c.Pids = int(stats.PidsStats.Current)
}

func (c *Docker) ReadMem(stats *api.Stats) {
	c.MemUsage = int64(stats.MemoryStats.Usage)
	c.MemLimit = int64(stats.MemoryStats.Limit)
	c.MemPercent = round((float64(c.MemUsage) / float64(c.MemLimit)) * 100)
}

func (c *Docker) ReadNet(stats *api.Stats) {
	var rx, tx int64
	for _, network := range stats.Networks {
		rx += int64(network.RxBytes)
		tx += int64(network.TxBytes)
	}
	c.NetRx, c.NetTx = rx, tx
}

func (c *Docker) ReadIO(stats *api.Stats) {
	var read, write int64
	for _, blk := range stats.BlkioStats.IOServiceBytesRecursive {
		if blk.Op == "Read" {
			read = int64(blk.Value)
		}
		if blk.Op == "Write" {
			write = int64(blk.Value)
		}
	}
	c.IOBytesRead, c.IOBytesWrite = read, write
}
