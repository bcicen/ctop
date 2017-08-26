package entity

import (
	"github.com/bcicen/ctop/logging"
	"github.com/bcicen/ctop/connector/collector"
	"github.com/bcicen/ctop/models"
)

var (
	log = logging.Init()
)

type Entity interface {
	SetState(s string)
	Logs() collector.LogCollector
	GetMetaEntity() Meta
	GetId() string
	GetMetrics() models.Metrics
}

