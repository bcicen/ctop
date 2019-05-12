package manager

type Manager interface {
	Start() error
	Stop() error
	Remove() error
	Pause() error
	Unpause() error
	Restart() error
	Exec(cmd []string) error
}
