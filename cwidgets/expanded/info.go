package expanded

import (
	ui "github.com/gizak/termui"
)

var displayInfo = []string{"id", "name", "image", "state"}

type Info struct {
	*ui.Table
	data map[string]string
}

func NewInfo(id string) *Info {
	p := ui.NewTable()
	p.Height = 4
	p.Width = colWidth[0]
	p.FgColor = ui.ColorWhite
	p.Seperator = false
	i := &Info{p, make(map[string]string)}
	i.Set("id", id)
	return i
}

func (w *Info) Set(k, v string) {
	w.data[k] = v
	// rebuild rows
	w.Rows = [][]string{}
	for _, k := range displayInfo {
		if v, ok := w.data[k]; ok {
			w.Rows = append(w.Rows, []string{k, v})
		}
	}
	w.Height = len(w.Rows) + 2
}
