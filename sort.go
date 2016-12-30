package main

import (
	"sort"
)

type ContainerSorter interface {
	sort.Interface
	Sort()
}

var Sorters = map[string][]Container{
	"id":   ByID{},
	"name": ByName{},
	"cpu":  ByCPU{},
	"mem":  ByMem{},
}

type ByID []*Container

func (a ByID) Len() int           { return len(a) }
func (a ByID) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByID) Less(i, j int) bool { return a[i].id < a[j].id }
func (a ByID) Sort()              { sort.Sort(a) } // Sort is a convenience method.

type ByName []*Container

func (a ByName) Len() int           { return len(a) }
func (a ByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByName) Less(i, j int) bool { return a[i].id < a[j].id }
func (a ByName) Sort()              { sort.Sort(a) } // Sort is a convenience method.

type ByCPU []*Container

func (a ByCPU) Len() int           { return len(a) }
func (a ByCPU) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByCPU) Less(i, j int) bool { return a[i].reader.CPUUtil < a[j].reader.CPUUtil }
func (a ByCPU) Sort()              { sort.Sort(a) } // Sort is a convenience method.

type ByMem []*Container

func (a ByMem) Len() int           { return len(a) }
func (a ByMem) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByMem) Less(i, j int) bool { return a[i].reader.MemUsage < a[j].reader.MemUsage }
func (a ByMem) Sort()              { sort.Sort(a) } // Sort is a convenience method.
