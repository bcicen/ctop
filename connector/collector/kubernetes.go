package collector

import (
	"encoding/json"
	"time"

	"k8s.io/metrics/pkg/apis/metrics/v1alpha1"
	clientset "k8s.io/metrics/pkg/client/clientset/versioned"

	"github.com/bcicen/ctop/config"
	"github.com/bcicen/ctop/models"
	"k8s.io/api/core/v1"

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
	}
}

func (k *Kubernetes) Start() {
	k.done = make(chan bool)
	k.stream = make(chan models.Metrics)

	go func() {
		k.running = false
		for {

			result := &v1alpha1.PodMetrics{}
			err := k.clientset.RESTClient().Get().AbsPath("/api/v1/namespaces/kube-system/services/http:heapster:/proxy/apis/metrics/v1alpha1/namespaces/" + config.GetVal("namespace") + "/pods/" + k.name).Do().Into(result)

			if err != nil {
				log.Errorf("has error %s here %s", k.name, err.Error())
				time.Sleep(1 * time.Second)
				continue
			}
			k.ReadCPU(result)
			k.ReadMem(result)
			txMetrics := &Response{}
			b, err := k.clientset.RESTClient().Get().AbsPath("/api/v1/namespaces/kube-system/services/heapster/proxy/api/v1/model/namespaces/" + config.GetVal("namespace") + "/pods/" + k.name + "/metrics/network/tx").Do().Raw()
			if err != nil {
				log.Errorf("has error %s here %s", k.name, err.Error())
				time.Sleep(1 * time.Second)
				continue
			}
			err = json.Unmarshal(b, txMetrics)
			if err != nil {
				log.Errorf("has error %s here %s", k.name, err.Error())
				time.Sleep(1 * time.Second)
				continue
			}
			rxMetrics := &Response{}
			b, err = k.clientset.RESTClient().Get().AbsPath("/api/v1/namespaces/kube-system/services/heapster/proxy/api/v1/model/namespaces/" + config.GetVal("namespace") + "/pods/" + k.name + "/metrics/network/rx").Do().Raw()
			if err != nil {
				log.Errorf("has error %s here %s", k.name, err.Error())
				time.Sleep(1 * time.Second)
				continue
			}
			err = json.Unmarshal(b, rxMetrics)
			if err != nil {
				log.Errorf("has error %s here %s", k.name, err.Error())
				time.Sleep(1 * time.Second)
				continue
			}
			k.ReadNet(rxMetrics, txMetrics)
			k.stream <- k.Metrics
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

func (k *Kubernetes) ReadCPU(metrics *v1alpha1.PodMetrics) {
	all := int64(0)
	for _, c := range metrics.Containers {
		v := c.Usage[v1.ResourceCPU]
		all += v.Value()
	}
	if all != 0 {
		k.CPUUtil = round(float64(all))
	}
}

func (k *Kubernetes) ReadMem(metrics *v1alpha1.PodMetrics) {
	all := int64(0)
	for _, c := range metrics.Containers {
		v := c.Usage[v1.ResourceMemory]
		a, ok := v.AsInt64()
		if ok {
			all += a
		}
	}
	k.MemUsage = all
	k.MemLimit = int64(0)
	//k.MemPercent = percent(float64(k.MemUsage), float64(k.MemLimit))
}

func (k *Kubernetes) ReadNet(rxR, txR *Response) {
	var rx, tx int64
	for _, network := range rxR.Metrics {
		rx += int64(network.Value)
	}
	for _, network := range txR.Metrics {
		tx += int64(network.Value)
	}
	k.NetRx, k.NetTx = rx, tx
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
