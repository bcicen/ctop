package config

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
	if _, ok := Global.params[k]; ok == true {
		return Global.params[k].val
	}
	return ""
}
