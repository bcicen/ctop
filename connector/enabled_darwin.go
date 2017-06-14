// +build !linux

package connector

var enabled = map[string]func() Connector{
	"docker": NewDocker,
}
