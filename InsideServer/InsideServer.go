package InsideServer

import (
	"fmt"
	"github.com/gaoze1998/GolangWebFramwork/BaseRouter"
	"github.com/gaoze1998/GolangWebFramwork/BaseSession"
	"net/http"
)

//InsideServer 内嵌服务器结构体
type InsideServer struct {
	Addr       string
	Req        *http.Request
	RespWriter http.ResponseWriter
	Router     *BaseRouter.Router
	SessionOn  bool
}

//InsideServerBaseConfig 内嵌服务器基础配置结构
type BaseConfig struct {
	Addr   string
	Router *BaseRouter.Router
	Son    bool
}

//ConfigInit 应用程序配置
func ConfigInit() BaseConfig {
	bc := BaseConfig{}
	bc.Addr = ":8081"
	bc.Son = true
	return bc
}

//InsideServer.Config 内嵌服务器配置
func (is *InsideServer) Config(config BaseConfig) {
	is.Addr = config.Addr
	is.Router = config.Router
	is.SessionOn = config.Son
	if config.Son {
		BaseSession.BaseMemorySession = make(map[int]*BaseSession.BaseSession)
	}
}

//InsideServer.ServeHTTP 默认路由
func (is *InsideServer) ServeHTTP(responseWriter http.ResponseWriter, r *http.Request) {
	for path, controller := range is.Router.PathToController {
		if path != r.URL.Path {
			continue
		}
		if r.Method == "GET" {
			controller.SetReq(r)
			controller.SetRespWriter(responseWriter)
			controller.Get()
		} else if r.Method == "POST" {
			controller.SetReq(r)
			controller.SetRespWriter(responseWriter)
			controller.Post()
		} else if r.Method == "DELETE" {
			controller.SetReq(r)
			controller.SetRespWriter(responseWriter)
			controller.Delete()
		} else if r.Method == "PUT" {
			controller.SetReq(r)
			controller.SetRespWriter(responseWriter)
			controller.Put()
		}
	}
}

//InsideServer.Serve 开启内嵌服务器
func (is *InsideServer) Serve() {
	err := http.ListenAndServe(is.Addr, is)
	if err != nil {
		fmt.Println("无法启动服务器，可能原因有：地址被占用；服务器资源分配空缺。\n解决方案：更改地址；重新分配服务器资源。")
	}
}
