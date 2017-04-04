package expanded

import (
	"strconv"
	"strings"

	ui "github.com/gizak/termui"
)

type Ports struct {
	*ui.Table
	Exposed []int
	Open    [][]int
}

func NewPorts() *Ports {
	t := ui.NewTable()
	t.BorderLabel = "Ports"
	t.Height = 4
	t.Width = colWidth[0]
	t.FgColor = ui.ThemeAttr("par.text.fg")
	t.Separator = false
	p := &Ports{t, nil, nil}
	return p
}

func (p *Ports) Update(exposed []int64, open [][]int64) {
	p.Rows = [][]string{}

	exp_string := ""
	for i, exp := range exposed {
		if i == 0 {
			exp_string = strconv.Itoa(int(exp))
		} else {
			exp_string = strings.Join([]string{exp_string, strconv.Itoa(int(exp))}, ", ")
		}
	}
	p.Rows = append(p.Rows, []string{"Exposed: ", exp_string})

	open_string := ""
	for i, op := range open {
		ported := strings.Join([]string{strconv.Itoa(int(op[0])), strconv.Itoa(int(op[1]))}, " -> ")
		if i == 0 {
			open_string = ported
		} else {
			open_string = strings.Join([]string{open_string, ported}, ", ")
		}
	}
	p.Rows = append(p.Rows, []string{"Open: ", open_string})
}
