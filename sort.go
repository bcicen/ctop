package main

import (
	"fmt"
	"regexp"

	"github.com/bcicen/ctop/config"
)

type sortMethod func(c1, c2 *Container) bool

var stateMap = map[string]int{
	"running": 3,
	"paused":  2,
	"exited":  1,
	"created": 0,
	"":        0,
}

var idSorter = func(c1, c2 *Container) bool { return c1.Id < c2.Id }
var nameSorter = func(c1, c2 *Container) bool { return c1.GetMeta("name") < c2.GetMeta("name") }

var Sorters = map[string]sortMethod{
	"id":   idSorter,
	"name": nameSorter,
	"cpu": func(c1, c2 *Container) bool {
		// Use secondary sort method if equal values
		if c1.CPUUtil == c2.CPUUtil {
			return nameSorter(c1, c2)
		}
		return c1.CPUUtil > c2.CPUUtil
	},
	"mem": func(c1, c2 *Container) bool {
		// Use secondary sort method if equal values
		if c1.MemUsage == c2.MemUsage {
			return nameSorter(c1, c2)
		}
		return c1.MemUsage > c2.MemUsage
	},
	"mem %": func(c1, c2 *Container) bool {
		// Use secondary sort method if equal values
		if c1.MemPercent == c2.MemPercent {
			return nameSorter(c1, c2)
		}
		return c1.MemPercent > c2.MemPercent
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
	"io": func(c1, c2 *Container) bool {
		sum1 := sumIO(c1)
		sum2 := sumIO(c2)
		// Use secondary sort method if equal values
		if sum1 == sum2 {
			return nameSorter(c1, c2)
		}
		return sum1 > sum2
	},
	"state": func(c1, c2 *Container) bool {
		// Use secondary sort method if equal values
		c1state := c1.GetMeta("state")
		c2state := c2.GetMeta("state")
		if c1state == c2state {
			return nameSorter(c1, c2)
		}
		return stateMap[c1state] > stateMap[c2state]
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

func (a Containers) Filter() {
	filter := config.GetVal("filterStr")
	re := regexp.MustCompile(fmt.Sprintf(".*%s", filter))

	for _, c := range a {
		c.display = true
		// Apply name filter
		if re.FindAllString(c.GetMeta("name"), 1) == nil {
			c.display = false
		}
		// Apply state filter
		if !config.GetSwitchVal("allContainers") && c.GetMeta("state") != "running" {
			c.display = false
		}
	}
}

func sumNet(c *Container) int64 { return c.NetRx + c.NetTx }

func sumIO(c *Container) int64 { return c.IOBytesRead + c.IOBytesWrite }
