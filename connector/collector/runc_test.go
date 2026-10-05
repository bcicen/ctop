//go:build linux
// +build linux

package collector

import (
	"runtime"
	"testing"

	"github.com/opencontainers/runc/libcontainer/cgroups"
	"github.com/opencontainers/runc/types"
)

func TestRuncReadCPUNilStats(t *testing.T) {
	c := &Runc{interval: 1}
	c.ReadCPU(nil)

	if c.NCpus != uint8(runtime.NumCPU()) {
		t.Errorf("expected NCpus fallback to %d, got %d", runtime.NumCPU(), c.NCpus)
	}
	if c.CPUUtil != 0 {
		t.Errorf("expected CPUUtil 0, got %d", c.CPUUtil)
	}
	if c.Pids != 0 {
		t.Errorf("expected Pids 0, got %d", c.Pids)
	}
}

func TestRuncReadCPUCgroupV2EmptyPercpu(t *testing.T) {
	c := &Runc{interval: 1}
	stats := &cgroups.Stats{}
	stats.CpuStats.CpuUsage.TotalUsage = 500000000
	stats.CpuStats.CpuUsage.PercpuUsage = nil
	stats.PidsStats.Current = 12

	c.ReadCPU(stats)
	if c.NCpus != uint8(runtime.NumCPU()) {
		t.Errorf("expected NCpus fallback to %d, got %d", runtime.NumCPU(), c.NCpus)
	}
	if c.Pids != 12 {
		t.Errorf("expected Pids 12, got %d", c.Pids)
	}
}

func TestRuncReadMemNilStats(t *testing.T) {
	c := &Runc{interval: 1}
	c.ReadMem(nil)

	if c.MemUsage != 0 || c.MemLimit != 0 || c.MemPercent != 0 {
		t.Errorf("expected zero mem metrics for nil stats, got usage=%d limit=%d percent=%d",
			c.MemUsage, c.MemLimit, c.MemPercent)
	}
}

func TestRuncReadMemPopulated(t *testing.T) {
	c := &Runc{interval: 1}
	stats := &cgroups.Stats{}
	stats.MemoryStats.Usage.Usage = 50000000
	stats.MemoryStats.Usage.Limit = 100000000

	c.ReadMem(stats)
	if c.MemUsage != 50000000 {
		t.Errorf("expected MemUsage 50000000, got %d", c.MemUsage)
	}
	if c.MemPercent != 50 {
		t.Errorf("expected MemPercent 50, got %d", c.MemPercent)
	}
}

func TestRuncReadNetNilInterfaces(t *testing.T) {
	c := &Runc{interval: 1}
	c.ReadNet(nil)

	if c.NetRx != 0 || c.NetTx != 0 {
		t.Errorf("expected zero net metrics for nil interfaces, got rx=%d tx=%d", c.NetRx, c.NetTx)
	}
}

func TestRuncReadNetNilElement(t *testing.T) {
	c := &Runc{interval: 1}
	interfaces := []*types.NetworkInterface{
		nil,
		{RxBytes: 100, TxBytes: 200},
	}
	c.ReadNet(interfaces)
	if c.NetRx != 100 || c.NetTx != 200 {
		t.Errorf("expected rx=100 tx=200, got rx=%d tx=%d", c.NetRx, c.NetTx)
	}
}

func TestRuncReadIONilStats(t *testing.T) {
	c := &Runc{interval: 1}
	c.ReadIO(nil)

	if c.IOBytesRead != 0 || c.IOBytesWrite != 0 {
		t.Errorf("expected zero IO metrics for nil stats, got read=%d write=%d",
			c.IOBytesRead, c.IOBytesWrite)
	}
}

func TestRuncReadIOPopulated(t *testing.T) {
	c := &Runc{interval: 1}
	stats := &cgroups.Stats{}
	stats.BlkioStats.IoServiceBytesRecursive = []cgroups.BlkioStatEntry{
		{Op: "Read", Value: 2048},
		{Op: "Write", Value: 4096},
	}

	c.ReadIO(stats)
	if c.IOBytesRead != 2048 || c.IOBytesWrite != 4096 {
		t.Errorf("expected read=2048 write=4096, got read=%d write=%d",
			c.IOBytesRead, c.IOBytesWrite)
	}
}
