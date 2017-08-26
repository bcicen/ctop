package entity

import (
	"github.com/bcicen/ctop/logging"
	"github.com/bcicen/ctop/connector/collector"
)

var (
	log = logging.Init()
)

type Entity interface {
	SetState(s string)
	Logs() collector.LogCollector
}
