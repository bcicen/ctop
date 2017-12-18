package main

import (
	"fmt"
	"reflect"
	"runtime"

	ui "github.com/gizak/termui"
	"github.com/bcicen/ctop/entity"
)

var mstats = &runtime.MemStats{}

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

func runtimeStats() {
	var msg string
	msg += fmt.Sprintf("cgo calls=%v", runtime.NumCgoCall())
	msg += fmt.Sprintf(" routines=%v", runtime.NumGoroutine())
	runtime.ReadMemStats(mstats)
	msg += fmt.Sprintf(" numgc=%v", mstats.NumGC)
	msg += fmt.Sprintf(" alloc=%v", mstats.Alloc)
	log.Debugf("runtime: %v", msg)
}

func runtimeStack() {
	buf := make([]byte, 32768)
	buf = buf[:runtime.Stack(buf, true)]
	log.Infof(fmt.Sprintf("stack:\n%v", string(buf)))
}

// log container, metrics, and widget state
func dumpContainer(c entity.Entity) {
	msg := fmt.Sprintf("logging state for container: %s\n", c.GetId())
	for k, v := range c.GetMetaEntity().Meta {
		msg += fmt.Sprintf("Meta.%s = %s\n", k, v)
	}
	msg += inspect(c.GetMetrics())
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
