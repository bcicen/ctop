package connector

import (
	"fmt"
	"strings"
	"sync"

	"github.com/bcicen/ctop/connector/collector"
	api "github.com/fsouza/go-dockerclient"
	"github.com/bcicen/ctop/config"
	"context"
	"github.com/bcicen/ctop/entity"
	"github.com/docker/docker/api/types/swarm"
)

type Docker struct {
	client       *api.Client
	containers   map[string]*entity.Container
	services     map[string]*entity.Service
	nodes        map[string]*entity.Node
	needsRefresh chan string // container IDs requiring refresh
	lock         sync.RWMutex
}

func NewDocker() Connector {
	// init docker client
	client, err := api.NewClientFromEnv()
	if err != nil {
		panic(err)
	}
	cm := &Docker{
		client:       client,
		containers:   make(map[string]*entity.Container),
		services:     make(map[string]*entity.Service),
		nodes:        make(map[string]*entity.Node),
		needsRefresh: make(chan string, 60),
		lock:         sync.RWMutex{},
	}
	go cm.Loop()
	cm.refreshAllContainers()
	if config.GetSwitchVal("swarmMode") {
		cm.refreshAllNodes()
		cm.refreshAllServices()
	}
	go cm.watchEvents()
	return cm
}

// Docker events watcher
func (cm *Docker) watchEvents() {
	log.Info("docker event listener starting")
	events := make(chan *api.APIEvents)
	cm.client.AddEventListener(events)

	for e := range events {
		if e.Type == "container" {
			log.Debugf("Container")
			switch e.Action {
			case "start", "die", "pause", "unpause":
				log.Debugf("handling docker event: action=%s id=%s", e.Action, e.ID)
				cm.needsRefresh <- e.ID
			case "destroy":
				log.Debugf("handling docker event: action=%s id=%s", e.Action, e.ID)
				cm.delByID(e.ID)
			}
		} else if e.Type == "node" {
			log.Debugf("NODE. Action: %s, ID: %s", e.Action, e.ID)
			cm.needsRefresh <- e.ID
		} else if e.Type == "service"{
			log.Debugf("SERVICE. Action: %s, ID: %s", e.Action, e.ID)
			cm.needsRefresh <- e.ID
		}
	}
}

func portsFormat(ports map[api.Port][]api.PortBinding) string {
	var exposed []string
	var published []string

	for k, v := range ports {
		if len(v) == 0 {
			exposed = append(exposed, string(k))
			continue
		}
		for _, binding := range v {
			s := fmt.Sprintf("%s:%s -> %s", binding.HostIP, binding.HostPort, k)
			published = append(published, s)
		}
	}

	return strings.Join(append(exposed, published...), "\n")
}

func (cm *Docker) refresh(c *entity.Container) {
	insp := cm.inspect(c.Id)
	// remove container if no longer exists
	if insp == nil {
		cm.delByID(c.Id)
		return
	}
	c.SetMeta("name", shortName(insp.Name))
	c.SetMeta("image", insp.Config.Image)
	c.SetMeta("ports", portsFormat(insp.NetworkSettings.Ports))
	c.SetMeta("created", insp.Created.Format("Mon Jan 2 15:04:05 2006"))
	c.SetMeta("health", insp.State.Health.Status)
	c.SetState(insp.State.Status)
}

func (cm *Docker) inspect(id string) *api.Container {
	c, err := cm.client.InspectContainer(id)
	if err != nil {
		if _, ok := err.(*api.NoSuchContainer); ok == false {
			log.Errorf(err.Error())
		}
	}
	return c
}

// Mark all container IDs for refresh
func (cm *Docker) refreshAllContainers() {
	opts := api.ListContainersOptions{All: true}
	allContainers, err := cm.client.ListContainers(opts)
	if err != nil {
		panic(err)
	}

	for _, i := range allContainers {
		c := cm.MustGetContainer(i.ID)
		c.SetMeta("name", shortName(i.Names[0]))
		c.SetState(i.State)
		cm.HealthCheck(i.ID)
		cm.needsRefresh <- c.Id
	}
}

func (cm *Docker) refreshAllServices() {
	log.Noticef("Refresh service start!!")
	ctx, cancel := context.WithCancel(context.Background())
	opts := api.ListServicesOptions{Context: ctx}
	allServices, err := cm.client.ListServices(opts)

	if err != nil {
		panic(err)
	}

	for _, i := range allServices {
		s := cm.MustGetService(i.ID)

		s.SetMeta("name", i.Spec.Annotations.Name)
		labels := ""
		for l := range i.Spec.Annotations.Labels {
			labels += l
		}
		s.SetMeta("labels", labels)
		cm.needsRefresh <- s.Id
	}
	cancel()
}

func (cm *Docker) refreshAllNodes() {
	log.Debugf("Start refreshing Nodes in swarm mode")
	ctx, cancel := context.WithCancel(context.Background())
	opt := api.ListNodesOptions{Context: ctx}
	allNodes, err := cm.client.ListNodes(opt)

	if err != nil {
		panic(err)
	}

	for _, i := range allNodes {
		n := cm.MustGetNode(i.ID)
		n.SetMeta("hostname", i.Description.Hostname)
		n.SetState(fmt.Sprintf("%s", i.Status.State))
		leader, reachable := leaderAndReachable(i.ManagerStatus)
		n.SetMeta("leader", leader)
		n.SetMeta("reachable", reachable)
		cm.needsRefresh <- n.Id
	}

	cancel()
}

func (cm *Docker) Loop() {
	for id := range cm.needsRefresh {
		c := cm.MustGetContainer(id)
		cm.refresh(c)
	}
}

// Get a single container, creating one anew if not existing
func (cm *Docker) MustGetContainer(id string) *entity.Container {
	c, ok := cm.GetContainer(id)
	// append container struct for new containers
	if !ok {
		// create collector
		collector := collector.NewDocker(cm.client, id)
		// create container
		c = entity.NewContainer(id, collector)
		cm.lock.Lock()
		cm.containers[id] = c
		cm.lock.Unlock()
	}
	return c
}

func (cm *Docker) MustGetService(id string) *entity.Service {
	s, ok := cm.GetService(id)

	if !ok {
		collector := collector.NewDocker(cm.client, id)
		s = entity.NewService(id, collector)
		cm.lock.Lock()
		cm.services[id] = s
		cm.lock.Unlock()
	}
	return s
}

func (cm *Docker) MustGetNode(id string) *entity.Node {
	n, ok := cm.GetNode(id)
	if !ok {
		collector := collector.NewDocker(cm.client, id)
		n = entity.NewNode(id, collector)
		cm.lock.Lock()
		cm.nodes[id] = n
		cm.lock.Unlock()
	}
	return n
}

// Get a single container, by ID
func (cm *Docker) GetContainer(id string) (*entity.Container, bool) {
	cm.lock.Lock()
	c, ok := cm.containers[id]
	cm.lock.Unlock()
	return c, ok
}

func (cm *Docker) GetService(id string) (*entity.Service, bool) {
	cm.lock.Lock()
	s, ok := cm.services[id]
	cm.lock.Unlock()
	return s, ok
}

func (cm *Docker) GetNode(id string) (*entity.Node, bool) {
	cm.lock.Lock()
	n, ok := cm.nodes[id]
	cm.lock.Unlock()
	return n, ok
}

// Remove containers by ID
func (cm *Docker) delByID(id string) {
	cm.lock.Lock()
	delete(cm.containers, id)
	cm.lock.Unlock()
	log.Infof("removed dead container: %s", id)
}

// Return array of all containers, sorted by field
func (cm *Docker) All() (containers entity.Containers, services entity.Services) {
	cm.lock.Lock()
	for _, c := range cm.containers {
		containers = append(containers, c)
		cm.lock.Unlock()
		cm.HealthCheck(c.Id)
		cm.lock.Lock()
	}

	for _, s := range cm.services {
		services = append(services, s)
	}
	containers.Sort()
	containers.Filter()
	cm.lock.Unlock()
	return containers, services
}

// use primary container name
func shortName(name string) string {
	return strings.Replace(name, "/", "", 1)
}

func (cm *Docker) HealthCheck(id string) {
	insp := cm.inspect(id)
	c := cm.MustGetContainer(id)
	c.SetMeta("health", insp.State.Health.Status)
}

func leaderAndReachable(n *swarm.ManagerStatus) (string, string) {
	reachable := fmt.Sprintf("%s",n.Reachability)
	if n.Leader {
		return "Leader", reachable
	}
	return "", reachable
}
