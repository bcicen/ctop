package collector

import (
	"context"
	"encoding/json"
	"github.com/bcicen/ctop/models"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"io"
)

// Docker collector
type Docker struct {
	models.Metrics
	id         string
	client     *client.Client
	running    bool
	stream     chan models.Metrics
	done       chan bool
	lastCpu    float64
	lastSysCpu float64
}

func NewDocker(client *client.Client, id string) *Docker {
	return &Docker{
		Metrics: models.Metrics{},
		id:      id,
		client:  client,
	}
}

func (c *Docker) Start() {
	c.done = make(chan bool)
	c.stream = make(chan models.Metrics)
	stats := make(chan *types.StatsJSON)

	go func() {
		ctx := context.Background()
		ss, err := c.client.ContainerStats(ctx, c.id, true)
		if err == nil {
			c.running = false
		}
		decoder := json.NewDecoder(ss.Body)
		cStats := new(types.StatsJSON)

		for err := decoder.Decode(cStats); err != io.EOF; err = decoder.Decode(cStats) {
			if err != nil {
				break
			}
			stats <- cStats
			cStats = new(types.StatsJSON)
		}
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

func (c *Docker) Stream() chan models.Metrics {
	return c.stream
}

func (c *Docker) Logs() LogCollector {
	return NewDockerLogs(c.id, c.client)
}

// Stop collector
func (c *Docker) Stop() {
	c.running = false
	c.done <- true
}

func (c *Docker) ReadCPU(stats *types.StatsJSON) {
	ncpus := uint8(len(stats.CPUStats.CPUUsage.PercpuUsage))
	total := float64(stats.CPUStats.CPUUsage.TotalUsage)
	system := float64(stats.CPUStats.SystemUsage)

	cpudiff := total - c.lastCpu
	syscpudiff := system - c.lastSysCpu

	c.NCpus = ncpus
	c.CPUUtil = percent(cpudiff, syscpudiff)
	c.lastCpu = total
	c.lastSysCpu = system
	c.Pids = int(stats.PidsStats.Current)
}

func (c *Docker) ReadMem(stats *types.StatsJSON) {
	c.MemUsage = int64(stats.MemoryStats.Usage - stats.MemoryStats.Stats["cache"])
	c.MemLimit = int64(stats.MemoryStats.Limit)
	c.MemPercent = percent(float64(c.MemUsage), float64(c.MemLimit))
}

func (c *Docker) ReadNet(stats *types.StatsJSON) {
	var rx, tx int64
	for _, network := range stats.Networks {
		rx += int64(network.RxBytes)
		tx += int64(network.TxBytes)
	}
	c.NetRx, c.NetTx = rx, tx
}

func (c *Docker) ReadIO(stats *types.StatsJSON) {
	var read, write int64
	for _, blk := range stats.BlkioStats.IoServiceBytesRecursive {
		if blk.Op == "Read" {
			read += int64(blk.Value)
		}
		if blk.Op == "Write" {
			write += int64(blk.Value)
		}
	}
	c.IOBytesRead, c.IOBytesWrite = read, write
}
