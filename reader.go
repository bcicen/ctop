package main

import (
	"math"

	"github.com/fsouza/go-dockerclient"
)

type StatReader struct {
	CPUUtil    int
	NetTx      int64
	NetRx      int64
	MemLimit   int64
	MemPercent int
	MemUsage   int64
	lastCpu    float64
	lastSysCpu float64
}

func (s *StatReader) Read(stats *docker.Stats) {
	s.ReadCPU(stats)
	s.ReadMem(stats)
	s.ReadNet(stats)
}

func (s *StatReader) ReadCPU(stats *docker.Stats) {
	ncpus := float64(len(stats.CPUStats.CPUUsage.PercpuUsage))
	total := float64(stats.CPUStats.CPUUsage.TotalUsage)
	system := float64(stats.CPUStats.SystemCPUUsage)

	cpudiff := total - s.lastCpu
	syscpudiff := system - s.lastSysCpu
	s.CPUUtil = round((cpudiff / syscpudiff * 100) * ncpus)
	s.lastCpu = total
	s.lastSysCpu = system
}

func (s *StatReader) ReadMem(stats *docker.Stats) {
	s.MemUsage = int64(stats.MemoryStats.Usage)
	s.MemLimit = int64(stats.MemoryStats.Limit)
	s.MemPercent = round((float64(s.MemUsage) / float64(s.MemLimit)) * 100)
}

func (s *StatReader) ReadNet(stats *docker.Stats) {
	s.NetTx, s.NetRx = 0, 0
	for _, network := range stats.Networks {
		s.NetTx += int64(network.TxBytes)
		s.NetRx += int64(network.RxBytes)
	}
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}
