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
	GlobalColumns  []*Column
	lock           sync.RWMutex
	log            = logging.Init()
)

func Init() {
	for _, p := range defaultParams {
		GlobalParams = append(GlobalParams, p)
		log.Infof("loaded default config param [%s]: %s", quote(p.Key), quote(p.Val))
	}
	for _, s := range defaultSwitches {
		GlobalSwitches = append(GlobalSwitches, s)
		log.Infof("loaded default config switch [%s]: %t", quote(s.Key), s.Val)
	}
	for _, c := range defaultColumns {
		x := c
		GlobalColumns = append(GlobalColumns, &x)
		log.Infof("loaded default widget config [%s]: %t", quote(x.Name), x.Enabled)
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
