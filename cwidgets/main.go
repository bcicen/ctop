package cwidgets

import (
	"github.com/bcicen/ctop/logging"
	"github.com/bcicen/ctop/models"
)

var log = logging.Init()

type WidgetUpdater interface {
	SetMeta(string, string)
	SetMetrics(models.Metrics)
}
