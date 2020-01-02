package config

import (
	"strings"
)

// defaults
var defaultWidgets = []*Widget{
	&Widget{
		Name:    "status",
		Enabled: true,
	},
	&Widget{
		Name:    "name",
		Enabled: true,
	},
	&Widget{
		Name:    "id",
		Enabled: true,
	},
	&Widget{
		Name:    "cpu",
		Enabled: true,
	},
	&Widget{
		Name:    "mem",
		Enabled: true,
	},
	&Widget{
		Name:    "net",
		Enabled: true,
	},
	&Widget{
		Name:    "io",
		Enabled: true,
	},
	&Widget{
		Name:    "pids",
		Enabled: true,
	},
}

type Widget struct {
	Name    string
	Enabled bool
}

// GetWidget returns a Widget by name
func GetWidget(name string) *Widget {
	lock.RLock()
	defer lock.RUnlock()

	for _, w := range GlobalWidgets {
		if w.Name == name {
			return w
		}
	}
	log.Errorf("widget name not found: %s", name)
	return &Widget{} // default
}

// Widgets returns a copy of all configurable Widgets, in order
func Widgets() []Widget {
	a := make([]Widget, len(GlobalWidgets))

	lock.RLock()
	defer lock.RUnlock()

	for n, w := range GlobalWidgets {
		a[n] = *w
	}
	return a
}

// EnabledWidgets returns an ordered array of enabled widget names
func EnabledWidgets() (a []string) {
	for _, w := range Widgets() {
		if w.Enabled {
			a = append(a, w.Name)
		}
	}
	return a
}

func UpdateWidget(name string, enabled bool) {
	w := GetWidget(name)
	oldVal := w.Enabled
	log.Noticef("config change [%s-enabled]: %t -> %t", name, oldVal, enabled)

	lock.Lock()
	defer lock.Unlock()
	w.Enabled = enabled
}

func ToggleWidgetEnabled(name string) {
	w := GetWidget(name)
	newVal := !w.Enabled
	log.Noticef("config change [%s-enabled]: %t -> %t", name, w.Enabled, newVal)

	lock.Lock()
	defer lock.Unlock()
	w.Enabled = newVal
}

// UpdateWidgets replaces existing ordered widgets with those provided
func UpdateWidgets(newWidgets []Widget) {
	oldOrder := widgetNames()
	lock.Lock()
	for n, w := range newWidgets {
		GlobalWidgets[n] = &w
	}
	lock.Unlock()
	log.Noticef("config change [widget-order]: %s -> %s", oldOrder, widgetNames())
}

func widgetNames() string {
	a := make([]string, len(GlobalWidgets))
	for n, w := range Widgets() {
		a[n] = w.Name
	}
	return strings.Join(a, ", ")
}
