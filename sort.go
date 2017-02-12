package main

import (
	"sort"

	"github.com/bcicen/ctop/config"
)

type sortMethod func(c1, c2 *Container) bool

var Sorters = map[string]sortMethod{
	"id":    func(c1, c2 *Container) bool { return c1.id < c2.id },
	"name":  func(c1, c2 *Container) bool { return c1.name < c2.name },
	"cpu":   func(c1, c2 *Container) bool { return c1.metrics.CPUUtil < c2.metrics.CPUUtil },
	"mem":   func(c1, c2 *Container) bool { return c1.metrics.MemUsage < c2.metrics.MemUsage },
	"mem %": func(c1, c2 *Container) bool { return c1.metrics.MemPercent < c2.metrics.MemPercent },
	"net":   func(c1, c2 *Container) bool { return sumNet(c1) < sumNet(c2) },
}

func SortFields() []string {
	a := sort.StringSlice{}
	for k := range Sorters {
		a = append(a, k)
	}
	sort.Sort(a)
	return a
}

type Containers []*Container

func (a Containers) Len() int      { return len(a) }
func (a Containers) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a Containers) Less(i, j int) bool {
	f := Sorters[config.Get("sortField")]
	if config.GetToggle("sortReversed") {
		return f(a[j], a[i])
	}
	return f(a[i], a[j])
}

func sumNet(c *Container) int64 { return c.metrics.NetRx + c.metrics.NetTx }
