package config

// defaults
var switches = []*Switch{
	&Switch{
		key:   "sortReversed",
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
	for _, sw := range GlobalSwitches {
		if sw.key == k {
			return sw.val
		}
	}
	return false // default
}

// Toggle a boolean switch
func Toggle(k string) {
	for _, sw := range GlobalSwitches {
		if sw.key == k {
			newVal := sw.val != true
			log.Noticef("config change: %s: %t -> %t", k, sw.val, newVal)
			sw.val = newVal
			return
		}
	}
	log.Errorf("ignoring toggle for non-existant switch: %s", k)
}
