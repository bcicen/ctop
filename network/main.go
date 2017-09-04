package network

import (
	"net/http"
	"encoding/json"
	"github.com/bcicen/ctop/models"
	"github.com/bcicen/ctop/entity"
	"github.com/bcicen/ctop/logging"
	"fmt"
	"time"
	//"strings"
)

var (
	log = logging.Init()
)

func Main() {
	http.HandleFunc("/metrics", Metrics)
	http.ListenAndServe(":9001", nil)
}

func Metrics(w http.ResponseWriter, r *http.Request) {
	metric := &models.Metrics{CPUUtil: 1,
		NetTx: 1,
		NetRx: 1,
		MemLimit: 2,
		MemPercent: 2,
		MemUsage: 2,
		IOBytesRead: 3,
		IOBytesWrite: 3,
		Pids: 4}
	js, err := json.Marshal(metric)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func TestDockerNetwork(services map[string]*entity.Task) {
	for {
		for k, v := range services {
			log.Debugf("Get service id " + k)
			//if !strings.Contains(v.GetMeta("name"), "CTOP_swarm"){
			//	continue
			//}
			url := "http://" + v.GetMeta("addr") + ":9001/metrics"
			r, err := http.Get(url)
			if err != nil {
				log.Error(fmt.Sprintf("Can't HTTP-GET %s, with error: %s", url, err))
			}
			log.Infof("Get response from %s. Response: \n %s", url, r)
		}
		time.Sleep(5 * time.Second)
	}
}
