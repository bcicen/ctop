package config

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/BurntSushi/toml"
)

var (
	xdgRe = regexp.MustCompile("^XDG_*")
)

type File struct {
	Options map[string]string `toml:"options"`
	Toggles map[string]bool   `toml:"toggles"`
}

func exportConfig() File {
	// update columns param from working config
	Update("columns", ColumnsString())

	lock.RLock()
	defer lock.RUnlock()

	c := File{
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

//
func Read() error {
	var config File

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

	// set working column config, if provided
	colStr := GetVal("columns")
	if len(colStr) > 0 {
		var colNames []string
		for _, s := range strings.Split(colStr, ",") {
			s = strings.TrimSpace(s)
			if s != "" {
				colNames = append(colNames, s)
			}
		}
		SetColumns(colNames)
	}

	return nil
}

func Write() (path string, err error) {
	path, err = getConfigPath()
	if err != nil {
		return path, err
	}

	cfgdir := filepath.Dir(path)
	// create config dir if not exist
	if _, err := os.Stat(cfgdir); err != nil {
		err = os.MkdirAll(cfgdir, 0755)
		if err != nil {
			return path, fmt.Errorf("failed to create config dir [%s]: %s", cfgdir, err)
		}
	}

	// remove prior to writing new file
	if err := os.Remove(path); err != nil {
		if !os.IsNotExist(err) {
			return path, err
		}
	}

	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return path, fmt.Errorf("failed to open config for writing: %s", err)
	}

	writer := toml.NewEncoder(file)
	err = writer.Encode(exportConfig())
	if err != nil {
		return path, fmt.Errorf("failed to write config: %s", err)
	}

	return path, nil
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
