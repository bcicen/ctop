package config

// defaults
var params = []*Param{
	&Param{
		key:   "dockerHost",
		val:   getEnv("DOCKER_HOST", "unix:///var/run/docker.sock"),
		label: "Docker API URL",
	},
	&Param{
		key:   "filterStr",
		val:   "",
		label: "Container Name or ID Filter",
	},
	&Param{
		key:   "sortField",
		val:   "id",
		label: "Container Sort Field",
	},
}

type Param struct {
	key   string
	val   string
	label string
}

// Return param value
func Get(k string) string {
	for _, p := range GlobalParams {
		if p.key == k {
			return p.val
		}
	}
	return "" // default
}

// Set param value
func Update(k, v string) {
	for _, p := range GlobalParams {
		if p.key == k {
			log.Noticef("config change: %s: %s -> %s", k, p.val, v)
			p.val = v
			return
		}
	}
	log.Errorf("ignoring update for non-existant parameter: %s", k)
}
