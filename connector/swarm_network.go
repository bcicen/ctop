package connector

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"io/ioutil"

	"github.com/bcicen/ctop/config"
	"github.com/bcicen/ctop/models"
)

var (
	conn      Connector
	DoneServe = make(chan bool)
)

func Serve(current Connector) {
	defer close(DoneServe)
	conn = current
	server := &http.Server{Addr: ":9001", Handler: nil}
	server.SetKeepAlivesEnabled(false)
	http.HandleFunc("/metrics", Metrics)
	log.Infof("start HTTP server, listen :9001/metrics")
	go func() {
		select {
		case <-DoneServe:
			server.Shutdown(context.Background())
		}
	}()
	server.ListenAndServe()
}

func Metrics(w http.ResponseWriter, r *http.Request) {
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
