package connector

import (
	"context"
	"fmt"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/bcicen/ctop/config"
	"github.com/bcicen/ctop/connector/collector"
	"github.com/bcicen/ctop/connector/manager"
	"github.com/bcicen/ctop/entity"
	"github.com/bcicen/ctop/widgets"
	api "github.com/fsouza/go-dockerclient"

	"github.com/bcicen/ctop/models"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
)

var (
	ctopSwarm   = "CTOP_swarm"
	ctopNetwork = "ctop_default"
)

type Docker struct {
	client                 *client.Client
	containers             map[string]*entity.Container
	services               map[string]*entity.Service
	nodes                  map[string]*entity.Node
	tasks                  map[string]*entity.Task
	needsRefreshNodes      chan string // node IDs requiring refresh
	needsRefreshContainers chan string // service IDs requiring refresh
	needsRefreshTasks      chan string // task IDs requiring refresh
	needsRefreshServices   chan string // container IDs requiring refresh
	lock                   sync.RWMutex
	networkSwarmId         string
	currentContext         context.Context
	cancel                 context.CancelFunc
	// sync swarm channels
	doneNode      chan bool
	doneService   chan bool
	doneTask      chan bool
	doneDiscovery chan bool
}

func NewDocker() Connector {
	// init docker client
	client, err := client.NewEnvClient()
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
		networkSwarmId:         "",

		doneNode:      make(chan bool),
		doneService:   make(chan bool),
		doneTask:      make(chan bool),
		doneDiscovery: make(chan bool),
	}
	cm.currentContext, cm.cancel = context.WithCancel(context.Background())
	cm.checkLoadedSwarm()
	if config.GetSwitchVal("swarmMode") {
		go cm.SwarmListen()
		go cm.LoopNode()
		go cm.LoopService()
		go cm.LoopTask()
		go cm.LoopDiscoveryTasks()
		cm.refreshAllNodes()
		cm.refreshAllServices()
		go Serve(cm)
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
	messages, err := cm.client.Events(cm.currentContext, types.EventsOptions{})

	go func() {
		for e := range err {
			log.Errorf("%s", e)
		}
	}()
	go func() {
		for e := range messages {
			log.Debugf("Action ", e)
			if e.Type == "service" {
				actionName := strings.Split(e.Action, ":")[0]
				switch actionName {
				case "update":
					cm.needsRefreshServices <- e.ID
					log.Debugf("SERVICE. Action: %s, ID: %s", e.Action, e.ID)
				}
			} else if e.Type == "container" {
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
	}()
}

func portsFormat(container *types.ContainerJSON) string {
	var exposed []string
	var published []string

	for k, v := range container.NetworkSettings.Ports {
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
	c.SetMeta("ports", portsFormat(insp))
	c.SetMeta("created", insp.Created)
	if insp.State.Health != nil {
		c.SetMeta("health", insp.State.Health.Status)
	}
	c.SetState(insp.State.Status)
}

func (cm *Docker) refreshNode(n *entity.Node) {
	insp := cm.inspectNode(n.Id)
	if insp == nil {
		cm.delByIDNode(n.Id)
		return
	}
	n.SetMeta("name", insp.Description.Hostname)
}

func (cm *Docker) refreshService(s *entity.Service) {
	insp := cm.inspectService(s.Id)
	if insp == nil {
		cm.delByIDService(s.Id)
		return
	}
	s.SetMeta("name", insp.Spec.Annotations.Name)
}

func (cm *Docker) refreshTask(t *entity.Task) {
	insp := cm.inspectTask(t.Id)
	if insp == nil {
		log.Debugf("Delete task")
		cm.delByIDTask(t.Id)
		return
	}
	node := cm.MustGetNode(insp.NodeID)
	t.SetMeta("node", node.GetMeta("name"))
	t.SetState(fmt.Sprintf("%s", insp.Status.State))
	t.SetMeta("service", insp.ServiceID)
}

func (cm *Docker) inspectContainer(id string) *types.ContainerJSON {
	c, err := cm.client.ContainerInspect(context.Background(), id)
	if err != nil {
		if _, ok := err.(*api.NoSuchContainer); ok == false {
			log.Errorf(err.Error())
		}
	}
	return &c
}
func (cm *Docker) inspectNode(id string) *swarm.Node {
	n, _, err := cm.client.NodeInspectWithRaw(context.Background(), id)
	if err != nil {
		if _, ok := err.(*api.NoSuchContainer); ok == false {
			log.Errorf(err.Error())
		}
	}
	return &n
}
func (cm *Docker) inspectService(id string) *swarm.Service {
	s, _, err := cm.client.ServiceInspectWithRaw(context.Background(), id)
	if err != nil {
		if _, ok := err.(*api.NoSuchService); ok == false {
			log.Errorf(err.Error())
		}
	}
	return &s
}
func (cm *Docker) inspectTask(id string) *swarm.Task {
	s, _, err := cm.client.TaskInspectWithRaw(context.Background(), id)
	if err != nil {
		if _, ok := err.(*api.NoSuchTask); ok == false {
			log.Errorf(err.Error())
		}
	}
	return &s
}

// Mark all container IDs for refresh
func (cm *Docker) refreshAllContainers() {
	allContainers, err := cm.client.ContainerList(cm.currentContext, types.ContainerListOptions{All: true})
	if err != nil {
		panic(fmt.Sprintf("Refreshing all containers:%s", err))
	}

	for _, i := range allContainers {
		insp := cm.inspectContainer(i.ID)
		c := cm.MustGetContainer(i.ID)
		c.SetMeta("name", shortName(insp.Name))
		c.SetState(insp.State.Status)
		cm.needsRefreshContainers <- insp.ID
	}
}

// Mark all nodes IDs for refresh
func (cm *Docker) refreshAllNodes() {
	allNodes, err := cm.client.NodeList(cm.currentContext, types.NodeListOptions{})

	if err != nil {
		panic(fmt.Sprintf("Refreshing all nodes:%s", err))
	}
	for _, i := range allNodes {
		n := cm.MustGetNode(i.ID)
		n.SetMeta("name", i.Description.Hostname)
		cm.needsRefreshNodes <- n.Id
	}
}

// Mark all services IDs for refresh
func (cm *Docker) refreshAllServices() {
	allServices, err := cm.client.ServiceList(cm.currentContext, types.ServiceListOptions{})

	if err != nil {
		panic(fmt.Sprintf("Refreshing all services:%s", err))
	}
	for _, i := range allServices {
		s := cm.MustGetService(i.ID)

		s.SetMeta("name", i.Spec.Annotations.Name)
		s.SetMeta("mode", modeService(i.Spec.Mode))
		s.SetState("service")
		cm.needsRefreshServices <- s.Id
	}
}

// Mark all tasks IDs for refresh
func (cm *Docker) refreshAllTasks() {
	allTasks, err := cm.client.TaskList(cm.currentContext, types.TaskListOptions{})

	if err != nil {
		cm.DownSwarmMode()
		log.Error(fmt.Sprintf("Refreshing all tasks:%s", err))
		return
	}

	if len(allTasks) == 0 {
		cm.tasks = make(map[string]*entity.Task)
		return
	}

	for _, i := range allTasks {
		t := cm.MustGetTask(i.ID)

		node := cm.MustGetNode(i.NodeID)
		service := cm.MustGetService(i.ServiceID)
		taskState := i.Status.State
		if service.GetMeta("mode") == "replicas" {
			t.SetMeta("name", service.GetMeta("name")+"."+fmt.Sprintf("%d", i.Slot))
		} else {
			t.SetMeta("name", service.GetMeta("name")+"."+fmt.Sprintf("%s", i.NodeID))
		}
		t.SetMeta("node", node.GetMeta("name"))
		t.SetState(fmt.Sprintf("%s", taskState))
		if len(i.NetworksAttachments) > 0 && len(i.NetworksAttachments[0].Addresses) > 0 {
			t.SetMeta("addr", strings.Split(i.NetworksAttachments[0].Addresses[0], "/")[0])
		}
		t.SetMeta("service", i.ServiceID)
		cm.needsRefreshTasks <- t.Id
	}
}

// Loop for discovery tasks
func (cm *Docker) LoopDiscoveryTasks() {
	defer close(cm.doneDiscovery)
	for {
		select {
		case <-cm.doneDiscovery:
			return
		default:
			time.Sleep(time.Second)
			cm.refreshAllTasks()
		}
	}
}

// Loop for discovery container
func (cm *Docker) LoopContainer() {
	for id := range cm.needsRefreshContainers {
		c := cm.MustGetContainer(id)
		cm.refreshContainer(c)
	}
}

// Loop for discovery node
func (cm *Docker) LoopNode() {
	var id string
	defer close(cm.needsRefreshNodes)
	defer close(cm.doneNode)
	for {
		select {
		case id = <-cm.needsRefreshNodes:
			cm.refreshNode(cm.MustGetNode(id))
			break
		case <-cm.doneNode:
			return
		}
		runtime.Gosched()
	}
}

// Loop for discovery service
func (cm *Docker) LoopService() {
	var id string
	defer close(cm.needsRefreshServices)
	defer close(cm.doneService)
	for {
		select {
		case id = <-cm.needsRefreshServices:
			cm.refreshService(cm.MustGetService(id))
		case <-cm.doneService:
			return
		}
		runtime.Gosched()
	}
}

// Loop for discovery task
func (cm *Docker) LoopTask() {
	var id string
	defer close(cm.needsRefreshTasks)
	defer close(cm.doneTask)
	for {
		select {
		case id = <-cm.needsRefreshTasks:
			cm.refreshTask(cm.MustGetTask(id))
		case <-cm.doneTask:
			return
		}
		runtime.Gosched()
	}
}

// Get a single container, creating one anew if not existing
func (cm *Docker) MustGetContainer(id string) *entity.Container {
	c, ok := cm.GetContainer(id)
	// append container struct for new containers
	if !ok {
		// create collector
		collector := collector.NewDocker(cm.client, id)
		// create manager
		manager := manager.NewDocker(cm.client, id)
		// create container
		c = entity.NewContainer(id, collector, manager)
		cm.lock.Lock()
		cm.containers[id] = c
		cm.lock.Unlock()
	}
	return c
}

// Get a single service, creating one anew if not existing
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

// Get a single task, creating one anew if not existing
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

// Get a single node, creating one anew if not existing
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

// Get a single service, by ID
func (cm *Docker) GetService(id string) (*entity.Service, bool) {
	cm.lock.Lock()
	s, ok := cm.services[id]
	cm.lock.Unlock()
	return s, ok
}

// Get a single task, by ID
func (cm *Docker) GetTask(id string) (*entity.Task, bool) {
	cm.lock.Lock()
	t, ok := cm.tasks[id]
	cm.lock.Unlock()
	return t, ok
}

// Get a single node, by ID
func (cm *Docker) GetNode(id string) (*entity.Node, bool) {
	cm.lock.Lock()
	n, ok := cm.nodes[id]
	cm.lock.Unlock()
	return n, ok
}

//del container by id
func (cm *Docker) delByIDContainer(id string) {
	cm.lock.Lock()
	delete(cm.containers, id)
	cm.lock.Unlock()
	log.Infof("removed dead container: %s", id)
}

//del node by id
func (cm *Docker) delByIDNode(id string) {
	cm.lock.Lock()
	delete(cm.nodes, id)
	cm.lock.Unlock()
	log.Infof("removed node: %s", id)
}

//del service by id
func (cm *Docker) delByIDService(id string) {
	cm.lock.Lock()
	delete(cm.services, id)
	cm.lock.Unlock()
	log.Infof("removed stopped service: %s", id)
}

//del task by id
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

func modeService(mode swarm.ServiceMode) string {
	if mode.Global == nil {
		return "replicas"
	} else {
		return "global"
	}
}

func (cm *Docker) StopSwarm() {
	cm.doneNode <- true
	cm.doneService <- true
	cm.doneTask <- true
	cm.doneDiscovery <- true
	DoneServe <- true
	go cm.LoopContainer()
	cm.refreshAllContainers()
}

func (cm *Docker) checkLoadedSwarm() {
	if !config.GetSwitchVal("swarmMode") {
		return
	}

	filter := fmt.Sprintf(`{"name":{"%s":true}}`, ctopSwarm)
	args, err := filters.FromParam(filter)
	if err != nil {
		log.Errorf("Can't parser filter %s for finding service: %s", filter, err.Error())
	}
	ctopService := types.ServiceListOptions{
		Filters: args,
	}
	services, err := cm.client.ServiceList(cm.currentContext, ctopService)
	log.Debugf("Found services: %+v", services)
	if err != nil {
		log.Errorf("Can't find service: %s", err.Error())
	}
	if len(services) > 0 {
		return
	}
	widgets.ShowNotifiation()
}

func (cm *Docker) SwarmListen() {
	networks, err := cm.client.NetworkList(cm.currentContext, types.NetworkListOptions{})
	if err != nil {
		log.Errorf("Can't load list networks: %s", err.Error())
	}

	cm.networkSwarmId = ""
	for _, n := range networks {
		if n.Name == ctopNetwork {
			cm.networkSwarmId = n.ID
		}
	}

	if cm.networkSwarmId == "" {
		log.Infof(fmt.Sprintf("Netfowks: %s", networks))
		networkOpt := types.NetworkCreate{
			Driver:     "overlay",
			Attachable: true,
		}
		net, err := cm.client.NetworkCreate(cm.currentContext, ctopNetwork, networkOpt)
		if err != nil {
			log.Error(fmt.Sprintf("%s", err))
			return
		}
		cm.networkSwarmId = net.ID
		log.Noticef("Create '%s' network: %s", ctopNetwork, net)
	}

	netConfig := swarm.NetworkAttachmentConfig{
		Target:  cm.networkSwarmId,
		Aliases: []string{"ctop"},
	}
	cont := swarm.ContainerSpec{
		Image: config.Get("image").Val,
		Mounts: []mount.Mount{
			{
				Type:     mount.TypeBind,
				Source:   "/var/run/docker.sock",
				Target:   "/var/run/docker.sock",
				ReadOnly: true,
			},
		},
		Env:     []string{"CTOP_DEBUG=1", "CTOP_DEBUG_TCP=1"},
		Command: []string{"/ctop", "-D", "-host", config.GetVal("host")},
	}
	serviceSpec := swarm.ServiceSpec{
		Annotations: swarm.Annotations{Name: ctopSwarm, Labels: make(map[string]string)},
		TaskTemplate: swarm.TaskSpec{
			ContainerSpec: cont,
			Networks:      []swarm.NetworkAttachmentConfig{netConfig},
		},
		Networks: []swarm.NetworkAttachmentConfig{netConfig},
		Mode:     swarm.ServiceMode{Global: &swarm.GlobalService{}},
		UpdateConfig: &swarm.UpdateConfig{Parallelism: 1,
			Delay:           time.Duration(10),
			Monitor:         time.Duration(60),
			MaxFailureRatio: 0.5,
		},
		EndpointSpec: &swarm.EndpointSpec{
			Mode: swarm.ResolutionModeVIP,
			Ports: []swarm.PortConfig{
				{
					Name:          "tcp4",
					Protocol:      swarm.PortConfigProtocolTCP,
					TargetPort:    9000,
					PublishedPort: 9002,
					PublishMode:   swarm.PortConfigPublishModeHost,
				},
			},
		},
	}
	serviceOpt := types.ServiceCreateOptions{}
	_, err = cm.client.ServiceCreate(cm.currentContext, serviceSpec, serviceOpt)
	if err != nil {
		log.Error(fmt.Sprintf("Error create service:\n %s", err))
	}
	cm.refreshAllContainers()
	var containerID string
	for _, c := range cm.containers {
		if c.GetMeta("name") == "ctop" {
			containerID = c.GetId()
		}
	}
	log.Debugf(fmt.Sprintf("Container %s", containerID))

	err = cm.client.NetworkConnect(cm.currentContext, cm.networkSwarmId, containerID, &network.EndpointSettings{
		NetworkID: cm.networkSwarmId,
	})
	if err != nil {
		log.Error(fmt.Sprintf("Can't connect to \n %s \n, with err:\n %s", cm.networkSwarmId, err))
	}
}

func (cm *Docker) DownSwarmMode() {
	for _, s := range cm.services {
		if s.GetMeta("name") == "CTOP_swarm" {
			log.Infof("Down services.")
			cm.client.ServiceRemove(cm.currentContext, s.GetId())
			cm.client.NetworkRemove(cm.currentContext, cm.networkSwarmId)
		}
	}
	if cm.cancel != nil {
		cm.cancel()
	}
}

func (cm *Docker) SetMetrics(metrics models.Metrics) {
	if config.GetSwitchVal("swarmMode") {
		task := cm.MustGetTask(metrics.Id)
		task.SetMetrics(metrics)
	} else {
		cont := cm.MustGetContainer(metrics.Id)
		cont.SetMetrics(metrics)
	}
}
