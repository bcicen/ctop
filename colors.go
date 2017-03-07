package main

import ui "github.com/gizak/termui"

/*
Valid colors:
	ui.ColorDefault
	ui.ColorBlack
	ui.ColorRed
	ui.ColorGreen
	ui.ColorYellow
	ui.ColorBlue
	ui.ColorMagenta
	ui.ColorCyan
	ui.ColorWhite
*/

var ColorMap = map[string]ui.Attribute{
	"fg":                 ui.ColorWhite,
	"bg":                 ui.ColorDefault,
	"block.bg":           ui.ColorDefault,
	"border.bg":          ui.ColorDefault,
	"border.fg":          ui.ColorWhite,
	"label.bg":           ui.ColorDefault,
	"label.fg":           ui.ColorGreen,
	"menu.text.fg":       ui.ColorWhite,
	"menu.text.bg":       ui.ColorDefault,
	"menu.border.fg":     ui.ColorCyan,
	"menu.label.fg":      ui.ColorGreen,
	"header.fg":          ui.ColorBlack,
	"header.bg":          ui.ColorWhite,
	"gauge.bar.bg":       ui.ColorGreen,
	"gauge.percent.fg":   ui.ColorWhite,
	"linechart.axes.fg":  ui.ColorDefault,
	"linechart.line.fg":  ui.ColorGreen,
	"mbarchart.bar.bg":   ui.ColorGreen,
	"mbarchart.num.fg":   ui.ColorWhite,
	"mbarchart.text.fg":  ui.ColorWhite,
	"par.text.fg":        ui.ColorWhite,
	"par.text.bg":        ui.ColorDefault,
	"sparkline.line.fg":  ui.ColorGreen,
	"sparkline.title.fg": ui.ColorWhite,
}
