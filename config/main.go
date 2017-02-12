package config

import (
	"os"

	"github.com/bcicen/ctop/logging"
)

var (
	Global = NewDefaultConfig()
	log    = logging.Init()
)

type Config struct {
	params  map[string]string
	toggles map[string]bool
	updates chan ConfigMsg
}

// Return param value
func Get(k string) string {
	if _, ok := Global.params[k]; ok == true {
		return Global.params[k]
	}
	return ""
}

// Return toggle value
func GetToggle(k string) bool {
	if _, ok := Global.toggles[k]; ok == true {
		return Global.toggles[k]
	}
	return false
}

// Toggle a boolean option
func Toggle(k string) {
	Global.toggles[k] = Global.toggles[k] != true
}

type ConfigMsg struct {
	key string
	val string
}

func Update(k, v string) {
	log.Noticef("config update: %s = %s", k, v)
	Global.updates <- ConfigMsg{k, v}
}

func NewDefaultConfig() Config {
	docker := os.Getenv("DOCKER_HOST")
	if docker == "" {
		docker = "unix:///var/run/docker.sock"
	}

	params := map[string]string{
		"dockerHost": docker,
		"filterStr":  "",
		"sortField":  "id",
	}

	toggles := map[string]bool{
		"sortReverse":    false,
		"enableHeader":   false,
		"loggingEnabled": true,
	}

	config := Config{params, toggles, make(chan ConfigMsg)}
	go func() {
		for m := range config.updates {
			config.params[m.key] = m.val
		}
	}()
	return config
}
