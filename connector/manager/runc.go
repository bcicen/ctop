package manager

type Runc struct{}

func NewRunc() *Runc {
	return &Runc{}
}

func (rc *Runc) Start() error {
	return nil
}

func (rc *Runc) Stop() error {
	return nil
}

func (rc *Runc) Remove() error {
	return nil
}
