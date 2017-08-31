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
	client                 *api.Client
	containers             map[string]*entity.Container
	services               map[string]*entity.Service
	nodes                  map[string]*entity.Node
	tasks                  map[string]*entity.Task
	needsRefreshNodes      chan string // node IDs requiring refresh
	needsRefreshContainers chan string // service IDs requiring refresh
	needsRefreshTasks      chan string // task IDs requiring refresh
	needsRefreshServices   chan string // container IDs requiring refresh
	lock                   sync.RWMutex
}

func NewDocker() Connector {
	// init docker client
	client, err := api.NewClientFromEnv()
	if err != nil {
		panic(fmt.Sprintf("NewDocker err:%s", err))
	}
	cm := &Docker{
		client:                 client,
		containers:             make(map[string]*entity.Container),
		services:               make(map[string]*entity.Service),
		nodes:                  make(map[string]*entity.Node),
		tasks:                  make(map[string]*entity.Task),
		needsRefreshNodes:      make(chan string, 60),
		needsRefreshContainers: make(chan string, 60),
		needsRefreshServices:   make(chan string, 60),
		needsRefreshTasks:      make(chan string, 60),
		lock:                   sync.RWMutex{},
	}
	if config.GetSwitchVal("swarmMode") {
		go cm.LoopNode()
		go cm.LoopService()
		go cm.LoopTask()
		cm.refreshAllNodes()
		cm.refreshAllServices()
		cm.refreshAllTasks()
	} else {
		go cm.LoopContainer()
		cm.refreshAllContainers()
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
		if config.GetSwitchVal("swarmMode") {
			log.Debugf("Action ", e)
			//if e.Type == "node" {
			//	log.Debugf("NODE. Action: %s, ID: %s", e.Action, e.ID)
			//	cm.needsRefreshNodes <- e.ID
			if e.Type == "service" {
				log.Debugf("Service")
				actionName := strings.Split(e.Action, ":")[0]
				log.Debugf("actionName %s", actionName)

				switch actionName {
				case "update":
					cm.needsRefreshServices <- e.ID
					log.Debugf("SERVICE. Action: %s, ID: %s", e.Action, e.ID)
					cm.refreshAllTasks()
					for _, t := range cm.tasks {
						if t.GetMeta("service") == e.ID {
							cm.needsRefreshTasks <- t.Id
						}
					}
				}
			}
		} else {
			if e.Type == "container" {
				actionName := strings.Split(e.Action, ":")[0]

				switch actionName {
				case "start", "die", "pause", "unpause", "health_status":
					log.Debugf("handling docker event: action=%s id=%s", e.Action, e.ID)
					cm.needsRefreshContainers <- e.ID
				case "destroy":
					log.Debugf("handling docker event: action=%s id=%s", e.Action, e.ID)
					cm.delByIDContainer(e.ID)
				}
			}
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

func (cm *Docker) refreshContainer(c *entity.Container) {
	log.Infof("refreshContainer")
	insp := cm.inspectContainer(c.Id)
	// remove container if no longer exists
	if insp == nil {
		cm.delByIDContainer(c.Id)
		return
	}
	c.SetMeta("name", shortName(insp.Name))
	c.SetMeta("image", insp.Config.Image)
	c.SetMeta("ports", portsFormat(insp.NetworkSettings.Ports))
	c.SetMeta("created", insp.Created.Format("Mon Jan 2 15:04:05 2006"))
	c.SetMeta("health", insp.State.Health.Status)
	c.SetState(insp.State.Status)
}

func (cm *Docker) refreshNode(n *entity.Node) {
	log.Debugf("1")
	insp := cm.inspectNode(n.Id)
	// remove container if no longer exists
	if insp == nil {
		cm.delByIDNode(n.Id)
		return
	}
	n.SetMeta("name", insp.Description.Hostname)
}

func (cm *Docker) refreshService(s *entity.Service) {
	log.Debugf("2")
	insp := cm.inspectService(s.Id)
	// remove container if no longer exists
	if insp == nil {
		cm.delByIDService(s.Id)
		return
	}
	s.SetMeta("name", insp.Spec.Annotations.Name)
}

func (cm *Docker) refreshTask(t *entity.Task) {
	insp := cm.inspectTask(t.Id)
	// remove task if no longer exists
	if insp == nil {
		cm.delByIDTask(t.Id)
		return
	}
	t.SetMeta("name", insp.Annotations.Name)
	t.SetState(fmt.Sprintf("%s", insp.Status.State))
}

func (cm *Docker) inspectContainer(id string) *api.Container {
	c, err := cm.client.InspectContainer(id)
	if err != nil {
		if _, ok := err.(*api.NoSuchContainer); ok == false {
			log.Errorf(err.Error())
		}
	}
	return c
}
func (cm *Docker) inspectNode(id string) *swarm.Node {
	n, err := cm.client.InspectNode(id)
	if err != nil {
		if _, ok := err.(*api.NoSuchContainer); ok == false {
			log.Errorf(err.Error())
		}
	}
	return n
}
func (cm *Docker) inspectService(id string) *swarm.Service {
	s, err := cm.client.InspectService(id)
	if err != nil {
		if _, ok := err.(*api.NoSuchService); ok == false {
			log.Errorf(err.Error())
		}
	}
	return s
}
func (cm *Docker) inspectTask(id string) *swarm.Task {
	s, err := cm.client.InspectTask(id)
	if err != nil {
		if _, ok := err.(*api.NoSuchTask); ok == false {
			log.Errorf(err.Error())
		}
	}
	return s
}

// Mark all container IDs for refresh
func (cm *Docker) refreshAllContainers() {
	opts := api.ListContainersOptions{All: true}
	allContainers, err := cm.client.ListContainers(opts)
	if err != nil {
		panic(fmt.Sprintf("Refreshing all containers:%s", err))
	}

	for _, i := range allContainers {
		c := cm.MustGetContainer(i.ID)
		c.SetMeta("name", shortName(i.Names[0]))
		c.SetState(i.State)
		cm.needsRefreshContainers <- c.Id
	}
}

func (cm *Docker) refreshAllNodes() {
	ctx, cancel := context.WithCancel(context.Background())
	opt := api.ListNodesOptions{Context: ctx}
	allNodes, err := cm.client.ListNodes(opt)

	if err != nil {
		panic(fmt.Sprintf("Refreshing all nodes:%s", err))
	}
	for _, i := range allNodes {
		n := cm.MustGetNode(i.ID)
		n.SetMeta("name", i.Description.Hostname)
		cm.needsRefreshNodes <- n.Id
	}

	if cancel != nil {
		cancel()
	}
}

func (cm *Docker) refreshAllServices() {
	ctx, cancel := context.WithCancel(context.Background())
	opts := api.ListServicesOptions{Context: ctx}
	allServices, err := cm.client.ListServices(opts)

	if err != nil {
		panic(fmt.Sprintf("Refreshing all services:%s", err))
	}
	for _, i := range allServices {
		s := cm.MustGetService(i.ID)

		s.SetMeta("name", i.Spec.Annotations.Name)
		s.SetState("service")
		cm.needsRefreshServices <- s.Id
	}
	if cancel != nil {
		cancel()
	}
}

func (cm *Docker) refreshAllTasks() {
	ctx, cancel := context.WithCancel(context.Background())
	opt := api.ListTasksOptions{Context: ctx}
	allTasks, err := cm.client.ListTasks(opt)

	if err != nil {
		panic(fmt.Sprintf("Refreshing all tasks:%s", err))
	}
	for n, i := range allTasks {
		t := cm.MustGetTask(i.ID)

		node := cm.MustGetNode(i.NodeID)
		service := cm.MustGetService(i.ServiceID)
		t.SetMeta("name", "\\"+service.GetMeta("name")+"."+fmt.Sprintf("%d", n))
		t.SetMeta("node", node.GetMeta("name"))
		t.SetState(fmt.Sprintf("%s", i.Status.State))
		t.SetMeta("service", i.ServiceID)
		log.Debugf("Service %s Node id %s Node name %s", t.GetMeta("name"), i.NodeID, node.GetMeta("name"))
		cm.needsRefreshTasks <- t.Id
	}

	if cancel != nil {
		cancel()
	}
}

func (cm *Docker) LoopContainer() {
	for id := range cm.needsRefreshContainers {
		c := cm.MustGetContainer(id)
		cm.refreshContainer(c)
	}
}
func (cm *Docker) LoopNode() {
	for id := range cm.needsRefreshNodes {
		n := cm.MustGetNode(id)
		cm.refreshNode(n)
	}
}
func (cm *Docker) LoopService() {
	for id := range cm.needsRefreshServices {
		s := cm.MustGetService(id)
		cm.refreshService(s)
	}
}
func (cm *Docker) LoopTask() {
	for id := range cm.needsRefreshTasks {
		t := cm.MustGetTask(id)
		cm.refreshTask(t)
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

func (cm *Docker) MustGetTask(id string) *entity.Task {
	n, ok := cm.GetTask(id)
	if !ok {
		collector := collector.NewDocker(cm.client, id)
		n = entity.NewTask(id, collector)
		cm.lock.Lock()
		cm.tasks[id] = n
		cm.lock.Unlock()
	}
	return n
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

func (cm *Docker) GetTask(id string) (*entity.Task, bool) {
	cm.lock.Lock()
	t, ok := cm.tasks[id]
	cm.lock.Unlock()
	return t, ok
}

func (cm *Docker) GetNode(id string) (*entity.Node, bool) {
	cm.lock.Lock()
	n, ok := cm.nodes[id]
	cm.lock.Unlock()
	return n, ok
}

//del by id
func (cm *Docker) delByIDContainer(id string) {
	cm.lock.Lock()
	delete(cm.containers, id)
	cm.lock.Unlock()
	log.Infof("removed dead container: %s", id)
}

func (cm *Docker) delByIDNode(id string) {
	cm.lock.Lock()
	delete(cm.nodes, id)
	cm.lock.Unlock()
	log.Infof("removed node: %s", id)
}

func (cm *Docker) delByIDService(id string) {
	cm.lock.Lock()
	delete(cm.services, id)
	cm.lock.Unlock()
	log.Infof("removed stopped service: %s", id)
}

func (cm *Docker) delByIDTask(id string) {
	cm.lock.Lock()
	delete(cm.tasks, id)
	cm.lock.Unlock()
	log.Infof("removed task: %s", id)
}

func (cm *Docker) AllNodes() (nodes entity.Nodes) {
	cm.lock.Lock()
	for _, node := range cm.nodes {
		nodes = append(nodes, node)
	}
	//nodes.Sort()
	nodes.Filter()
	cm.lock.Unlock()
	return nodes
}

func (cm *Docker) AllServices() (services entity.Services) {
	cm.lock.Lock()
	for _, service := range cm.services {
		services = append(services, service)
	}

	services.Sort()
	services.Filter()
	cm.lock.Unlock()
	return services
}

func (cm *Docker) AllTasks() (tasks entity.Tasks) {
	cm.lock.Lock()
	for _, task := range cm.tasks {
		tasks = append(tasks, task)
	}

	tasks.Sort()
	tasks.Filter()
	cm.lock.Unlock()
	return tasks
}

func (cm *Docker) AllContainers() (containers entity.Containers) {
	cm.lock.Lock()
	for _, container := range cm.containers {
		containers = append(containers, container)
	}

	containers.Sort()
	containers.Filter()
	cm.lock.Unlock()
	return containers
}

// use primary container name
func shortName(name string) string {
	return strings.Replace(name, "/", "", 1)
}
