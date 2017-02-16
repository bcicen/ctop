package config

import (
	"os"

	"github.com/bcicen/ctop/logging"
)

var (
	GlobalParams   []*Param
	GlobalSwitches []*Switch
	log            = logging.Init()
)

func Init() {
	for _, p := range params {
		GlobalParams = append(GlobalParams, p)
		log.Debugf("loaded config param: \"%s\": \"%s\"", p.key, p.val)
	}

	for _, s := range switches {
		GlobalSwitches = append(GlobalSwitches, s)
		log.Debugf("loaded config switch: \"%s\": %t", s.key, s.val)
	}
}

// Return env var value if set, else return defaultVal
func getEnv(key, defaultVal string) string {
	val := os.Getenv(key)
	if val != "" {
		return val
	}
	return defaultVal
}
