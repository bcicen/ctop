package network

import (
	"net/http"
	"encoding/json"
	"github.com/bcicen/ctop/models"
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
