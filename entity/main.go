package entity

import "github.com/bcicen/ctop/logging"

var (
  log = logging.Init()
)

type Entity interface {
  SetState() (s string)
}
