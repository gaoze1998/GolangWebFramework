package Distrabution

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"sync"
	"time"
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
	SP.ServerNameNextIndex[args.Name] = 0
	*reply = 1
	return nil
}

//ServerPool服务提供者集合
type ServerPool struct {
	ServerNameBackendsMap map[string][]*Backend
	ServerNameNextIndex   map[string]int
}

func ping(Url string) bool {
	timeout := 2 * time.Second
	conn, err := net.DialTimeout("tcp", Url, timeout)
	if err != nil {
		fmt.Println("Site unreachable, error: ", err)
		return false
	}
	_ = conn.Close() // 不需要维护连接，把它关闭
	return true
}

//CheckBackendHealth检查服务提供者的健康
func CheckBackendHealth() {
	for {
		time.Sleep(time.Minute * 10)
		for _, v := range SP.ServerNameBackendsMap {
			for _, vv := range v {
				alive := ping(vv.URL.Host)
				vv.SetAlive(alive)
			}
		}
	}
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
	ABackendIndex := SP.ServerNameNextIndex[name]
	for ; ABackendIndex < len(backends); ABackendIndex++ {
		if backends[ABackendIndex].IsAlive() {
			SP.ServerNameNextIndex[name] = (ABackendIndex + 1) % len(backends)
			return backends[ABackendIndex]
		}
	}
	SP.ServerNameNextIndex[name] = 0
	return nil
}

var SP ServerPool

// ServerPoolInit初始化服务提供者集合
func ServerPoolInit() {
	SP.ServerNameBackendsMap = make(map[string][]*Backend)
	SP.ServerNameNextIndex = make(map[string]int)
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
