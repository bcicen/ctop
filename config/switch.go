package config

// defaults
var defaultSwitches = []*Switch{
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
}

type Switch struct {
	Key   string
	Val   bool
	Label string
}

// GetSwitch returns Switch by key
func GetSwitch(k string) *Switch {
	lock.RLock()
	defer lock.RUnlock()

	for _, sw := range GlobalSwitches {
		if sw.Key == k {
			return sw
		}
	}
	return &Switch{} // default
}

// GetSwitchVal returns Switch value by key
func GetSwitchVal(k string) bool {
	return GetSwitch(k).Val
}

func UpdateSwitch(k string, val bool) {
	sw := GetSwitch(k)

	lock.Lock()
	defer lock.Unlock()

	if sw.Val != val {
		log.Noticef("config change [%s]: %t -> %t", k, sw.Val, val)
		sw.Val = val
	}
}

// Toggle a boolean switch
func Toggle(k string) {
	sw := GetSwitch(k)

	lock.Lock()
	defer lock.Unlock()

	sw.Val = !sw.Val
	log.Noticef("config change [%s]: %t -> %t", k, !sw.Val, sw.Val)
	//log.Errorf("ignoring toggle for non-existant switch: %s", k)
}
