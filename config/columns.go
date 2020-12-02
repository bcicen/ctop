package config

import (
	"strings"
)

// defaults
var defaultColumns = []Column{
	Column{
		Name:    "status",
		Label:   "Status Indicator",
		Enabled: true,
	},
	Column{
		Name:    "name",
		Label:   "Container Name",
		Enabled: true,
	},
	Column{
		Name:    "id",
		Label:   "Container ID",
		Enabled: true,
	},
	Column{
		Name:    "cpu",
		Label:   "CPU Usage",
		Enabled: true,
	},
	Column{
		Name:    "cpus",
		Label:   "CPU Usage (% of system total)",
		Enabled: false,
	},
	Column{
		Name:    "mem",
		Label:   "Memory Usage",
		Enabled: true,
	},
	Column{
		Name:    "net",
		Label:   "Network RX/TX",
		Enabled: true,
	},
	Column{
		Name:    "io",
		Label:   "Disk IO Read/Write",
		Enabled: true,
	},
	Column{
		Name:    "pids",
		Label:   "Container PID Count",
		Enabled: true,
	},
}

type Column struct {
	Name    string
	Label   string
	Enabled bool
}

// ColumnsString returns an ordered and comma-delimited string of currently enabled Columns
func ColumnsString() string { return strings.Join(EnabledColumns(), ",") }

// EnabledColumns returns an ordered array of enabled column names
func EnabledColumns() (a []string) {
	lock.RLock()
	defer lock.RUnlock()
	for _, col := range GlobalColumns {
		if col.Enabled {
			a = append(a, col.Name)
		}
	}
	return a
}

// ColumnToggle toggles the enabled status of a given column name
func ColumnToggle(name string) {
	col := GlobalColumns[colIndex(name)]
	col.Enabled = !col.Enabled
	log.Noticef("config change [column-%s]: %t -> %t", col.Name, !col.Enabled, col.Enabled)
}

// ColumnLeft moves the column with given name up one position, if possible
func ColumnLeft(name string) {
	idx := colIndex(name)
	if idx > 0 {
		swapCols(idx, idx-1)
	}
}

// ColumnRight moves the column with given name up one position, if possible
func ColumnRight(name string) {
	idx := colIndex(name)
	if idx < len(GlobalColumns)-1 {
		swapCols(idx, idx+1)
	}
}

// Set Column order and enabled status from one or more provided Column names
func SetColumns(names []string) {
	var (
		n          int
		curColStr  = ColumnsString()
		newColumns = make([]*Column, len(GlobalColumns))
	)

	lock.Lock()

	// add enabled columns by name
	for _, name := range names {
		newColumns[n] = popColumn(name)
		newColumns[n].Enabled = true
		n++
	}

	// extend with omitted columns as disabled
	for _, col := range GlobalColumns {
		newColumns[n] = col
		newColumns[n].Enabled = false
		n++
	}

	GlobalColumns = newColumns
	lock.Unlock()

	log.Noticef("config change [columns]: %s -> %s", curColStr, ColumnsString())
}

func swapCols(i, j int) { GlobalColumns[i], GlobalColumns[j] = GlobalColumns[j], GlobalColumns[i] }

func popColumn(name string) *Column {
	idx := colIndex(name)
	if idx < 0 {
		panic("no such column name: " + name)
	}
	col := GlobalColumns[idx]
	GlobalColumns = append(GlobalColumns[:idx], GlobalColumns[idx+1:]...)
	return col
}

// return index of column with given name, if any
func colIndex(name string) int {
	for n, c := range GlobalColumns {
		if c.Name == name {
			return n
		}
	}
	return -1
}
