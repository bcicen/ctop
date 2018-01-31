// +build linux

package collector

import (
	linuxproc "github.com/c9s/goprocinfo/linux"
	"github.com/opencontainers/runc/libcontainer/system"
)

var sysMemTotal = getSysMemTotal()
var clockTicksPerSecond = uint64(system.GetClockTicks())

const nanoSecondsPerSecond = 1e9

func getSysMemTotal() int64 {
	stat, err := linuxproc.ReadMemInfo("/proc/meminfo")
	if err != nil {
		log.Errorf("error reading system stats: %s", err)
		return 0
	}
	return int64(stat.MemTotal * 1024)
}

// return cumulative system cpu usage in nanoseconds
func getSysCPUUsage() uint64 {
	stat, err := linuxproc.ReadStat("/proc/stat")
	if err != nil {
		log.Errorf("error reading system stats: %s", err)
		return 0
	}

	sum := stat.CPUStatAll.User +
		stat.CPUStatAll.Nice +
		stat.CPUStatAll.System +
		stat.CPUStatAll.Idle +
		stat.CPUStatAll.IOWait +
		stat.CPUStatAll.IRQ +
		stat.CPUStatAll.SoftIRQ +
		stat.CPUStatAll.Steal +
		stat.CPUStatAll.Guest +
		stat.CPUStatAll.GuestNice

	return (sum * nanoSecondsPerSecond) / clockTicksPerSecond
}
