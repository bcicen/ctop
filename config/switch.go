package config

// defaults
var switches = []*Switch{
	&Switch{
		Key:   "sortReversed",
		Val:   false,
		Label: "Reverse sort order",
	},
	&Switch{
		Key:   "allContainers",
		Val:   true,
		Label: "Show all containers",
	},
	&Switch{
		Key:   "fullRowCursor",
		Val:   true,
		Label: "Highlight entire cursor row (vs. name only)",
	},
	&Switch{
		Key:   "enableHeader",
		Val:   true,
		Label: "Enable status header",
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

func UpdateSwitch(k string, val bool) {
	sw := GetSwitch(k)
	if sw.Val != val {
		log.Noticef("config change: %s: %t -> %t", k, sw.Val, val)
		sw.Val = val
	}
}

// Toggle a boolean switch
func Toggle(k string) {
	sw := GetSwitch(k)
	newVal := !sw.Val
	log.Noticef("config change: %s: %t -> %t", k, sw.Val, newVal)
	sw.Val = newVal
	//log.Errorf("ignoring toggle for non-existant switch: %s", k)
}
