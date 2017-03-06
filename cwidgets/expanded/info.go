package expanded

import (
	ui "github.com/gizak/termui"
)

type Info struct {
	*ui.Table
	data map[string]string
}

func NewInfo(id string) *Info {
	p := ui.NewTable()
	p.Height = 4
	p.Width = 50
	p.FgColor = ui.ColorWhite
	p.Seperator = false
	i := &Info{p, make(map[string]string)}
	i.Set("ID", id)
	return i
}

func (w *Info) Set(k, v string) {
	w.data[k] = v
	// rebuild rows
	w.Rows = [][]string{}
	for k, v := range w.data {
		w.Rows = append(w.Rows, []string{k, v})
	}
}
