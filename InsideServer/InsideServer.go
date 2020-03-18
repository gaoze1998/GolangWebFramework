package InsideServer

import (
	"fmt"
	"github.com/gaoze1998/GolangWebFramework/BaseRouter"
	"github.com/gaoze1998/GolangWebFramework/BaseSession"
	"github.com/gaoze1998/GolangWebFramework/Distrabution"
	"log"
	"net"
	"net/http"
	"net/rpc"
)

//InsideServer 内嵌服务器结构体
type InsideServer struct {
	Addr         string
	Req          *http.Request
	RespWriter   http.ResponseWriter
	Router       *BaseRouter.Router
	SessionOn    bool
	RegistryOn   bool
	RegistryAddr string
}

//InsideServerBaseConfig 内嵌服务器基础配置结构
type BaseConfig struct {
	Addr         string
	Router       *BaseRouter.Router
	Son          bool
	Registry     bool
	RegistryAddr string
}

//ConfigInit 应用程序配置
func ConfigInit() BaseConfig {
	bc := BaseConfig{}
	bc.Addr = ":8080"
	bc.Son = true
	bc.Registry = false
	return bc
}

//InsideServer.Config 内嵌服务器配置
func (is *InsideServer) Config(config BaseConfig) {
	is.Addr = config.Addr
	is.Router = config.Router
	is.SessionOn = config.Son
	is.RegistryOn = config.Registry
	is.RegistryAddr = config.RegistryAddr
	if config.Son {
		BaseSession.BaseMemorySession = make(map[int]*BaseSession.BaseSession)
	}
}

//InsideServer.ServeHTTP 默认路由
func (is *InsideServer) ServeHTTP(responseWriter http.ResponseWriter, r *http.Request) {
	orgin := r.Header.Get("Access-Control-Allow-Origin")
	if orgin == "" {
		orgin = r.Header.Get("Origin")
	}
	responseWriter.Header().Set("Access-Control-Allow-Origin", orgin)
	responseWriter.Header().Add("Access-Control-Allow-Headers", "X-Requested-With,Content-Type,x-csrftoken") //header的类型
	responseWriter.Header().Set("Access-Control-Max-Age", "86400")
	responseWriter.Header().Set("Access-Control-Allow-Methods", "*")
	responseWriter.Header().Set("Access-Control-Allow-Credentials", "true")
	responseWriter.Header().Set("content-type", "application/json;charset=UTF-8") //返回数据格式是json
	if is.RegistryOn {
		Distrabution.LB(responseWriter, r)
		return
	}
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
	if is.RegistryOn {
		Distrabution.ServerPoolInit()
		rpc.Register(new(Distrabution.AddServerI)) // 注册rpc服务
		rpc.HandleHTTP()                           // 采用http协议作为rpc载体

		lis, err := net.Listen("tcp", is.RegistryAddr)
		if err != nil {
			log.Fatalln("fatal error: ", err)
		}

		go http.Serve(lis, nil)
		go Distrabution.CheckBackendHealth()
		http.ListenAndServe(is.Addr, http.HandlerFunc(Distrabution.LB))
	} else {
		err := http.ListenAndServe(is.Addr, is)
		if err != nil {
			fmt.Println("无法启动服务器，可能原因有：地址被占用；服务器资源分配空缺。\n" +
				"解决方案：更改地址；重新分配服务器资源。")
		}
	}
}
