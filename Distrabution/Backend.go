package Distrabution

import (
	"fmt"
	"net/http/httputil"
	"net/rpc"
	"net/url"
	"sync"
)

//Backend服务提供者
type Backend struct {
	URL          *url.URL
	Alive        bool
	mux          sync.RWMutex
	ReverseProxy *httputil.ReverseProxy
}

//SetAlive修改存活状态
func (b *Backend) SetAlive(alive bool) {
	b.mux.Lock()
	b.Alive = alive
	b.mux.Unlock()
}

// 如果后端还活着，IsAlive 返回 true
func (b *Backend) IsAlive() (alive bool) {
	b.mux.RLock()
	alive = b.Alive
	b.mux.RUnlock()
	return
}

func AddBackendToRegistry(registryAddr string, servername string, backendUrl *url.URL) {
	conn, err := rpc.DialHTTP("tcp", registryAddr)
	if err != nil {
		fmt.Println(err)
		return
	}

	req := AddBackEndArgs{
		Name: servername,
		URL:  backendUrl,
	}
	var res int
	err = conn.Call("AddServerI.AddBackEnd", req, &res)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("注册结果是 %d\n", res)
}

func RegistBackend(serverAddr string, serverName string, registryAddr string) {
	if serverAddr[0] == ':' {
		serverAddr = "localhost" + serverAddr
	}
	serverAddr = "http://" + serverAddr
	backendUrl, err := url.Parse(serverAddr)
	if err != nil {
		fmt.Println("解析地址失败: ", err)
		return
	}
	AddBackendToRegistry(registryAddr, serverName, backendUrl)
}
