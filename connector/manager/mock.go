package manager

type Mock struct{}

func NewMock() *Mock {
	return &Mock{}
}

func (m *Mock) Start() error {
	return nil
}

func (m *Mock) Stop() error {
	return nil
}

func (m *Mock) Remove() error {
	return nil
}

func (m *Mock) Pause() error {
	return nil
}

func (m *Mock) Unpause() error {
	return nil
}

func (m *Mock) Restart() error {
	return nil
}

func (m *Mock) Exec(cmd []string) error {
	return nil
}
