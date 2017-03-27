package main

import (
	ui "github.com/gizak/termui"
)

// Common action keybindings
// Common action keybindings
var keyMap = map[string][]string{
	"up": []string{
		"/sys/kbd/<up>",
		"/sys/kbd/k",
	},
	"down": []string{
		"/sys/kbd/<down>",
		"/sys/kbd/j",
	},
	"pgup": []string{
		"/sys/kbd/<previous>",
		"/sys/kbd/C-<up>",
	},
	"pgdown": []string{
		"/sys/kbd/<next>",
		"/sys/kbd/C-<down>",
    },
	"exit": []string{
		"/sys/kbd/q",
		"/sys/kbd/C-c",
	},
	"help": []string{
		"/sys/kbd/h",
		"/sys/kbd/?",
	},
}

// Apply a common handler function to all given keys
func HandleKeys(i string, f func()) {
	for _, k := range keyMap[i] {
		ui.Handle(k, func(ui.Event) { f() })
	}
}
