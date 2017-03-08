// +build !release

package main

import (
	"math/rand"
	"sort"
	"strings"
	"time"

	"github.com/bcicen/ctop/metrics"
	"github.com/jgautheron/codename-generator"
	"github.com/nu7hatch/gouuid"
)

type MockContainerSource struct {
	containers Containers
}

func NewMockContainerSource() *MockContainerSource {
	cs := &MockContainerSource{}
	go cs.Init()
	go cs.Loop()
	return cs
}

// Create Mock containers
func (cs *MockContainerSource) Init() {
	total := 20
	rand.Seed(int64(time.Now().Nanosecond()))

	for i := 0; i < total; i++ {
		//time.Sleep(1 * time.Second)
		collector := metrics.NewMock()
		c := NewContainer(makeID(), collector)
		c.SetMeta("name", makeName())
		c.SetState(makeState())
		cs.containers = append(cs.containers, c)
	}

}

func (cs *MockContainerSource) Loop() {
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
func (cs *MockContainerSource) Get(id string) (*Container, bool) {
	for _, c := range cs.containers {
		if c.Id == id {
			return c, true
		}
	}
	return nil, false
}

// Return array of all containers, sorted by field
func (cs *MockContainerSource) All() Containers {
	sort.Sort(cs.containers)
	cs.containers.Filter()
	return cs.containers
}

// Remove containers by ID
func (cs *MockContainerSource) delByID(id string) {
	for n, c := range cs.containers {
		if c.Id == id {
			cs.del(n)
			return
		}
	}
}

// Remove one or more containers by index
func (cs *MockContainerSource) del(idx ...int) {
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
