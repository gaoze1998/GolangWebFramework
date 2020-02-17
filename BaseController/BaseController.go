package BaseController

import (
	"fmt"
	"net/http"
)

//Controller 控制器接口
type Controller interface {
	Get()
	Post()
	Put()
	Delete()
	SetReq(req *http.Request)
	SetRespWriter(respWriter http.ResponseWriter)
}

//BaseContruller 基础控制器结构体
type BaseController struct {
	Req        *http.Request
	RespWriter http.ResponseWriter
}

func (bc BaseController) Get() {
	fmt.Println("这是Get接口")
}

func (bc BaseController) Post() {
	fmt.Println("这是Post接口")
}

func (bc BaseController) Put() {
	fmt.Println("这是Put接口")
}

func (bc BaseController) Delete() {
	fmt.Println("这是Delete接口")
}

func (bc *BaseController) SetReq(req *http.Request) {
	bc.Req = req
}

func (bc *BaseController) SetRespWriter(respWriter http.ResponseWriter) {
	bc.RespWriter = respWriter
}
