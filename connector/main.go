package connector

import (
	"fmt"
	"sort"

	"github.com/bcicen/ctop/container"
	"github.com/bcicen/ctop/logging"
)

var (
	log     = logging.Init()
	enabled = make(map[string]func() Connector)
)

// return names for all enabled connectors on the current platform
func Enabled() (a []string) {
	for k, _ := range enabled {
		a = append(a, k)
	}
	sort.Strings(a)
	return a
}

func ByName(s string) (Connector, error) {
	if cfn, ok := enabled[s]; ok {
		return cfn(), nil
	}
	return nil, fmt.Errorf("invalid connector type \"%s\"", s)
}

type Connector interface {
	All() container.Containers
	Get(string) (*container.Container, bool)
}
