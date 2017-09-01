package entity

import (
	"fmt"
	"regexp"
	"sort"

	"github.com/bcicen/ctop/config"
	"strings"
)

type sortMethod func(c1, c2 Entity) bool

var stateMap = map[string]int{
	"new":                6, //>
	"running":            5, //>
	"rollback_completed": 5,
	"ready":              4, //>
	"rollback_paused":    4,
	"updating":           4,
	"starting":           4,
	"paused":             3, //>
	"exited":             2, //>
	"shutdown":           2,
	"failed":             1, //>
	"created":            0, //>
	"":                   0,
}

var idSorter = func(c1, c2 Entity) bool { return c1.GetId() < c2.GetId() }

var nameSorter = func(c1, c2 Entity) bool {
	name1 := strings.TrimSpace(c1.GetMeta("name"))
	name1 = strings.Trim(name1, "\\_")

	name2 := strings.TrimSpace(c2.GetMeta("name"))
	name2 = strings.Trim(name2, "\\_")
	if name1 == name2 {
		c1state := c1.GetMeta("state")
		c2state := c2.GetMeta("state")
		return stateMap[c1state] > stateMap[c2state]
	}
	if name1 == "" || name2 == "" {
		return idSorter(c1, c2)
	}
	return c1.GetMeta("name") < c2.GetMeta("name")
}

var Sorters = map[string]sortMethod{
	"id":   idSorter,
	"name": nameSorter,
	"cpu": func(c1, c2 Entity) bool {
		// Use secondary sort method if equal values
		if c1.GetMetrics().CPUUtil == c2.GetMetrics().CPUUtil {
			return nameSorter(c1, c2)
		}
		return c1.GetMetrics().CPUUtil > c2.GetMetrics().CPUUtil
	},
	"mem": func(c1, c2 Entity) bool {
		// Use secondary sort method if equal values
		if c1.GetMetrics().MemUsage == c2.GetMetrics().MemUsage {
			return nameSorter(c1, c2)
		}
		return c1.GetMetrics().MemUsage > c2.GetMetrics().MemUsage
	},
	"mem %": func(c1, c2 Entity) bool {
		// Use secondary sort method if equal values
		if c1.GetMetrics().MemPercent == c2.GetMetrics().MemPercent {
			return nameSorter(c1, c2)
		}
		return c1.GetMetrics().MemPercent > c2.GetMetrics().MemPercent
	},
	"net": func(c1, c2 Entity) bool {
		sum1 := sumNet(c1)
		sum2 := sumNet(c2)
		// Use secondary sort method if equal values
		if sum1 == sum2 {
			return nameSorter(c1, c2)
		}
		return sum1 > sum2
	},
	"pids": func(c1, c2 Entity) bool {
		// Use secondary sort method if equal values
		if c1.GetMetrics().Pids == c2.GetMetrics().Pids {
			return nameSorter(c1, c2)
		}
		return c1.GetMetrics().Pids > c2.GetMetrics().Pids
	},
	"io": func(c1, c2 Entity) bool {
		sum1 := sumIO(c1)
		sum2 := sumIO(c2)
		// Use secondary sort method if equal values
		if sum1 == sum2 {
			return nameSorter(c1, c2)
		}
		return sum1 > sum2
	},
	"state": func(c1, c2 Entity) bool {
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

func (c Containers) Sort()         { sort.Sort(c) }
func (c Containers) Len() int      { return len(c) }
func (c Containers) Swap(i, j int) { c[i], c[j] = c[j], c[i] }
func (c Containers) Less(i, j int) bool {
	f := Sorters[config.GetVal("sortField")]
	if config.GetSwitchVal("sortReversed") {
		return f(c[j], c[i])
	}
	return f(c[i], c[j])
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

type Nodes []*Node

func (n Nodes) Filter() {
	filter := config.GetVal("filterStr")
	re := regexp.MustCompile(fmt.Sprintf(".*%s", filter))

	for _, c := range n {
		c.Display = true
		// Apply name filter
		if re.FindAllString(c.GetMeta("name"), 1) == nil {
			c.Display = false
		}
	}
}

type Services []*Service

func (s Services) Sort()         { sort.Sort(s) }
func (s Services) Len() int      { return len(s) }
func (s Services) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s Services) Less(i, j int) bool {
	f := Sorters[config.GetVal("sortField")]
	if config.GetSwitchVal("sortReversed") {
		return f(s[j], s[i])
	}
	return f(s[i], s[j])
}

func (s Services) Filter() {
	filter := config.GetVal("filterStr")
	re := regexp.MustCompile(fmt.Sprintf(".*%s", filter))

	for _, c := range s {
		c.Display = true
		// Apply name filter
		if re.FindAllString(c.GetMeta("name"), 1) == nil {
			c.Display = false
		}
	}
}

type Tasks []*Task

func (t Tasks) Filter() {
	filter := config.GetVal("filterStr")
	re := regexp.MustCompile(fmt.Sprintf(".*%s", filter))

	for _, c := range t {
		c.Display = true
		// Apply name filter
		if re.FindAllString(c.GetMeta("name"), 1) == nil {
			c.Display = false
		}
	}
}

func (t Tasks) Sort()         { sort.Sort(t) }
func (t Tasks) Len() int      { return len(t) }
func (t Tasks) Swap(i, j int) { t[i], t[j] = t[j], t[i] }
func (t Tasks) Less(i, j int) bool {
	f := Sorters[config.GetVal("sortField")]
	if config.GetSwitchVal("sortReversed") {
		return f(t[j], t[i])
	}
	return f(t[i], t[j])
}

type Entities []Entity

func sumNet(c Entity) int64 { return c.GetMetrics().NetRx + c.GetMetrics().NetTx }

func sumIO(c Entity) int64 { return c.GetMetrics().IOBytesRead + c.GetMetrics().IOBytesWrite }
