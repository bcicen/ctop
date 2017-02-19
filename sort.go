package main

import (
	"fmt"
	"regexp"

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

func SortFields() (fields []string) {
	for k := range Sorters {
		fields = append(fields, k)
	}
	return fields
}

type Containers []*Container

func (a Containers) Len() int      { return len(a) }
func (a Containers) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a Containers) Less(i, j int) bool {
	f := Sorters[config.GetVal("sortField")]
	if config.GetSwitchVal("sortReversed") {
		return f(a[j], a[i])
	}
	return f(a[i], a[j])
}

func (a Containers) Filter() (filtered []*Container) {
	filter := config.GetVal("filterStr")
	re := regexp.MustCompile(fmt.Sprintf(".*%s", filter))

	for _, c := range a {
		// Apply name filter
		if re.FindAllString(c.name, 1) == nil {
			continue
		}
		// Apply state filter
		if !config.GetSwitchVal("allContainers") && c.state != "running" {
			continue
		}
		filtered = append(filtered, c)
	}

	return filtered
}

func sumNet(c *Container) int64 { return c.metrics.NetRx + c.metrics.NetTx }
