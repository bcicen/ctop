package single

import (
	"time"

	"github.com/bcicen/ctop/models"
	ui "github.com/gizak/termui"
)

type LogLines struct {
	ts   []time.Time
	data []string
}

func NewLogLines(max int) *LogLines {
	ll := &LogLines{
		ts:   make([]time.Time, max),
		data: make([]string, max),
	}
	return ll
}

func (ll *LogLines) tail(n int) []string {
	lines := make([]string, n)
	for i := 0; i < n; i++ {
		lines = append(lines, ll.data[len(ll.data)-i])
	}
	return lines
}
func (ll *LogLines) getLines(start, end int) []string {
	if end < 0 {
		return ll.data[start:]
	}
	return ll.data[start:end]
}

func (ll *LogLines) add(l models.Log) {
	if len(ll.data) == cap(ll.data) {
		ll.data = append(ll.data[:0], ll.data[1:]...)
		ll.ts = append(ll.ts[:0], ll.ts[1:]...)
	}
	ll.ts = append(ll.ts, l.Timestamp)
	ll.data = append(ll.data, l.Message)
	log.Debugf("recorded log line: %v", l)
}

type Logs struct {
	*ui.List
	lines *LogLines
}

func NewLogs(stream chan models.Log) *Logs {
	p := ui.NewList()
	p.Y = ui.TermHeight() / 2
	p.X = 0
	p.Height = ui.TermHeight() - p.Y
	p.Width = ui.TermWidth()
	//p.Overflow = "wrap"
	p.ItemFgColor = ui.ThemeAttr("par.text.fg")
	i := &Logs{p, NewLogLines(4098)}
	go func() {
		for line := range stream {
			i.lines.add(line)
			ui.Render(i)
		}
	}()
	return i
}

func (w *Logs) Align() {
	w.X = colWidth[0]
	w.List.Align()
}

func (w *Logs) Buffer() ui.Buffer {
	maxLines := w.Height - 2
	offset := len(w.lines.data) - maxLines
	w.Items = w.lines.getLines(offset, -1)
	return w.List.Buffer()
}

// number of rows a line will occupy at current panel width
func (w *Logs) lineHeight(s string) int { return (len(s) / w.InnerWidth()) + 1 }
