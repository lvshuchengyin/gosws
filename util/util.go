// util
package util

import "reflect"

func CloneValue(source interface{}) (destin interface{}) {
	x := reflect.ValueOf(source)
	if x.Kind() == reflect.Ptr {
		starX := x.Elem()
		y := reflect.New(starX.Type())
		starY := y.Elem()
		starY.Set(starX)
		//reflect.ValueOf(destin).Elem().Set(y.Elem())
		//destin = starY.Interface()
		destin = y.Interface()
	} else {
		destin = x.Interface()
	}

	return destin
}
