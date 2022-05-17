package tglogger

import (
	"fmt"
	"reflect"
)

type Fields map[string]interface{}

func (f Fields) String() (s string) {
	for k, v := range f {
		if t := reflect.TypeOf(v); t != nil {
			switch {
			case t.Kind() == reflect.Func, t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Func:
				v = "tgLogger - can not add field"
			}
		}

		s += fmt.Sprintf("%s: %v\n", k, v)
	}
	return
}
