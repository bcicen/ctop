package collector

import (
	"k8s.io/metrics/pkg/client/clientset_generated/clientset"

	"github.com/bcicen/ctop/config"
	"github.com/bcicen/ctop/models"
	"k8s.io/client-go/kubernetes"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

			cm, err := k.client.Metrics().PodMetricses("akozlenkov").List(metav1.ListOptions{})
			if err != nil {
				log.Errorf(">>>>>> %s here %s", k.name, err.Error())
				continue
			}
			log.Debugf(">>>> %+v", cm)
			//for _, m := range cm.Containers {
			//	log.Debugf(">>>> %+v", m)
			//}
		}
	}()

	//go func() {
	//	defer close(c.stream)
	//	for s := range stats {
	//		c.ReadCPU(s)
	//		c.ReadMem(s)
	//		c.ReadNet(s)
	//		c.ReadIO(s)
	//		c.stream <- c.Metrics
	//	}
	//	log.Infof("collector stopped for container: %s", c.id)
	//}()

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

//
//func (c *Kubernetes) ReadCPU(stats *api.Stats) {
//	ncpus := float64(len(stats.CPUStats.CPUUsage.PercpuUsage))
//	total := float64(stats.CPUStats.CPUUsage.TotalUsage)
//	system := float64(stats.CPUStats.SystemCPUUsage)
//
//	cpudiff := total - c.lastCpu
//	syscpudiff := system - c.lastSysCpu
//
//	if c.scaleCpu {
//		c.CPUUtil = round((cpudiff / syscpudiff * 100))
//	} else {
//		c.CPUUtil = round((cpudiff / syscpudiff * 100) * ncpus)
//	}
//	c.lastCpu = total
//	c.lastSysCpu = system
//	c.Pids = int(stats.PidsStats.Current)
//}

//func (c *Kubernetes) ReadMem(stats *api.Stats) {
//	c.MemUsage = int64(stats.MemoryStats.Usage - stats.MemoryStats.Stats.Cache)
//	c.MemLimit = int64(stats.MemoryStats.Limit)
//	c.MemPercent = percent(float64(c.MemUsage), float64(c.MemLimit))
//}

//func (c *Kubernetes) ReadNet(stats *api.Stats) {
//	var rx, tx int64
//	for _, network := range stats.Networks {
//		rx += int64(network.RxBytes)
//		tx += int64(network.TxBytes)
//	}
//	c.NetRx, c.NetTx = rx, tx
//}
//
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
