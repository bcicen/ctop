package collector

import (
	"github.com/bcicen/ctop/models"
	api "github.com/fsouza/go-dockerclient"
	"github.com/docker/docker/client"
	"context"
	"github.com/bcicen/ctop/config"
	swarmnet "github.com/bcicen/ctop/network"
	"github.com/docker/docker/api/types"
	"io/ioutil"
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

type bc struct {
	stats types.ContainerStats
	err   error
}

func (c *Docker) Start(id string) {
	if config.GetSwitchVal("swarmMode") {
		return
	}
	c.done = make(chan bool)
	c.stream = make(chan models.Metrics)
	stats := make(chan bc)
	go func() {
		ctx, cl := context.WithCancel(context.Background())
		resp, err := c.client.ContainerStats(ctx, id, false)
		stats <- bc{resp, err}
		defer cl()
		defer func() { c.running = false }()
	}()

	go func() {
		defer close(c.stream)
		for s := range stats {
			//c.ReadCPU(s)
			//c.ReadMem(s)
			//c.ReadNet(s)
			//c.ReadIO(s)
			if s.err != nil {
				continue
				log.Errorf("%s", s.err)
			}
			b, _ := ioutil.ReadAll(s.stats.Body)
			log.Infof("%s", string(b))
			s.stats.Body.Close()
			swarmnet.TestDockerNetwork(&c.Metrics)
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
	return &DockerLogs{c.id, c.client, make(chan bool)}
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
	c.MemUsage = int64(stats.MemoryStats.Usage - stats.MemoryStats.Stats.Cache)
	c.MemLimit = int64(stats.MemoryStats.Limit)
	c.MemPercent = percent(float64(c.MemUsage), float64(c.MemLimit))
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
