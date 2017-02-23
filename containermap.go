package main

import (
	"sort"
	"strings"

	"github.com/bcicen/ctop/config"
	"github.com/bcicen/ctop/metrics"
	"github.com/bcicen/ctop/widgets"
	"github.com/fsouza/go-dockerclient"
)

type ContainerMap struct {
	client     *docker.Client
	containers Containers
	collectors map[string]metrics.Collector
}

func NewContainerMap() *ContainerMap {
	// init docker client
	client, err := docker.NewClient(config.GetVal("dockerHost"))
	if err != nil {
		panic(err)
	}
	cm := &ContainerMap{
		client:     client,
		collectors: make(map[string]metrics.Collector),
	}
	//cm.Refresh()
	return cm
}

func (cm *ContainerMap) Refresh() {
	var id, name string

	opts := docker.ListContainersOptions{All: true}
	allContainers, err := cm.client.ListContainers(opts)
	if err != nil {
		panic(err)
	}

	// add new containers
	states := make(map[string]string)
	for _, c := range allContainers {
		id = c.ID[:12]
		states[id] = c.State

		if _, ok := cm.Get(id); ok == false {
			name = strings.Replace(c.Names[0], "/", "", 1) // use primary container name
			newc := &Container{
				id:      id,
				name:    name,
				widgets: widgets.NewCompact(id, name),
			}
			cm.containers = append(cm.containers, newc)
		}

		if _, ok := cm.collectors[id]; ok == false {
			cm.collectors[id] = metrics.NewDocker(cm.client, id)
		}

	}

	var removeIdxs []int
	for n, c := range cm.containers {

		// mark stale internal containers
		if _, ok := states[c.id]; ok == false {
			removeIdxs = append(removeIdxs, n)
			continue
		}

		c.SetState(states[c.id])
		// start collector if needed
		//collector := cm.collectors[id]
		if c.state == "running" && !cm.collectors[c.id].Running() {
			cm.collectors[c.id].Start()
			c.Read(cm.collectors[c.id].Stream())
		}
	}

	// delete removed containers
	cm.Del(removeIdxs...)
}

// Kill a container by ID
func (cm *ContainerMap) Kill(id string, sig docker.Signal) error {
	opts := docker.KillContainerOptions{
		ID:     id,
		Signal: sig,
	}
	return cm.client.KillContainer(opts)
}

// Return number of containers/rows
func (cm *ContainerMap) Len() uint {
	return uint(len(cm.containers))
}

// Get a single container, by ID
func (cm *ContainerMap) Get(id string) (*Container, bool) {
	for _, c := range cm.containers {
		if c.id == id {
			return c, true
		}
	}
	return nil, false
}

// Remove one or more containers by index
func (cm *ContainerMap) Del(idx ...int) {
	for _, i := range idx {
		cm.containers = append(cm.containers[:i], cm.containers[i+1:]...)
	}
}

// Return array of all containers, sorted by field
func (cm *ContainerMap) All() []*Container {
	sort.Sort(cm.containers)
	return cm.containers
}
