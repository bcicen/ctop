package main

import (
	"fmt"
	"reflect"

	ui "github.com/gizak/termui"
)

func logEvent(e ui.Event) {
	var s string
	s += fmt.Sprintf("Type=%s", quote(e.Type))
	s += fmt.Sprintf(" Path=%s", quote(e.Path))
	s += fmt.Sprintf(" From=%s", quote(e.From))
	if e.To != "" {
		s += fmt.Sprintf(" To=%s", quote(e.To))
	}
	log.Debugf("new event: %s", s)
}

// log container, metrics, and widget state
func dumpContainer(c *Container) {
	msg := fmt.Sprintf("logging state for container: %s\n", c.Id)
	for k, v := range c.Meta {
		msg += fmt.Sprintf("Meta.%s = %s\n", k, v)
	}
	msg += inspect(&c.Metrics)
	log.Infof(msg)
}

func inspect(i interface{}) (s string) {
	val := reflect.ValueOf(i)
	elem := val.Type().Elem()

	eName := elem.String()
	for i := 0; i < elem.NumField(); i++ {
		field := elem.Field(i)
		fieldVal := reflect.Indirect(val).FieldByName(field.Name)
		s += fmt.Sprintf("%s.%s = ", eName, field.Name)
		s += fmt.Sprintf("%v (%s)\n", fieldVal, field.Type)
	}
	return s
}

func quote(s string) string {
	return fmt.Sprintf("\"%s\"", s)
}
