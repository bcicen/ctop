package main

import (
	"fmt"
	"regexp"

	"github.com/bcicen/ctop/config"
)

type sortMethod func(c1, c2 *Container) bool

var idSorter = func(c1, c2 *Container) bool { return c1.id < c2.id }
var nameSorter = func(c1, c2 *Container) bool { return c1.name < c2.name }

var Sorters = map[string]sortMethod{
	"id":   idSorter,
	"name": nameSorter,
	"cpu": func(c1, c2 *Container) bool {
		// Use secondary sort method if equal values
		if c1.metrics.CPUUtil == c2.metrics.CPUUtil {
			return nameSorter(c1, c2)
		}
		return c1.metrics.CPUUtil > c2.metrics.CPUUtil
	},
	"mem": func(c1, c2 *Container) bool {
		// Use secondary sort method if equal values
		if c1.metrics.MemUsage == c2.metrics.MemUsage {
			return nameSorter(c1, c2)
		}
		return c1.metrics.MemUsage > c2.metrics.MemUsage
	},
	"mem %": func(c1, c2 *Container) bool {
		// Use secondary sort method if equal values
		if c1.metrics.MemPercent == c2.metrics.MemPercent {
			return nameSorter(c1, c2)
		}
		return c1.metrics.MemPercent > c2.metrics.MemPercent
	},
	"net": func(c1, c2 *Container) bool {
		sum1 := sumNet(c1)
		sum2 := sumNet(c2)
		// Use secondary sort method if equal values
		if sum1 == sum2 {
			return nameSorter(c1, c2)
		}
		return sum1 > sum2
	},
	"state": func(c1, c2 *Container) bool {
		// Use secondary sort method if equal values
		if c1.state == c2.state {
			return nameSorter(c1, c2)
		}
		if c1.state == "running" {
			return true
		}
		if c2.state == "running" {
			return false
		}
		if c2.state == "paused" {
			return false
		}
		return true
	},
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
