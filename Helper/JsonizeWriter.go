package Helper

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
)

// 将结果集rst json序列化，并响应请求
func JsonizeWriter(respw http.ResponseWriter, rst []reflect.Value) {
	fmt.Fprintf(respw, "[\n")
	for i, v := range rst{
		vv := v.Interface()
		b, err := json.Marshal(vv)
		if err != nil {
			fmt.Println("Json序列化错误\n")
		}
		//fmt.Println(string(b))
		if i != len(rst) - 1 {
			fmt.Fprintf(respw, "%s,\n", b)
		}else{
			fmt.Fprintf(respw, "%s\n", b)
		}
	}
	fmt.Fprintf(respw, "]")
}
