package main

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/bcicen/ctop/collector"
	"github.com/bcicen/ctop/widgets"
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
	var id, name string
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
			name = strings.Replace(c.Names[0], "/", "", 1) // use primary container name
			cm.containers[id] = &Container{
				id:      id,
				name:    name,
				done:    make(chan bool),
				collect: collector.NewDocker(cm.client, id),
				widgets: widgets.NewCompact(id, name),
			}
			cm.containers[id].Collect()
		}
	}
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
func (cm *ContainerMap) Get(id string) *Container {
	return cm.containers[id]
}

// Return array of all containers, sorted by field
func (cm *ContainerMap) All() []*Container {
	var containers Containers
	filter := GlobalConfig["filterStr"]
	re := regexp.MustCompile(fmt.Sprintf(".*%s", filter))

	for _, c := range cm.containers {
		if re.FindAllString(c.name, 1) != nil {
			containers = append(containers, c)
		}
	}

	sort.Sort(containers)
	return containers
}
