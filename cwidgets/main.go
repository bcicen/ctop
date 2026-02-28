package cwidgets

import (
	"github.com/lordoverlord/ctop/logging"
	"github.com/lordoverlord/ctop/models"
)

var log = logging.Init()

type WidgetUpdater interface {
	SetMeta(models.Meta)
	SetMetrics(models.Metrics)
}

type NullWidgetUpdater struct{}

// NullWidgetUpdater implements WidgetUpdater
func (wu NullWidgetUpdater) SetMeta(models.Meta) {}

// NullWidgetUpdater implements WidgetUpdater
func (wu NullWidgetUpdater) SetMetrics(models.Metrics) {}
