package connector

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"runtime"
	"time"

	"io/ioutil"

	"github.com/bcicen/ctop/config"
	"github.com/bcicen/ctop/models"
	"github.com/bcicen/ctop/widgets"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/swarm"
)

const (
	replicas = "replicas"
	global   = "global"
)

var (
	conn Connector
	// DoneServe chan bool for stopped HTTP server
	DoneServe = make(chan bool)
)

func serve(current Connector) {
	defer close(DoneServe)
	conn = current
	server := &http.Server{Addr: ":9001", Handler: nil}
	server.SetKeepAlivesEnabled(false)
	http.HandleFunc("/metrics", metrics)
	log.Infof("start HTTP server, listen :9001/metrics")
	go func() {
		select {
		case <-DoneServe:
			server.Shutdown(context.Background())
		}
	}()
	server.ListenAndServe()
}

func metrics(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		keys, ok := r.URL.Query()["id"]
		if !ok || len(keys) < 1 {
			log.Errorf("Url Param 'id' is missing")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		log.Debugf("Get metrics for task id %s", keys[0])
		if keys[0] == "" {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		docker, ok := conn.(*Docker)
		metrics, ok := docker.GetTaskMetrics(keys[0])
		if !ok {
			log.Errorf("Not found task with id %s", keys[0])
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		b := new(bytes.Buffer)
		json.NewEncoder(b).Encode(metrics)

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		status, err := w.Write(b.Bytes())
		if err != nil {
			log.Error("Cant write response: %d, err: %s", status, err.Error())
			return
		}
		log.Debugf("Successful response: %d", status)

	case "POST":
		defer r.Body.Close()
		bytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Fprintf(w, "Hello, POST method. ParseForm() err: %s", err)
			return
		}
		var metrics models.Metrics
		err = json.Unmarshal(bytes, &metrics)
		if err != nil {
			log.Errorf("Can't decode Metrics %s", err)
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			conn.SetMetrics(metrics)
			log.Infof("POST Metrics: %s", metrics)
			w.WriteHeader(http.StatusOK)
		}
	default:
		if len(config.GetVal("host")) == 0 {
			log.Infof("Sorry, only GET methods are supported.")
			return
		}
		log.Infof("Sorry, only POST methods are supported.")
	}
}

func (cm *Docker) collectMetrics() {
	ips, err := net.LookupIP("tasks." + ctopSwarm)
	if err != nil {
		log.Errorf("Errors: %s", err.Error())
	}
	for id := range cm.tasks {
		if len(id) == 0 {
			continue
		}
		for _, ip := range ips {
			go func() {
				if id == "" {
					return
				}
				log.Debugf("ip: %s, task: %s", ip, id)
				resp, err := http.Get(fmt.Sprintf("http://%s:9001/metrics?id=%s", ip, id))
				if err != nil {
					log.Debugf("bad request: %s", err.Error())
					return
				}
				if resp.StatusCode != http.StatusOK {
					log.Debugf("Bad status: %d", resp.StatusCode)
					return
				}
				defer resp.Body.Close()
				bytes, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					log.Debugf("Cannon read bytes from body %s", err.Error())
				}
				var metrics models.Metrics
				err = json.Unmarshal(bytes, &metrics)
				if err != nil {
					log.Debugf("Can't decode Metrics %s", err.Error())
					return
				}
				conn.SetMetrics(metrics)
				log.Debugf("Metrics: %+v", metrics)
			}()
		}
	}
}

func (cm *Docker) swarmListen() {
	networks, err := cm.client.NetworkList(cm.currentContext, types.NetworkListOptions{})
	if err != nil {
		log.Errorf("Can't load list networks: %s", err.Error())
	}

	cm.networkSwarmID = ""
	for _, n := range networks {
		if n.Name == ctopNetwork {
			cm.networkSwarmID = n.ID
		}
	}

	if cm.networkSwarmID == "" {
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
		cm.networkSwarmID = net.ID
		log.Noticef("Create '%s' network: %s", ctopNetwork, net)
	}

	netConfig := swarm.NetworkAttachmentConfig{
		Target:  cm.networkSwarmID,
		Aliases: []string{"ctop"},
	}
	command := []string{"/ctop", "-D"}
	if len(config.GetVal("host")) > 0 {
		command = []string{"/ctop", "-D", "-host", config.GetVal("host")}
	}
	cont := swarm.ContainerSpec{
		Image: config.GetVal("image"),
		Mounts: []mount.Mount{
			{
				Type:     mount.TypeBind,
				Source:   "/var/run/docker.sock",
				Target:   "/var/run/docker.sock",
				ReadOnly: true,
			},
		},
		Env:     []string{"CTOP_DEBUG=1", "CTOP_DEBUG_TCP=1"},
		Command: command,
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
			break
		}
	}
	log.Debugf(fmt.Sprintf("Container %s", containerID))

	err = cm.client.NetworkConnect(cm.currentContext, cm.networkSwarmID, containerID, &network.EndpointSettings{
		NetworkID: cm.networkSwarmID,
	})
	if err != nil {
		log.Error(fmt.Sprintf("Can't connect to \n %s \n, with err:\n %s", cm.networkSwarmID, err))
	}
	go func() {
		for {
			cm.collectMetrics()
			time.Sleep(time.Millisecond)
		}
	}()
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

// replace swarm.SericeMode with string constant
func modeService(mode swarm.ServiceMode) string {
	if mode.Global == nil {
		return replicas
	}
	return global
}

func (cm *Docker) stopSwarm() {
	cm.doneNode <- true
	cm.doneService <- true
	cm.doneTask <- true
	cm.doneDiscovery <- true
	DoneServe <- true
	go cm.LoopContainer()
	cm.refreshAllContainers()
}

// LoopDiscoveryTasks loop for discovery tasks
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
		runtime.Gosched()
	}
}
