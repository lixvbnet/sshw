package sshlib

import (
	"fmt"
	"reflect"
)

const debug = false

// CoverDefaults assigns default values to zero-valued fields in node if override is set to false.
// It overrides node values for all non-zero fields in defaults if override is set to true.
// The two arguments should be of the same type.
func CoverDefaults(node interface{}, defaults interface{}, override bool) {
	val := reflect.ValueOf(node).Elem()
	valDefaults := reflect.ValueOf(defaults).Elem()

	for i := 0; i < val.NumField(); i++ {
		name := val.Type().Field(i).Name
		valueField := val.Field(i)
		if debug {
			value := valueField.Interface()
			fmt.Print(name, "=", value, "\t")
		}

		defaultValueField := valDefaults.FieldByName(name)
		if !defaultValueField.IsZero() {
			if override || valueField.IsZero() {
				if debug {
					fmt.Print("Set to default value ", defaultValueField.Interface())
				}
				valueField.Set(defaultValueField)
			}
		}

		if debug {
			fmt.Println()
		}
	}
}
