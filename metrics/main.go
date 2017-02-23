package metrics

import (
	"math"
)

type Metrics struct {
	CPUUtil    int
	NetTx      int64
	NetRx      int64
	MemLimit   int64
	MemPercent int
	MemUsage   int64
}

type Collector interface {
	Stream() chan Metrics
	Running() bool
	Start()
	Stop()
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}
