package main

import (
	"math"

	"github.com/fsouza/go-dockerclient"
)

type Metrics struct {
	CPUUtil    int
	NetTx      int64
	NetRx      int64
	MemLimit   int64
	MemPercent int
	MemUsage   int64
}

type MetricsReader struct {
	Metrics
	lastCpu    float64
	lastSysCpu float64
}

func NewMetricsReader() *MetricsReader {
	return &MetricsReader{}
}

func (m *MetricsReader) Read(statsCh chan *docker.Stats) chan Metrics {
	stream := make(chan Metrics)

	go func() {
		for s := range statsCh {
			m.ReadCPU(s)
			m.ReadMem(s)
			m.ReadNet(s)
			stream <- m.Metrics
		}
	}()

	return stream
}

func (m *MetricsReader) ReadCPU(stats *docker.Stats) {
	ncpus := float64(len(stats.CPUStats.CPUUsage.PercpuUsage))
	total := float64(stats.CPUStats.CPUUsage.TotalUsage)
	system := float64(stats.CPUStats.SystemCPUUsage)

	cpudiff := total - m.lastCpu
	syscpudiff := system - m.lastSysCpu
	m.CPUUtil = round((cpudiff / syscpudiff * 100) * ncpus)
	m.lastCpu = total
	m.lastSysCpu = system
}

func (m *MetricsReader) ReadMem(stats *docker.Stats) {
	m.MemUsage = int64(stats.MemoryStats.Usage)
	m.MemLimit = int64(stats.MemoryStats.Limit)
	m.MemPercent = round((float64(m.MemUsage) / float64(m.MemLimit)) * 100)
}

func (m *MetricsReader) ReadNet(stats *docker.Stats) {
	var rx, tx int64
	for _, network := range stats.Networks {
		rx += int64(network.RxBytes)
		tx += int64(network.TxBytes)
	}
	m.NetRx, m.NetTx = rx, tx
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}
