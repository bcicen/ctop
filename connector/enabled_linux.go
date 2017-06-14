// +build !darwin

package connector

var enabled = map[string]func() Connector{
	"docker": NewDocker,
	"runc":   NewRunc,
}
