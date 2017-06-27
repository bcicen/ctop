package collector

import (
	"math"

	"github.com/bcicen/ctop/logging"
	"github.com/bcicen/ctop/models"
)

var log = logging.Init()

type Collector interface {
	Stream() chan models.Metrics
	StreamLogs() (chan string, error)
	Running() bool
	Start()
	Stop()
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

// return rounded percentage
func percent(val float64, total float64) int {
	if total <= 0 {
		return 0
	}
	return round((val / total) * 100)
}
