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
		log.Infof("loaded config param: \"%s\": \"%s\"", p.Key, p.Val)
	}

	for _, s := range switches {
		GlobalSwitches = append(GlobalSwitches, s)
		log.Infof("loaded config switch: \"%s\": %t", s.Key, s.Val)
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
