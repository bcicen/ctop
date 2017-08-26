package entity

import (
	"github.com/bcicen/ctop/cwidgets/compact"
	"github.com/bcicen/ctop/cwidgets"
)

type Meta struct {
	Meta    map[string]string
	Widgets *compact.Compact
	Display bool // display this service in compact view
	updater cwidgets.WidgetUpdater
}

func NewMeta(id string, ) Meta {
	widgets := compact.NewCompact(id)
	return Meta{
		Meta:    make(map[string]string),
		Widgets: widgets,
		updater: widgets,
	}
}

func (m *Meta) SetUpdater(u cwidgets.WidgetUpdater) {
	m.updater = u
	for k, v := range m.Meta {
		m.updater.SetMeta(k, v)
	}
}

func (m *Meta) SetMeta(k, v string) {
	m.Meta[k] = v
	m.updater.SetMeta(k, v)
}

func (m *Meta) GetMeta(k string) string {
	if v, ok := m.Meta[k]; ok {
		return v
	}
	return ""
}
