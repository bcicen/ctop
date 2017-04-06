package config

import (
	"reflect"
)

// defaults
type Switches struct {
	SortReversed  bool // Reverse Sort Order
	AllContainers bool // Show All Containers
	EnableHeader  bool // Enable Status Header
}

func initSwitches() {
	Config.SortReversed = false
	Config.AllContainers = true
	Config.EnableHeader = true
}

// Get Param by key
func GetSwitch(k string) reflect.Value {
	return reflect.ValueOf(&Config).Elem().FieldByName(k)
}

// Get Param value by key
func GetSwitchVal(k string) bool {
	p := Get(k)
	if p.Kind().String() != "bool" {
		log.Errorf("Tried to access a " + p.Kind().String() + " named " + k + " as a Bool Param")
		return false
	}
	return p.Bool()
}

// Set param value
func Toggle(k string) {
	p := Get(k)
	if p.CanSet() && p.Kind().String() == "bool" {
		log.Noticef("config change: %s: %b -> %b", k, p.Bool(), p.Bool() != true)
		p.SetBool(p.Bool() != true)
	} else {
		log.Errorf("ignoring toggle for non-existent switch: %s", k)
	}
}
