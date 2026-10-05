package collector

import (
	"runtime"
	"testing"

	api "github.com/fsouza/go-dockerclient"
)

func TestDockerReadCPUNilStats(t *testing.T) {
	c := NewDocker(nil, "test-container")
	// Should not panic on nil stats
	c.ReadCPU(nil)

	if c.NCpus != uint8(runtime.NumCPU()) {
		t.Errorf("expected NCpus to fallback to %d, got %d", runtime.NumCPU(), c.NCpus)
	}
	if c.CPUUtil != 0 {
		t.Errorf("expected CPUUtil 0, got %d", c.CPUUtil)
	}
	if c.Pids != 0 {
		t.Errorf("expected Pids 0, got %d", c.Pids)
	}
}

func TestDockerReadCPUCgroupV2MissingPercpu(t *testing.T) {
	c := NewDocker(nil, "test-container")
	stats := &api.Stats{}
	stats.CPUStats.CPUUsage.TotalUsage = 1000
	stats.CPUStats.SystemCPUUsage = 10000
	// OnlineCPUs is 0 and PercpuUsage is empty (cgroup v2 behavior)
	stats.CPUStats.OnlineCPUs = 0
	stats.CPUStats.CPUUsage.PercpuUsage = nil
	stats.PidsStats.Current = 5

	c.ReadCPU(stats)

	if c.NCpus != uint8(runtime.NumCPU()) {
		t.Errorf("expected NCpus fallback to NumCPU %d, got %d", runtime.NumCPU(), c.NCpus)
	}
	if c.Pids != 5 {
		t.Errorf("expected Pids 5, got %d", c.Pids)
	}
}

func TestDockerReadCPUHappyPath(t *testing.T) {
	c := NewDocker(nil, "test-container")
	stats := &api.Stats{}
	stats.CPUStats.OnlineCPUs = 4
	stats.CPUStats.CPUUsage.PercpuUsage = []uint64{250, 250, 250, 250}
	stats.CPUStats.CPUUsage.TotalUsage = 1000
	stats.CPUStats.SystemCPUUsage = 10000

	c.ReadCPU(stats)
	if c.NCpus != 4 {
		t.Errorf("expected NCpus 4, got %d", c.NCpus)
	}

	// Second sample with progress
	stats.CPUStats.CPUUsage.TotalUsage = 2000
	stats.CPUStats.SystemCPUUsage = 12000
	c.ReadCPU(stats)
	if c.CPUUtil <= 0 {
		t.Errorf("expected positive CPUUtil, got %d", c.CPUUtil)
	}
}

func TestDockerReadMemNilStats(t *testing.T) {
	c := NewDocker(nil, "test-container")
	c.ReadMem(nil)

	if c.MemUsage != 0 || c.MemLimit != 0 || c.MemPercent != 0 {
		t.Errorf("expected zero memory metrics for nil stats, got usage=%d limit=%d percent=%d",
			c.MemUsage, c.MemLimit, c.MemPercent)
	}
}

func TestDockerReadMemCgroupV2InactiveFile(t *testing.T) {
	c := NewDocker(nil, "test-container")
	stats := &api.Stats{}
	stats.MemoryStats.Usage = 100000000
	stats.MemoryStats.Limit = 200000000
	// Under cgroup v2, Cache is 0 and InactiveFile holds cache
	stats.MemoryStats.Stats.Cache = 0
	stats.MemoryStats.Stats.InactiveFile = 20000000

	c.ReadMem(stats)
	expectedUsage := int64(80000000)
	if c.MemUsage != expectedUsage {
		t.Errorf("expected MemUsage %d, got %d", expectedUsage, c.MemUsage)
	}
	if c.MemLimit != 200000000 {
		t.Errorf("expected MemLimit 200000000, got %d", c.MemLimit)
	}
	if c.MemPercent != 40 {
		t.Errorf("expected MemPercent 40, got %d", c.MemPercent)
	}
}

func TestDockerReadMemUnderflowProtection(t *testing.T) {
	c := NewDocker(nil, "test-container")
	stats := &api.Stats{}
	// Cache reported greater than Usage due to kernel accounting race
	stats.MemoryStats.Usage = 5000
	stats.MemoryStats.Limit = 100000
	stats.MemoryStats.Stats.Cache = 8000

	c.ReadMem(stats)
	if c.MemUsage != 5000 {
		t.Errorf("expected MemUsage fallback to 5000 without underflow, got %d", c.MemUsage)
	}
}

func TestDockerReadNetNilStats(t *testing.T) {
	c := NewDocker(nil, "test-container")
	c.ReadNet(nil)

	if c.NetRx != 0 || c.NetTx != 0 {
		t.Errorf("expected zero net metrics for nil stats, got rx=%d tx=%d", c.NetRx, c.NetTx)
	}
}

func TestDockerReadNetHappyPath(t *testing.T) {
	c := NewDocker(nil, "test-container")
	stats := &api.Stats{}
	stats.Networks = map[string]api.NetworkStats{
		"eth0": {RxBytes: 1234, TxBytes: 5678},
		"eth1": {RxBytes: 1000, TxBytes: 2000},
	}

	c.ReadNet(stats)
	if c.NetRx != 2234 || c.NetTx != 7678 {
		t.Errorf("expected rx=2234 tx=7678, got rx=%d tx=%d", c.NetRx, c.NetTx)
	}
}

func TestDockerReadIONilStats(t *testing.T) {
	c := NewDocker(nil, "test-container")
	c.ReadIO(nil)

	if c.IOBytesRead != 0 || c.IOBytesWrite != 0 {
		t.Errorf("expected zero IO metrics for nil stats, got read=%d write=%d",
			c.IOBytesRead, c.IOBytesWrite)
	}
}

func TestDockerReadIOHappyPath(t *testing.T) {
	c := NewDocker(nil, "test-container")
	stats := &api.Stats{}
	stats.BlkioStats.IOServiceBytesRecursive = []api.BlkioStatsEntry{
		{Op: "Read", Value: 500},
		{Op: "Write", Value: 1500},
		{Op: "Read", Value: 300},
	}

	c.ReadIO(stats)
	if c.IOBytesRead != 800 || c.IOBytesWrite != 1500 {
		t.Errorf("expected read=800 write=1500, got read=%d write=%d",
			c.IOBytesRead, c.IOBytesWrite)
	}
}
