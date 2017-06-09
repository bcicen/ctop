package metrics

import (
	"time"

	"github.com/opencontainers/runc/libcontainer"
	"github.com/opencontainers/runc/libcontainer/cgroups"
)

// Runc collector
type Runc struct {
	Metrics
	id         string
	libc       libcontainer.Container
	stream     chan Metrics
	done       bool
	running    bool
	interval   int // collection interval, in seconds
	lastCpu    float64
	lastSysCpu float64
}

func NewRunc(libc libcontainer.Container) *Runc {
	c := &Runc{
		Metrics:  Metrics{},
		id:       libc.ID(),
		libc:     libc,
		interval: 1,
	}
	return c
}

func (c *Runc) Running() bool {
	return c.running
}

func (c *Runc) Start() {
	c.done = false
	c.stream = make(chan Metrics)
	go c.run()
}

func (c *Runc) Stop() {
	c.done = true
}

func (c *Runc) Stream() chan Metrics {
	return c.stream
}

func (c *Runc) run() {
	c.running = true
	defer close(c.stream)

	for {
		stats, err := c.libc.Stats()
		if err != nil {
			log.Errorf("failed to collect stats for container %s:\n%s", c.id, err)
			break
		}

		c.ReadCPU(stats.CgroupStats)
		c.ReadMem(stats.CgroupStats)
		c.ReadNet(stats.Interfaces)

		c.stream <- c.Metrics
		if c.done {
			break
		}
		time.Sleep(1 * time.Second)
	}

	c.running = false
}

func (c *Runc) ReadCPU(stats *cgroups.Stats) {
	u := stats.CpuStats.CpuUsage
	ncpus := float64(len(u.PercpuUsage))
	total := float64(u.TotalUsage)
	system := float64(getSysCPUUsage())

	cpudiff := total - c.lastCpu
	syscpudiff := system - c.lastSysCpu

	c.CPUUtil = round((cpudiff / syscpudiff * 100) * ncpus)
	c.lastCpu = total
	c.lastSysCpu = system
	c.Pids = int(stats.PidsStats.Current)
}

func (c *Runc) ReadMem(stats *cgroups.Stats) {
	c.MemUsage = int64(stats.MemoryStats.Usage.Usage)
	c.MemLimit = int64(stats.MemoryStats.Usage.Limit)
	if c.MemLimit > sysMemTotal && sysMemTotal > 0 {
		c.MemLimit = sysMemTotal
	}
	c.MemPercent = round((float64(c.MemUsage) / float64(c.MemLimit)) * 100)
}

func (c *Runc) ReadNet(interfaces []*libcontainer.NetworkInterface) {
	var rx, tx int64
	for _, network := range interfaces {
		rx += int64(network.RxBytes)
		tx += int64(network.TxBytes)
	}
	c.NetRx, c.NetTx = rx, tx
}

func (c *Runc) ReadIO(stats *cgroups.Stats) {
	var read, write int64
	for _, blk := range stats.BlkioStats.IoServiceBytesRecursive {
		if blk.Op == "Read" {
			read = int64(blk.Value)
		}
		if blk.Op == "Write" {
			write = int64(blk.Value)
		}
	}
	c.IOBytesRead, c.IOBytesWrite = read, write
}
