// +build !release

package connector

import (
	"math/rand"
	"strings"
	"time"

	"github.com/bcicen/ctop/connector/collector"
	"github.com/bcicen/ctop/connector/manager"
	"github.com/bcicen/ctop/entity"
	"github.com/jgautheron/codename-generator"
	"github.com/nu7hatch/gouuid"
)

// Mock for connector
type Mock struct {
	containers entity.Containers
	services   entity.Services
	nodes      entity.Nodes
}

// NewMock return new instance of Mock
func NewMock() *Mock {
	cs := &Mock{}
	go cs.Init()
	go cs.Loop()
	return cs
}

// Init create Mock containers
func (cs *Mock) Init() {
	rand.Seed(int64(time.Now().Nanosecond()))

	for i := 0; i < 4; i++ {
		cs.makeContainer(3)
	}

	for i := 0; i < 16; i++ {
		cs.makeContainer(1)
	}

}

func (cs *Mock) makeContainer(aggression int64) {
	collector := collector.NewMock(aggression)
	manager := manager.NewMock()
	c := entity.NewContainer(makeID(), collector, manager)
	c.SetMeta("name", makeName())
	c.SetState(makeState())
	cs.containers = append(cs.containers, c)
}

// Loop container of mock
func (cs *Mock) Loop() {
	iter := 0
	for {
		// Change state for random container
		if iter%5 == 0 && len(cs.containers) > 0 {
			randC := cs.containers[rand.Intn(len(cs.containers))]
			randC.SetState(makeState())
		}
		iter++
		time.Sleep(3 * time.Second)
	}
}

// GetContainer get a single container, by ID
func (cs *Mock) GetContainer(id string) (*entity.Container, bool) {
	for _, c := range cs.containers {
		if c.Id == id {
			return c, true
		}
	}
	return nil, false
}

// GetTask return a single task by ID
func (cs *Mock) GetTask(id string) (*entity.Task, bool) {
	return nil, false
}

// AllNodes Return slice of all containers, sorted by field
func (cs *Mock) AllNodes() entity.Nodes {
	//cs.nodes.Sort()
	//cs.nodes.Filter()
	return cs.nodes
}

// AllServices return slice of all Service
func (cs *Mock) AllServices() entity.Services {
	//cs.services.Sort()
	//cs.services.Filter()
	return cs.services
}

// AllContainers return slice of all Container
func (cs *Mock) AllContainers() entity.Containers {
	cs.containers.Sort()
	cs.containers.Filter()
	return cs.containers
}

// Remove containers by ID
func (cs *Mock) delByID(id string) {
	for n, c := range cs.containers {
		if c.Id == id {
			cs.del(n)
			return
		}
	}
}

// Remove one or more containers by index
func (cs *Mock) del(idx ...int) {
	for _, i := range idx {
		cs.containers = append(cs.containers[:i], cs.containers[i+1:]...)
	}
	log.Infof("removed %d dead containers", len(idx))
}

func makeID() string {
	u, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}
	return strings.Replace(u.String(), "-", "", -1)[:12]
}

func makeName() string {
	n, err := codename.Get(codename.Sanitized)
	nsp := strings.Split(n, "-")
	if len(nsp) > 2 {
		n = strings.Join(nsp[:2], "-")
	}
	if err != nil {
		panic(err)
	}
	return strings.Replace(n, "-", "_", -1)
}

func makeState() string {
	switch rand.Intn(10) {
	case 0, 1, 2:
		return "exited"
	case 3:
		return "paused"
	}
	return "running"
}

// Down mock
func (cs *Mock) Down() {
	log.Warning("Not implent Down for Mock")
}
