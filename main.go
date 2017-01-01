package main

import (
	"os"
)

func main() {
	g := NewGrid()
	for {
		exit := Display(g)
		if exit {
			os.Exit(0)
		}
	}
}
