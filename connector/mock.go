// +build !release

package connector

import (
	"math/rand"
	"strings"
	"time"

	"github.com/bcicen/ctop/connector/collector"
	"github.com/bcicen/ctop/connector/manager"
	"github.com/bcicen/ctop/container"
	"github.com/jgautheron/codename-generator"
	"github.com/nu7hatch/gouuid"
)

func init() { enabled["mock"] = NewMock }

type Mock struct {
	containers container.Containers
	noneStack  *container.Stack
}

func NewMock() (Connector, error) {
	cs := &Mock{
		noneStack: container.NewStack("", "none"),
	}
	go cs.Init()
	go cs.Loop()
	return cs, nil
}

// Create Mock containers
func (cs *Mock) Init() {
	rand.Seed(int64(time.Now().Nanosecond()))
	stack1 := container.NewStack("stack1", "compose")
	stack2 := container.NewStack("stack2", "compose")

	for i := 0; i < 2; i++ {
		cs.makeContainer(3, true, cs.noneStack)
		cs.makeContainer(3, true, stack1)
	}

	for i := 0; i < 8; i++ {
		cs.makeContainer(1, false, cs.noneStack)
		cs.makeContainer(1, false, stack2)
	}

}

func (cs *Mock) Wait() struct{} {
	ch := make(chan struct{})
	go func() {
		time.Sleep(30 * time.Second)
		close(ch)
	}()
	return <-ch
}

var healthStates = []string{"starting", "healthy", "unhealthy"}

func (cs *Mock) makeContainer(aggression int64, health bool, stack *container.Stack) {
	collector := collector.NewMock(aggression)
	manager := manager.NewMock()
	c := container.New(makeID(), collector, manager)
	c.Stack = stack
	c.Stack.Count++
	c.SetMeta("name", makeName())
	c.SetState(makeState())
	if health {
		var i int
		c.SetMeta("health", healthStates[i])
		go func() {
			for {
				i++
				if i >= len(healthStates) {
					i = 0
				}
				c.SetMeta("health", healthStates[i])
				time.Sleep(12 * time.Second)
			}
		}()
	}
	cs.containers = append(cs.containers, c)
}

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

// Get a single container, by ID
func (cs *Mock) Get(id string) (*container.Container, bool) {
	for _, c := range cs.containers {
		if c.Id == id {
			return c, true
		}
	}
	return nil, false
}

// All returns array of all containers, sorted by field
func (cs *Mock) All() container.Containers {
	cs.containers.Sort()
	cs.containers.Filter()
	return cs.containers
}

// Remove containers by ID
func (cs *Mock) delByID(id string) {
	for n, c := range cs.containers {
		if c.Id == id {
			c.Stack.Count--
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
