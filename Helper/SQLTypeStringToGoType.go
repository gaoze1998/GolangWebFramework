package Helper

import "reflect"

// 转换SQL类型字符串到go语言类型
func SQLTypeStringToGoType(p string) reflect.Type {
	switch p {
	case "int":
		return reflect.TypeOf(1)
	case "varchar":
		return reflect.TypeOf("")
	}
	return nil
}
