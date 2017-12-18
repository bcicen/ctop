package manager

type Manager interface {
	Start() error
	Stop() error
	Remove() error
}
