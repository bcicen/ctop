package connector

import (
	"github.com/bcicen/ctop/container"
	"github.com/bcicen/ctop/logging"
)

var log = logging.Init()

type ContainerSource interface {
	All() container.Containers
	Get(string) (*container.Container, bool)
}
