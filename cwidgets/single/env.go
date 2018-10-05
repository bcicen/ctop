package single

import (
	ui "github.com/gizak/termui"
	"regexp"
)

var envPattern = regexp.MustCompile(`(?P<KEY>[^=]+)=(?P<VALUJE>.*)`)

type Env struct {
	*ui.Table
	data map[string]string
}

func NewEnv() *Env {
	p := ui.NewTable()
	p.Height = 4
	p.Width = colWidth[0]
	p.FgColor = ui.ThemeAttr("par.text.fg")
	p.Separator = false
	i := &Env{p, make(map[string]string)}
	i.BorderLabel = "Env"
	return i
}

func (w *Env) Set(k, v string) {
	match := envPattern.FindStringSubmatch(v)
	key := match[1]
	value := match[2]
	w.data[key] = value

	w.Rows = [][]string{}
	w.Rows = append(w.Rows, mkInfoRows(key, value)...)

	w.Height = len(w.Rows) + 2
}
