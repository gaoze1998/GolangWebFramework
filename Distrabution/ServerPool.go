package Distrabution

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"sync"
)

type AddServerI struct {
}

//AddBackEnd添加后端
func (s *AddServerI) AddBackEnd(args AddBackEndArgs, reply *int) error {
	fmt.Println("添加一个后端: ", args.URL)
	backend := Backend{
		URL:          args.URL,
		Alive:        true,
		mux:          sync.RWMutex{},
		ReverseProxy: httputil.NewSingleHostReverseProxy(args.URL),
	}
	SP.ServerNameBackendsMap[args.Name] = append(SP.ServerNameBackendsMap[args.Name], &backend)
	*reply = 1
	return nil
}

//ServerPool服务提供者集合
type ServerPool struct {
	ServerNameBackendsMap map[string][]*Backend
}

//AddBackEndArgs 添加后端参数
type AddBackEndArgs struct {
	Name string
	URL  *url.URL
}

//GetNextPeer返回下一个可用的服务器
func GetNextPeer(name string) *Backend {
	// 遍历后端列表，找到可用的服务器
	backends := SP.ServerNameBackendsMap[name]
	for _, v := range backends {
		if v.IsAlive() {
			return v
		}
	}
	return nil
}

var SP ServerPool

// ServerPoolInit初始化服务提供者集合
func ServerPoolInit() {
	SP.ServerNameBackendsMap = make(map[string][]*Backend)
}

// LB 对入向请求进行负载均衡
func LB(w http.ResponseWriter, r *http.Request) {
	serverName := strings.Split(r.URL.Path, "/")[1]
	peer := GetNextPeer(serverName)
	if peer != nil {
		peer.ReverseProxy.ServeHTTP(w, r)
		return
	}
	http.Error(w, "Service not available", http.StatusServiceUnavailable)
}
