package main

import (
	"sort"
)

var SortFields = []string{"id", "name", "cpu", "mem"}

type ByID []*Container

func (a ByID) Len() int           { return len(a) }
func (a ByID) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByID) Less(i, j int) bool { return a[i].id < a[j].id }

type ByName []*Container

func (a ByName) Len() int           { return len(a) }
func (a ByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByName) Less(i, j int) bool { return a[i].name < a[j].name }

type ByCPU []*Container

func (a ByCPU) Len() int           { return len(a) }
func (a ByCPU) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByCPU) Less(i, j int) bool { return a[i].reader.CPUUtil < a[j].reader.CPUUtil }

type ByMem []*Container

func (a ByMem) Len() int           { return len(a) }
func (a ByMem) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByMem) Less(i, j int) bool { return a[i].reader.MemUsage < a[j].reader.MemUsage }

// Return array of containers, sorted by field
func (cm *ContainerMap) Sorted() []*Container {
	containers := cm.All()

	switch cm.sortField {
	case "id":
		sort.Sort(ByID(containers))
	case "name":
		sort.Sort(ByName(containers))
	case "cpu":
		sort.Sort(sort.Reverse(ByCPU(containers)))
	case "mem":
		sort.Sort(sort.Reverse(ByMem(containers)))
	default:
		sort.Sort(ByID(containers))
	}

	return containers
}
