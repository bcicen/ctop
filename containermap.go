package main

import (
	"strings"

	"github.com/fsouza/go-dockerclient"
)

var filters = map[string][]string{
	"status": []string{"running"},
}

func NewContainerMap() *ContainerMap {
	config := DefaultConfig

	// init docker client
	client, err := docker.NewClient(config.dockerHost)
	if err != nil {
		panic(err)
	}

	cm := &ContainerMap{
		config:     config,
		client:     client,
		containers: make(map[string]*Container),
	}
	cm.Refresh()
	return cm
}

type ContainerMap struct {
	config     Config
	client     *docker.Client
	containers map[string]*Container
}

func (cm *ContainerMap) Refresh() {
	opts := docker.ListContainersOptions{
		Filters: filters,
	}
	containers, err := cm.client.ListContainers(opts)
	if err != nil {
		panic(err)
	}
	for _, c := range containers {
		if _, ok := cm.containers[c.ID[:12]]; ok == false {
			cm.Add(c)
		}
	}
}

// Return number of containers/rows
func (cm *ContainerMap) Len() uint {
	return uint(len(cm.containers))
}

func (cm *ContainerMap) Add(c docker.APIContainers) {
	id := c.ID[:12]
	name := strings.Replace(c.Names[0], "/", "", 1) // use primary container name
	cm.containers[id] = &Container{
		id:      id,
		name:    name,
		done:    make(chan bool),
		stats:   make(chan *docker.Stats),
		widgets: NewWidgets(id, name),
		reader:  &StatReader{},
	}
	cm.containers[id].Collect(cm.client)
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
	SortContainers(cm.config.sortField, containers)
	return containers
}
