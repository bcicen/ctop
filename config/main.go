package config

import (
	"os"

	"github.com/bcicen/ctop/logging"
)

var (
	Global     = NewDefaultConfig()
	log        = logging.Init()
	configChan = make(chan ConfigMsg)
)

type Config map[string]string

type ConfigMsg struct {
	key string
	val string
}

func Update(k, v string) {
	log.Noticef("config update: %s = %s", k, v)
	configChan <- ConfigMsg{k, v}
}

// Toggle a boolean option
func Toggle(k string) {
	if Global[k] == "0" {
		Global[k] = "1"
	} else {
		Global[k] = "0"
	}
}

func NewDefaultConfig() Config {
	docker := os.Getenv("DOCKER_HOST")
	if docker == "" {
		docker = "unix:///var/run/docker.sock"
	}
	config := Config{
		"dockerHost":     docker,
		"filterStr":      "",
		"sortField":      "id",
		"sortReverse":    "0",
		"enableHeader":   "0",
		"loggingEnabled": "1",
	}
	go func() {
		for m := range configChan {
			config[m.key] = m.val
		}
	}()
	return config
}
