package manager

import models "github.com/bcicen/ctop/models"

type Mock struct{}

func NewMock() *Mock {
	return &Mock{}
}

func (m *Mock) Start() error {
	return ActionNotImplErr
}

func (m *Mock) Stop() error {
	return ActionNotImplErr
}

func (m *Mock) Remove() error {
	return ActionNotImplErr
}

func (m *Mock) Pause() error {
	return ActionNotImplErr
}

func (m *Mock) Unpause() error {
	return ActionNotImplErr
}

func (m *Mock) Restart() error {
	return ActionNotImplErr
}

func (m *Mock) Exec(cmd []string) error {
	return ActionNotImplErr
}

func (m *Mock) Inspect() (models.Meta, error) {
	return nil, nil
}
