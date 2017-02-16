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
	params   map[string]*Param
	switches map[string]*Switch
	updates  chan ConfigMsg
}

type ConfigMsg struct {
	key string
	val string
}

func Update(k, v string) {
	Global.updates <- ConfigMsg{k, v}
}

func NewDefaultConfig() Config {
	config := Config{
		params:   make(map[string]*Param),
		switches: make(map[string]*Switch),
		updates:  make(chan ConfigMsg),
	}

	for _, p := range params {
		config.params[p.key] = p
		log.Debugf("loaded config param: \"%s\": \"%s\"", p.key, p.val)
	}

	for _, t := range switches {
		config.switches[t.key] = t
		log.Debugf("loaded config switch: \"%s\": %t", t.key, t.val)
	}

	go func() {
		for m := range config.updates {
			config.params[m.key].val = m.val
			log.Noticef("config change: %s: %s", m.key, m.val)
		}
	}()

	return config
}

// Return env var value if set, else return defaultVal
func getEnv(key, defaultVal string) string {
	val := os.Getenv(key)
	if val != "" {
		return val
	}
	return defaultVal
}
