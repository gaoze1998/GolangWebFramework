package Helper

import (
	"reflect"
	"strconv"
)

// ValueToString 转化Value为string
func ValueToString(elem reflect.Value) string {
	if elem.Kind() == reflect.Int {
		return strconv.Itoa(elem.Interface().(int))
	} else if elem.Kind() == reflect.String {
		return elem.String()
	}
	return ""
}
