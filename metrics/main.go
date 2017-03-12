package metrics

import (
	"math"

	"github.com/bcicen/ctop/logging"
)

var log = logging.Init()

type Metrics struct {
	CPUUtil      int
	NetTx        int64
	NetRx        int64
	MemLimit     int64
	MemPercent   int
	MemUsage     int64
	IOBytesRead  int64
	IOBytesWrite int64
	Pids         int
}

func NewMetrics() Metrics {
	return Metrics{
		CPUUtil:           -1,
		NetTx:             -1,
		NetRx:             -1,
		MemUsage:          -1,
		MemPercent:        -1,
		IOBytesRead:       -1,
		IOBytesWrite:      -1,
		Pids:              -1,
	}
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
