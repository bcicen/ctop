package config

var switches = []*Switch{
	&Switch{
		key:   "sortReverse",
		val:   false,
		label: "Reverse Sort Order",
	},
	&Switch{
		key:   "allContainers",
		val:   false,
		label: "Show All Containers",
	},
	&Switch{
		key:   "enableHeader",
		val:   false,
		label: "Enable cTop Status Line",
	},
	&Switch{
		key:   "loggingEnabled",
		val:   true,
		label: "Enable Logging Server",
	},
}

type Switch struct {
	key   string
	val   bool
	label string
}

// Return toggle value
func GetSwitch(k string) bool {
	if _, ok := Global.switches[k]; ok == true {
		return Global.switches[k].val
	}
	return false // default
}

// Toggle a boolean switch
func Toggle(k string) {
	Global.switches[k].val = Global.switches[k].val != true
	log.Noticef("config change: %s: %t", k, Global.switches[k].val)
}
