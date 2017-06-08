package config

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/bcicen/ctop/logging"
	"io/ioutil"
	"os"
)

var configFile = os.Getenv("HOME") + "/.ctop.conf"

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
	loadConfig()
}

func loadConfig() {
	cfile, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Noticef("Config File %s does not appear to exist", configFile)
		return
	}
	_, err = toml.Decode(string(cfile), &Config)
	if err != nil {
		log.Noticef("Config File %s does appears to have errors: %s", configFile, err.Error())
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
