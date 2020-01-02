package config

import (
	"fmt"
	"os"
	"sync"

	"github.com/bcicen/ctop/logging"
)

var (
	GlobalParams   []*Param
	GlobalSwitches []*Switch
	GlobalWidgets  []*Widget
	lock           sync.RWMutex
	log            = logging.Init()
)

func Init() {
	for _, p := range defaultParams {
		GlobalParams = append(GlobalParams, p)
		log.Infof("loaded default config param: %s: %s", quote(p.Key), quote(p.Val))
	}
	for _, s := range defaultSwitches {
		GlobalSwitches = append(GlobalSwitches, s)
		log.Infof("loaded default config switch: %s: %t", quote(s.Key), s.Val)
	}
	for _, w := range defaultWidgets {
		GlobalWidgets = append(GlobalWidgets, w)
		log.Infof("loaded default widget: %s: %t", quote(w.Name), w.Enabled)
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
