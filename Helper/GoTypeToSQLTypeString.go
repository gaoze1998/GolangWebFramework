package Helper

import "reflect"

// 转化Go语言类型至sql相应类型字符串
func GoTypeToSQLTypeString(p reflect.Type) string {
	switch p.Kind() {
	case reflect.Int:
		return "int"
	case reflect.String:
		return "varchar(20)"
	}
	return ""
}
