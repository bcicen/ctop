package main

import (
	"github.com/fsouza/go-dockerclient"
)

var filters = map[string][]string{
	"status": []string{"running"},
}

func NewContainerMap() *ContainerMap {
	// init docker client
	client, err := docker.NewClient(GlobalConfig["dockerHost"])
	if err != nil {
		panic(err)
	}

	cm := &ContainerMap{
		client:     client,
		containers: make(map[string]*Container),
	}
	cm.Refresh()
	return cm
}

type ContainerMap struct {
	client     *docker.Client
	containers map[string]*Container
}

func (cm *ContainerMap) Refresh() {
	var id string
	opts := docker.ListContainersOptions{
		Filters: filters,
	}
	containers, err := cm.client.ListContainers(opts)
	if err != nil {
		panic(err)
	}
	for _, c := range containers {
		id = c.ID[:12]
		if _, ok := cm.containers[id]; ok == false {
			cm.containers[id] = NewContainer(c)
			cm.containers[id].Collect(cm.client)
		}
	}
}

// Return number of containers/rows
func (cm *ContainerMap) Len() uint {
	return uint(len(cm.containers))
}

// Get a single container, by ID
func (cm *ContainerMap) Get(id string) *Container {
	return cm.containers[id]
}

// Return array of all containers, sorted by field
func (cm *ContainerMap) All() []*Container {
	var containers []*Container
	for _, c := range cm.containers {
		containers = append(containers, c)
	}
	SortContainers(GlobalConfig["sortField"], containers)
	return containers
}
