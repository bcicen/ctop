package manager

import (
	"fmt"
	api "github.com/fsouza/go-dockerclient"
	"strings"
)

func PortsFormat(ports map[api.Port][]api.PortBinding) string {
	var exposed []string
	var published []string

	for k, v := range ports {
		if len(v) == 0 {
			// 3306/tcp
			exposed = append(exposed, string(k))
			continue
		}
		for _, binding := range v {
			// 0.0.0.0:3307 -> 3306/tcp
			s := fmt.Sprintf("%s:%s -> %s", binding.HostIP, binding.HostPort, k)
			published = append(published, s)
		}
	}

	return strings.Join(append(exposed, published...), "\n")
}

func PortsFormatArr(ports []api.APIPort) string {
	var exposed []string
	var published []string
	for _, binding := range ports {
		if binding.PublicPort != 0 {
			// 0.0.0.0:3307 -> 3306/tcp
			s := fmt.Sprintf("%s:%d -> %d/%s", binding.IP, binding.PublicPort, binding.PrivatePort, binding.Type)
			published = append(published, s)
		} else {
			// 3306/tcp
			s := fmt.Sprintf("%d/%s", binding.PrivatePort, binding.Type)
			exposed = append(exposed, s)
		}
	}

	return strings.Join(append(exposed, published...), "\n")
}

func IpsFormat(networks map[string]api.ContainerNetwork) string {
	var ips []string

	for k, v := range networks {
		s := fmt.Sprintf("%s:%s", k, v.IPAddress)
		ips = append(ips, s)
	}

	return strings.Join(ips, "\n")
}

// use primary container name
func ShortName(name string) string {
	return strings.TrimPrefix(name, "/")
}
