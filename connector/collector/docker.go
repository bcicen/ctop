package collector

import (
	"github.com/bcicen/ctop/models"
	api "github.com/fsouza/go-dockerclient"
	"github.com/docker/docker/client"
	"context"
	"github.com/bcicen/ctop/config"
	"github.com/docker/docker/api/types"
	"io/ioutil"
	"encoding/json"
	"time"
	"bytes"
	"net/http"
	"strings"
)

// Docker collector
type Docker struct {
	models.Metrics
	id         string
	client     *client.Client
	running    bool
	stream     chan models.Metrics
	done       chan bool
	lastCpu    float64
	lastSysCpu float64
	httpClient http.Client
	url        string
}

func NewDocker(client *client.Client, id string) *Docker {
	return &Docker{
		Metrics:    models.Metrics{},
		id:         id,
		client:     client,
		httpClient: http.Client{},
		url:        "http://" + config.GetVal("host") + ":9001/metrics",
	}
}

type containerStats struct {
	stats   types.ContainerStats
	cont    types.ContainerJSON
	err     error
	errCont error
}

func (c *Docker) Start(id string) {
	if config.GetSwitchVal("swarmMode") {
		return
	}
	c.done = make(chan bool)
	c.stream = make(chan models.Metrics)
	stats := make(chan containerStats)
	go func() {
		ctx, closeCtx := context.WithCancel(context.Background())
		for {
			resp, err := c.client.ContainerStats(ctx, id, false)
			contJson, errCont := c.client.ContainerInspect(ctx, id)
			stats <- containerStats{resp, contJson, err, errCont}
			time.Sleep(time.Microsecond)
			if <-c.done {
				break
			}
		}
		defer close(stats)
		defer closeCtx()
		defer func() { c.running = false }()
	}()

	go func() {
		defer close(c.stream)
		for s := range stats {
			if s.err != nil {
				continue
				log.Errorf("%s", s.err)
			}
			b, _ := ioutil.ReadAll(s.stats.Body)
			s.stats.Body.Close()
			var apiStats api.Stats
			if err := json.Unmarshal(b, &apiStats); err != nil {
				log.Errorf("Unmarshal Stats error. err: %s", err)
			}
			log.Debugf("Api stats. %s", apiStats)
			c.ReadCPU(&apiStats)
			c.ReadMem(&apiStats)
			c.ReadNet(&apiStats)
			c.ReadIO(&apiStats)

			nameWords := strings.Split(s.cont.ContainerJSONBase.Name, ".")
			if len(nameWords) == 3 {
				c.Metrics.Id = nameWords[2]
			} else {
				c.Metrics.Id = id
			}
			c.done <- false
			c.stream <- c.Metrics
			go c.sendMetrics(&c.Metrics)
		}
		log.Infof("collector stopped for container: %s", c.id)
	}()

	c.running = true
	log.Infof("collector started for container: %s", c.id)
}

func (c *Docker) Running() bool {
	return c.running
}

func (c *Docker) Stream() chan models.Metrics {
	return c.stream
}

func (c *Docker) Logs() LogCollector {
	return &DockerLogs{c.id, c.client, make(chan bool)}
}

// Stop collector
func (c *Docker) Stop() {
	c.done <- true
}

func (c *Docker) ReadCPU(stats *api.Stats) {
	ncpus := float64(len(stats.CPUStats.CPUUsage.PercpuUsage))
	total := float64(stats.CPUStats.CPUUsage.TotalUsage)
	system := float64(stats.CPUStats.SystemCPUUsage)

	cpudiff := total - c.lastCpu
	syscpudiff := system - c.lastSysCpu

	c.CPUUtil = round((cpudiff / syscpudiff * 100) * ncpus)
	c.lastCpu = total
	c.lastSysCpu = system
	c.Pids = int(stats.PidsStats.Current)
}

func (c *Docker) ReadMem(stats *api.Stats) {
	c.MemUsage = int64(stats.MemoryStats.Usage - stats.MemoryStats.Stats.Cache)
	c.MemLimit = int64(stats.MemoryStats.Limit)
	c.MemPercent = percent(float64(c.MemUsage), float64(c.MemLimit))
}

func (c *Docker) ReadNet(stats *api.Stats) {
	var rx, tx int64
	for _, network := range stats.Networks {
		rx += int64(network.RxBytes)
		tx += int64(network.TxBytes)
	}
	c.NetRx, c.NetTx = rx, tx
}

func (c *Docker) ReadIO(stats *api.Stats) {
	var read, write int64
	for _, blk := range stats.BlkioStats.IOServiceBytesRecursive {
		if blk.Op == "Read" {
			read = int64(blk.Value)
		}
		if blk.Op == "Write" {
			write = int64(blk.Value)
		}
	}
	c.IOBytesRead, c.IOBytesWrite = read, write
}

func (c *Docker) sendMetrics(metric *models.Metrics) {
	if config.GetSwitchVal("enableDisplay") {
		return
	}
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(metric)
	req, err := http.NewRequest("POST", c.url, b)
	if err != nil {
		log.Errorf("%s", err)
	}
	req.Close = true
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	res, err := c.httpClient.Do(req)
	if err != nil {
		log.Errorf("%s", err)
		return
	}
	if res != nil {
		log.Debugf("Response: %s", res.Body)
	}
	defer res.Body.Close()
}
