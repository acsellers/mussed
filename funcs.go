package mussed

import (
	"fmt"
	"html/template"
	"reflect"
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
	"mussedUnescape": func(i ...interface{}) template.HTML {
		if len(i) == 1 {
			return template.HTML(fmt.Sprint(i[0]))
		}
		return template.HTML("")
	},
	"mussedUpscope":   upscope,
	"mussedDownscope": downscope,
}
