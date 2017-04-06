package config

import (
	"fmt"
	"os"

	"github.com/bcicen/ctop/logging"
)

var Config struct {
	Params
	Switches
}

var (
	log = logging.Init()
)

func Init() {
	initParams()
	initSwitches()
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
