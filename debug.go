package main

import (
	"fmt"
	"reflect"
)

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
