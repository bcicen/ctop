package config

// defaults
var switches = []*Switch{
	&Switch{
		Key:   "sortReversed",
		Val:   false,
		Label: "Reverse Sort Order",
	},
	&Switch{
		Key:   "allContainers",
		Val:   true,
		Label: "Show All Containers",
	},
	&Switch{
		Key:   "enableHeader",
		Val:   true,
		Label: "Enable Status Header",
	},
	&Switch{
		Key:   "scaleCpu",
		Val:   false,
		Label: "Show CPU as %% of system total",
	},
}

type Switch struct {
	Key   string
	Val   bool
	Label string
}

// Return Switch by key
func GetSwitch(k string) *Switch {
	for _, sw := range GlobalSwitches {
		if sw.Key == k {
			return sw
		}
	}
	return &Switch{} // default
}

// Return Switch value by key
func GetSwitchVal(k string) bool {
	return GetSwitch(k).Val
}

// Toggle a boolean switch
func Toggle(k string) {
	sw := GetSwitch(k)
	newVal := sw.Val != true
	log.Noticef("config change: %s: %t -> %t", k, sw.Val, newVal)
	sw.Val = newVal
	//log.Errorf("ignoring toggle for non-existant switch: %s", k)
}
