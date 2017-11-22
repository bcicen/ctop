package manager

type Mock struct {}

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
