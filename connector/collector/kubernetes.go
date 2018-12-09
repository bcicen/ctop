package collector

import (
	"encoding/json"
	"time"

	clientset "k8s.io/metrics/pkg/client/clientset/versioned"

	"github.com/bcicen/ctop/config"
	"github.com/bcicen/ctop/models"

	"k8s.io/client-go/kubernetes"
)

// Kubernetes collector
type Kubernetes struct {
	models.Metrics
	name       string
	client     clientset.Interface
	clientset  *kubernetes.Clientset
	running    bool
	stream     chan models.Metrics
	done       chan bool
	lastCpu    float64
	lastSysCpu float64
	scaleCpu   bool
	interval   time.Duration
}

type Metric struct {
	Timestamp time.Time `json:"timestamp"`
	Value     int64     `json:"value"`
}

type Response struct {
	Metrics         []Metric  `json:"metrics"`
	LatestTimestamp time.Time `json:"latest_timestamp"`
}

func NewKubernetes(client *kubernetes.Clientset, name string) *Kubernetes {
	return &Kubernetes{
		Metrics:   models.Metrics{},
		name:      name,
		client:    clientset.New(client.RESTClient()),
		clientset: client,
		scaleCpu:  config.GetSwitchVal("scaleCpu"),
		interval:  time.Duration(30) * time.Second,
	}
}

func buildURL(namespace, podName string) string {
	return "/api/v1/namespaces/kube-system/services/heapster/proxy/api/v1/model/namespaces/" + namespace + "/pods/" + podName
}

func (k *Kubernetes) Start() {
	k.done = make(chan bool)
	k.stream = make(chan models.Metrics)

	go func() {
		k.running = false
		for {
			log.Debugf("collect k8s metrics %s\n", k.name)
			k.ReadCPU()
			k.ReadMem()
			k.ReadNetRx()
			k.ReadNetTx()
			k.stream <- k.Metrics
			time.Sleep(k.interval)
		}
	}()

	k.running = true
	log.Infof("collector started for container: %s", k.name)
}

func (c *Kubernetes) Running() bool {
	return c.running
}

func (c *Kubernetes) Stream() chan models.Metrics {
	return c.stream
}

func (c *Kubernetes) Logs() LogCollector {
	return NewKubernetesLogs(c.name, c.clientset)
}

// Stop collector
func (c *Kubernetes) Stop() {
	c.done <- true
}

func (k *Kubernetes) ReadCPU() {
	cpu, err := k.read("/cpu/usage_rate")

	if err != nil {
		log.Errorf("collecte network cpu metric has error %s here %s", k.name, err.Error())
		time.Sleep(1 * time.Second)
		return
	}

	// TODO: heapster returning usage CPU in micro values without point so 0.004 is 4
	// because k8s calculate percent usage of all available CPU in cluster
	if cpu != 0 {
		k.CPUUtil = round(float64(cpu))
	}
}

func (k *Kubernetes) ReadMem() {
	usage, err := k.read("/memory/usage")
	if err != nil {
		log.Errorf("collecte network memory metric has error %s here %s", k.name, err.Error())
		time.Sleep(1 * time.Second)
		return
	}
	cache, err := k.read("/memory/cache")
	if err != nil {
		log.Errorf("collecte network memory metric has error %s here %s", k.name, err.Error())
		time.Sleep(1 * time.Second)
		return
	}
	k.MemUsage = usage - cache

	limit, err := k.read("/memory/limit")
	if err != nil {
		log.Errorf("collecte network memory metric has error %s here %s", k.name, err.Error())
		time.Sleep(1 * time.Second)
		return
	}
	k.MemLimit = limit
	//k.MemPercent = percent(float64(k.MemUsage), float64(k.MemLimit))
}

func (k *Kubernetes) ReadNetRx() {
	rx, err := k.read("/network/rx_rate")
	if err != nil {
		log.Errorf("collecte network rx_rate metric has error %s here %s", k.name, err.Error())
		time.Sleep(1 * time.Second)
		return
	}
	k.NetRx = rx
}

func (k *Kubernetes) ReadNetTx() {
	tx, err := k.read("/network/tx_rate")
	if err != nil {
		log.Errorf("collecte network tx_rate metric has error %s here %s", k.name, err.Error())
		time.Sleep(1 * time.Second)
		return
	}
	k.NetTx = tx
}

func (k *Kubernetes) read(name string) (int64, error) {
	m := &Response{}
	url := buildURL(config.GetVal("namespace"), k.name) + "/metrics" + name
	log.Debugf("get metrics: %s", url)
	b, err := k.clientset.RESTClient().Get().AbsPath(url).Do().Raw()
	if err != nil {
		return 0, err
	}
	err = json.Unmarshal(b, m)
	if err != nil {
		return 0, err
	}
	return m.Metrics[len(m.Metrics)-1].Value, nil
}

//func (c *Kubernetes) ReadIO(stats *api.Stats) {
//	var read, write int64
//	for _, blk := range stats.BlkioStats.IOServiceBytesRecursive {
//		if blk.Op == "Read" {
//			read = int64(blk.Value)
//		}
//		if blk.Op == "Write" {
//			write = int64(blk.Value)
//		}
//	}
//	c.IOBytesRead, c.IOBytesWrite = read, write
//}
