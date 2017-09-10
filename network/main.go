package network

import (
	"net/http"
	"encoding/json"
	"fmt"
	"bytes"

	"github.com/bcicen/ctop/models"
	"github.com/bcicen/ctop/logging"
	"github.com/bcicen/ctop/config"
)

var (
	log = logging.Init()
)

func Main() {
	log.Infof("start HTTP server, listen :9001/metrics")
	http.HandleFunc("/metrics", Metrics)
	http.ListenAndServe(":9001", nil)
}

func Metrics(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		log.Infof("GET")
	case "POST":
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "Hello, POT method. ParseForm() err: %v", err)
			return
		}
		decoder := json.NewDecoder(r.Body)
		defer r.Body.Close()
		var metrics models.Metrics
		if err := decoder.Decode(&metrics); err != nil {
			log.Error(fmt.Sprintf("Can't decode Metrics %s", err))
		}
		log.Infof("POST Metrics: %s", metrics)
	default:
		log.Infof("Sorry, only GET and POST methods are supported.")
	}
}

func TestDockerNetwork(metric *models.Metrics) {
	log.Infof("send test docker")
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(metric)
	res, err := http.Post("http://"+config.GetVal("host")+":9001/metrics", "application/json; charset=utf-8", b)
	defer res.Body.Close()
	if err != nil {
		log.Error(fmt.Sprintf("Cna't POST %s", err))
		return
	}
	log.Infof("Response: %s", res.Body)
}
