package manager

type Runc struct{}

func NewRunc() *Runc {
	return &Runc{}
}

func (rc *Runc) Start() error {
	return ActionNotImplErr
}

func (rc *Runc) Stop() error {
	return ActionNotImplErr
}

func (rc *Runc) Remove() error {
	return ActionNotImplErr
}

func (rc *Runc) Pause() error {
	return ActionNotImplErr
}

func (rc *Runc) Unpause() error {
	return ActionNotImplErr
}

func (rc *Runc) Restart() error {
	return ActionNotImplErr
}

func (rc *Runc) Exec(cmd []string) error {
	return ActionNotImplErr
}
