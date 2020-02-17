package Helper

import (
	"database/sql"
	"reflect"
	"strconv"
)

// SetElemValues辅助设置各种类型的值
func SetElemValues(elem *reflect.Value, values []sql.RawBytes) {
	var value string
	for i, col := range values {
		value = string(col)
		if elem.Field(i).Kind() == reflect.Int {
			v, _ := strconv.Atoi(value)
			elem.Field(i).SetInt(int64(v))
		} else if elem.Field(i).Kind() == reflect.String {
			elem.Field(i).SetString(value)
		}
	}
}