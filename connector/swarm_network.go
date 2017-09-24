package connector

import (
	"net/http"
	"encoding/json"
	"fmt"

	"github.com/bcicen/ctop/models"
	"io/ioutil"
)

var (
	conn Connector
)

func StartListen(current Connector) {
	conn = current
	server := &http.Server{Addr: ":9001", Handler: nil}
	server.SetKeepAlivesEnabled(false)
	http.HandleFunc("/metrics", Metrics)
	log.Infof("start HTTP server, listen :9001/metrics")
	server.ListenAndServe()
	//http.ListenAndServe(":9001", nil)
}

func Metrics(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		defer r.Body.Close()
		bytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Fprintf(w, "Hello, POT method. ParseForm() err: %s", err)
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
		log.Infof("Sorry, only POST methods are supported.")
	}
}
