package mussed

import (
	"reflect"
	"text/template"
)

var RequiredFuncs = template.FuncMap{
	"mussedIsCollection": func(i interface{}) bool {
		it := reflect.TypeOf(i)
		switch it.Kind() {
		case reflect.Array, reflect.Slice:
			return true
		default:
			return false
		}
	},
}
