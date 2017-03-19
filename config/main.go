package config

import (
	"fmt"
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
		log.Infof("loaded config param: %s: %s", quote(p.Key), quote(p.Val))
	}
	for _, s := range switches {
		GlobalSwitches = append(GlobalSwitches, s)
		log.Infof("loaded config switch: %s: %t", quote(s.Key), s.Val)
	}
}

func quote(s string) string {
	return fmt.Sprintf("\"%s\"", s)
}

// Return env var value if set, else return defaultVal
func getEnv(key, defaultVal string) string {
	val := os.Getenv(key)
	if val != "" {
		return val
	}
	return defaultVal
}
