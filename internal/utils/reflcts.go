package utils

import (
	"reflect"
	"strconv"
)

func ParseGoTagToStruct(goTag string, getter func(string) (string, bool), cfg interface{}) {
	val := reflect.ValueOf(cfg).Elem()
	typ := val.Type()
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		tag := field.Tag.Get(goTag)
		if tag == "" {
			continue
		}
		labelVal, ok := getter(tag)
		if !ok {
			continue
		}
		fv := val.Field(i)
		if !fv.CanSet() {
			continue
		}
		switch fv.Kind() {
		case reflect.String:
			fv.SetString(labelVal)
		case reflect.Bool:
			b, _ := strconv.ParseBool(labelVal)
			fv.SetBool(b)
		case reflect.Int, reflect.Int64:
			n, _ := strconv.ParseInt(labelVal, 10, 64)
			fv.SetInt(n)
		default:
			panic("unhandled default case")
		}
	}
}
