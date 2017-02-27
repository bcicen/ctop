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
	collectors map[string]metrics.Collector
}

func NewMockContainerSource() *MockContainerSource {
	cs := &MockContainerSource{
		collectors: make(map[string]metrics.Collector),
	}
	cs.Init()
	go cs.Loop()
	return cs
}

// Create Mock containers
func (cs *MockContainerSource) Init() {
	total := 10
	rand.Seed(int64(time.Now().Nanosecond()))

	for i := 0; i < total; i++ {
		c := NewContainer(makeID(), makeName())
		lock.Lock()
		cs.containers = append(cs.containers, c)
		lock.Unlock()
		cs.collectors[c.id] = metrics.NewMock()

		c.SetState(makeState())
	}

}

func (cs *MockContainerSource) Loop() {
	iter := 0
	for {
		for _, c := range cs.containers {
			// Change state for random container
			if iter%5 == 0 {
				randC := cs.containers[rand.Intn(len(cs.containers))]
				randC.SetState(makeState())
			}

			isCollecting := cs.collectors[c.id].Running()
			//log.Infof("id=%s state=%s collector=%t", c.id, c.state, isCollecting)

			// start collector if needed
			if c.state == "running" && !isCollecting {
				cs.collectors[c.id].Start()
				c.Read(cs.collectors[c.id].Stream())
			}
			// stop collector if needed
			if c.state != "running" && isCollecting {
				cs.collectors[c.id].Stop()
			}

		}
		iter++
		time.Sleep(3 * time.Second)
	}
}

// Get a single container, by ID
func (cs *MockContainerSource) Get(id string) (*Container, bool) {
	for _, c := range cs.containers {
		if c.id == id {
			return c, true
		}
	}
	return nil, false
}

// Remove containers by ID
func (cs *MockContainerSource) delByID(id string) {
	for n, c := range cs.containers {
		if c.id == id {
			cs.del(n)
			return
		}
	}
}

// Remove one or more containers by index
func (cs *MockContainerSource) del(idx ...int) {
	lock.Lock()
	defer lock.Unlock()
	for _, i := range idx {
		cs.containers = append(cs.containers[:i], cs.containers[i+1:]...)
	}
	log.Infof("removed %d dead containers", len(idx))
}

// Return array of all containers, sorted by field
func (cs *MockContainerSource) All() []*Container {
	sort.Sort(cs.containers)
	return cs.containers
}

func makeID() string {
	u, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}
	return strings.Replace(u.String(), "-", "", -1)
}

func makeName() string {
	n, err := codename.Get(codename.Sanitized)
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
