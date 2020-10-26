package manager

import "errors"

var ActionNotImplErr = errors.New("action not implemented")

type Manager interface {
	Start() error
	Stop() error
	Remove() error
	Pause() error
	Unpause() error
	Restart() error
	Exec(cmd []string) error
}
