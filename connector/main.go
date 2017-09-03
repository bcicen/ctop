package connector

import (
	"fmt"

	"github.com/bcicen/ctop/entity"
	"github.com/bcicen/ctop/logging"
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
	AllNodes() entity.Nodes
	AllServices() entity.Services
	AllContainers() entity.Containers
	AllTasks() entity.Tasks
	GetContainer(string) (*entity.Container, bool)
	GetService(string) (*entity.Service, bool)
	DownSwarmMode()
}
