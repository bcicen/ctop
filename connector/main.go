package connector

import (
	"fmt"

	"github.com/bcicen/ctop/container"
	"github.com/bcicen/ctop/logging"
	"github.com/bcicen/ctop/service"
)

var log = logging.Init()

func ByName(s string) (Connector, error) {
	if _, ok := enabled[s]; !ok {
		msg := fmt.Sprintf("invalid connector type \"%s\"\nconnector must be one of:", s)
		for k, _ := range enabled {
			msg += fmt.Sprintf("\n  %s", k)
		}
		return nil, fmt.Errorf(msg)
	}
	return enabled[s](), nil
}

type Connector interface {
	All() (container.Containers, service.Services)
	GetContainer(string) (*container.Container, bool)
	GetService(string) (*service.Service, bool)
}
