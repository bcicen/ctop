package container

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

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

func cmpByProject(c1, c2 *Container) int {
	return strings.Compare(c1.Project.Name, c2.Project.Name)
}

func cmpByName(c1, c2 *Container) bool {
	return c1.GetMeta("name") < c2.GetMeta("name")
}

var idSorter = func(c1, c2 *Container) bool { return c1.Id < c2.Id }

var Sorters = map[string]sortMethod{
	"id": idSorter,
	"name": func(c1, c2 *Container) bool {
		if projCmp := cmpByProject(c1, c2); projCmp != 0 {
			return projCmp < 0
		}
		return cmpByName(c1, c2)
	},
	"cpu": func(c1, c2 *Container) bool {
		// Use secondary sort method if equal values
		if c1.CPUUtil == c2.CPUUtil {
			return cmpByName(c1, c2)
		}
		return c1.CPUUtil > c2.CPUUtil
	},
	"mem": func(c1, c2 *Container) bool {
		// Use secondary sort method if equal values
		if c1.MemUsage == c2.MemUsage {
			return cmpByName(c1, c2)
		}
		return c1.MemUsage > c2.MemUsage
	},
	"mem %": func(c1, c2 *Container) bool {
		// Use secondary sort method if equal values
		if c1.MemPercent == c2.MemPercent {
			return cmpByName(c1, c2)
		}
		return c1.MemPercent > c2.MemPercent
	},
	"net": func(c1, c2 *Container) bool {
		sum1 := c1.SumNet()
		sum2 := c2.SumNet()
		// Use secondary sort method if equal values
		if sum1 == sum2 {
			return cmpByName(c1, c2)
		}
		return sum1 > sum2
	},
	"pids": func(c1, c2 *Container) bool {
		// Use secondary sort method if equal values
		if c1.Pids == c2.Pids {
			return cmpByName(c1, c2)
		}
		return c1.Pids > c2.Pids
	},
	"io": func(c1, c2 *Container) bool {
		sum1 := c1.SumIO()
		sum2 := c2.SumIO()
		// Use secondary sort method if equal values
		if sum1 == sum2 {
			return cmpByName(c1, c2)
		}
		return sum1 > sum2
	},
	"state": func(c1, c2 *Container) bool {
		// Use secondary sort method if equal values
		c1state := c1.GetMeta("state")
		c2state := c2.GetMeta("state")
		if c1state == c2state {
			return cmpByName(c1, c2)
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

func (a Containers) Sort()         { sort.Sort(a) }
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
		c.Display = true
		// Apply name filter
		if re.FindAllString(c.GetMeta("name"), 1) == nil {
			c.Display = false
		}
		// Apply state filter
		if !config.GetSwitchVal("allContainers") && c.GetMeta("state") != "running" {
			c.Display = false
		}
	}
}
