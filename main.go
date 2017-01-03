package main

import (
	ui "github.com/gizak/termui"
)

func main() {
	if err := ui.Init(); err != nil {
		panic(err)
	}
	defer ui.Close()

	g := NewGrid()
	for {
		exit := Display(g)
		if exit {
			return
		}
	}
}
