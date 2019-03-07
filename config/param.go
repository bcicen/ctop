package config

// defaults
var params = []*Param{
	&Param{
		Key:   "filterStr",
		Val:   "",
		Label: "Container Name or ID Filter",
	},
	&Param{
		Key:   "sortField",
		Val:   "state",
		Label: "Container Sort Field",
	},
}

type Param struct {
	Key   string
	Val   string
	Label string
}

// Get Param by key
func Get(k string) *Param {
	for _, p := range GlobalParams {
		if p.Key == k {
			return p
		}
	}
	return &Param{} // default
}

// GetVal gets Param value by key
func GetVal(k string) string {
	return Get(k).Val
}

// Set param value
func Update(k, v string) {
	p := Get(k)
	log.Noticef("config change: %s: %s -> %s", k, quote(p.Val), quote(v))
	p.Val = v
	// log.Errorf("ignoring update for non-existant parameter: %s", k)
}
