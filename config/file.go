package config

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/BurntSushi/toml"
)

var (
	xdgRe = regexp.MustCompile("^XDG_*")
)

type ConfigFile struct {
	Options map[string]string `toml:"options"`
	Toggles map[string]bool   `toml:"toggles"`
}

func exportConfig() ConfigFile {
	c := ConfigFile{
		Options: make(map[string]string),
		Toggles: make(map[string]bool),
	}
	for _, p := range GlobalParams {
		c.Options[p.Key] = p.Val
	}
	for _, sw := range GlobalSwitches {
		c.Toggles[sw.Key] = sw.Val
	}
	return c
}

func Read() error {
	var config ConfigFile

	path, err := getConfigPath()
	if err != nil {
		return err
	}

	if _, err := toml.DecodeFile(path, &config); err != nil {
		return err
	}

	for k, v := range config.Options {
		Update(k, v)
	}
	for k, v := range config.Toggles {
		UpdateSwitch(k, v)
	}
	return nil
}

func Write() error {
	path, err := getConfigPath()
	if err != nil {
		return err
	}

	cfgdir := basedir(path)
	// create config dir if not exist
	if _, err := os.Stat(cfgdir); err != nil {
		err = os.MkdirAll(cfgdir, 0755)
		if err != nil {
			return fmt.Errorf("failed to create config dir [%s]: %s", cfgdir, err)
		}
	}

	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("failed to open config for writing: %s", err)
	}

	writer := toml.NewEncoder(file)
	err = writer.Encode(exportConfig())
	if err != nil {
		return fmt.Errorf("failed to write config: %s", err)
	}

	return nil
}

// determine config path from environment
func getConfigPath() (path string, err error) {
	homeDir, ok := os.LookupEnv("HOME")
	if !ok {
		return path, fmt.Errorf("$HOME not set")
	}

	// use xdg config home if possible
	if xdgSupport() {
		xdgHome, ok := os.LookupEnv("XDG_CONFIG_HOME")
		if !ok {
			xdgHome = fmt.Sprintf("%s/.config", homeDir)
		}
		path = fmt.Sprintf("%s/ctop/config", xdgHome)
	} else {
		path = fmt.Sprintf("%s/.ctop", homeDir)
	}

	return path, nil
}

// test for environemnt supporting XDG spec
func xdgSupport() bool {
	for _, e := range os.Environ() {
		if xdgRe.FindAllString(e, 1) != nil {
			return true
		}
	}
	return false
}

func basedir(path string) string {
	parts := strings.Split(path, "/")
	return strings.Join((parts[0 : len(parts)-1]), "/")
}
