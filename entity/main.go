package entity

import (
	"github.com/bcicen/ctop/logging"
	"github.com/bcicen/ctop/connector/collector"
	"github.com/bcicen/ctop/models"
	"github.com/bcicen/ctop/cwidgets"
)

var (
	log = logging.Init()
)

type Entity interface {
	SetState(s string)
	Logs() collector.LogCollector
	GetMetaEntity() Meta
	SetUpdater(updater cwidgets.WidgetUpdater)
	GetMeta(v string) string
	GetId() string
	GetMetrics() models.Metrics
}

