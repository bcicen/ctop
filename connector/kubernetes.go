package connector

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/bcicen/ctop/connector/collector"
	"github.com/bcicen/ctop/connector/manager"
	"github.com/bcicen/ctop/container"
	api "github.com/fsouza/go-dockerclient"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func init() { enabled["kubernetes"] = NewKubernetes }

var namespace = "akozlenkov"

type Kubernetes struct {
	clientset    *kubernetes.Clientset
	containers   map[string]*container.Container
	needsRefresh chan string // container IDs requiring refresh
	lock         sync.RWMutex
}

func NewKubernetes() Connector {
	var kubeconfig string
	//if home := homeDir(); home != "" {
	//	kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	//} else {
	//	kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	//}
	//flag.Parse()
	kubeconfig = filepath.Join(homeDir(), ".kube", "config")

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Error(err.Error())
		return nil
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Error(err.Error())
		return nil
	}

	// init docker client
	k := &Kubernetes{
		clientset:    clientset,
		containers:   make(map[string]*container.Container),
		needsRefresh: make(chan string, 60),
		lock:         sync.RWMutex{},
	}
	go k.Loop()
	k.refreshAll()
	go k.watchEvents()
	return k
}

func (k *Kubernetes) watchEvents() {
	for {
		log.Info("kubernetes event listener starting")
		allEvents, err := k.clientset.CoreV1().Events(namespace).List(metav1.ListOptions{})
		if err != nil {
			log.Error(err.Error())
			return
		}

		for _, e := range allEvents.Items {
			if e.Kind != "pod" {
				continue
			}

			actionName := strings.Split(e.Action, ":")[0]

			switch actionName {
			case "start", "die", "pause", "unpause", "health_status":
				log.Debugf("handling docker event: action=%s id=%s", e.Action, e.UID)
				k.needsRefresh <- e.Name
			case "destroy":
				log.Debugf("handling docker event: action=%s id=%s", e.Action, e.UID)
				k.delByID(e.Name)
			default:
				log.Debugf("handling docker event: %v", e)
				k.needsRefresh <- e.Name
			}
		}
		time.Sleep(1 * time.Second)
	}
}
func (k *Kubernetes) Loop() {
	for id := range k.needsRefresh {
		c := k.MustGet(id)
		k.refresh(c)
	}
}

// Get a single container, creating one anew if not existing
func (k *Kubernetes) MustGet(name string) *container.Container {
	c, ok := k.Get(name)
	// append container struct for new containers
	if !ok {
		// create collector
		collector := collector.NewKubernetes(k.clientset, name)
		// create manager
		manager := manager.NewKubernetes(k.clientset, name)
		// create container
		c = container.New(name, collector, manager)
		k.lock.Lock()
		k.containers[name] = c
		k.lock.Unlock()
	}
	return c
}

func (k *Kubernetes) refresh(c *container.Container) {
	insp := k.inspect(c.Id)
	// remove container if no longer exists
	if insp == nil {
		k.delByID(c.Id)
		return
	}
	c.SetMeta("name", insp.Name)
	if len(insp.Spec.Containers) >= 1 {
		c.SetMeta("image", insp.Spec.Containers[0].Image)
		c.SetMeta("ports", k8sPort(insp.Spec.Containers[0].Ports))
		for _, env := range insp.Spec.Containers[0].Env {
			c.SetMeta("[ENV-VAR]", env.Name+"="+env.Value)
		}
	}
	c.SetMeta("IPs", insp.Status.PodIP)
	c.SetMeta("created", insp.CreationTimestamp.Format("Mon Jan 2 15:04:05 2006"))
	c.SetMeta("health", string(insp.Status.Phase))
	c.SetState("running")
}

func k8sPort(ports []v1.ContainerPort) string {
	str := []string{}
	for _, p := range ports {
		str = append(str, fmt.Sprintf("%s:%d -> %d", p.HostIP, p.HostPort, p.ContainerPort))
	}
	return strings.Join(str, "\n")
}

func (k *Kubernetes) inspect(id string) *v1.Pod {
	p, err := k.clientset.CoreV1().Pods(namespace).Get(id, metav1.GetOptions{})
	if err != nil {
		if _, ok := err.(*api.NoSuchContainer); !ok {
			log.Errorf(err.Error())
		}
	}
	return p
}

// Remove containers by ID
func (k *Kubernetes) delByID(name string) {
	k.lock.Lock()
	delete(k.containers, name)
	k.lock.Unlock()
	log.Infof("removed dead container: %s", name)
}

func (k *Kubernetes) Get(name string) (c *container.Container, ok bool) {
	k.lock.Lock()
	c, ok = k.containers[name]
	k.lock.Unlock()
	return
}

// Mark all container IDs for refresh
func (k *Kubernetes) refreshAll() {
	allPods, err := k.clientset.CoreV1().Pods(namespace).List(metav1.ListOptions{})
	if err != nil {
		log.Error(err.Error())
		return
	}

	for _, pod := range allPods.Items {
		c := k.MustGet(pod.Name)
		c.SetMeta("uid", string(pod.UID))
		c.SetMeta("name", pod.Name)
		if pod.Initializers != nil && pod.Initializers.Result != nil {
			c.SetState(pod.Initializers.Result.Status)
		} else {
			c.SetState(string(pod.Status.Phase))
		}
		k.needsRefresh <- c.Id
	}
}

func (k *Kubernetes) All() (containers container.Containers) {
	k.lock.Lock()
	for _, c := range k.containers {
		containers = append(containers, c)
	}

	containers.Sort()
	containers.Filter()
	k.lock.Unlock()
	return containers
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
