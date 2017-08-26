package entity

import (
  "github.com/bcicen/ctop/logging"
  "github.com/bcicen/ctop/connector/collector"
  "github.com/bcicen/ctop/cwidgets"
)

var (
  log = logging.Init()
)

type Entity interface {
  SetUpdater(u cwidgets.WidgetUpdater)
  SetState(s string)
  SetMeta(k, v string)
  GetMeta(k string) string
  Logs() collector.LogCollector
}
