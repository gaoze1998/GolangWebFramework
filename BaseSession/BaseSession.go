package BaseSession

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

//SessionInterface Session接口，只要符合此接口，可不必修改应用
type SessionInterface interface {
	Create()
	Save(w http.ResponseWriter)
	Read(req *http.Request)
	Delete()
}

//BaseSession 基础Session模型
type BaseSession struct {
	sid        int
	KV         map[string]interface{}
	ExpireTime time.Duration
}

//BaseMemorySession 基础内存型Session存储
var BaseMemorySession map[int]*BaseSession

func (baseSession *BaseSession) Create() {
	baseSession.sid = rand.Int()
	baseSession.KV = make(map[string]interface{})
	baseSession.ExpireTime = 30 * time.Minute
}

func (baseSession BaseSession) Save(w http.ResponseWriter) {
	v := strconv.FormatInt(int64(baseSession.sid), 10)
	if BaseMemorySession == nil {
		fmt.Println("Session并未开启")
		return
	}
	BaseMemorySession[baseSession.sid] = &baseSession
	http.SetCookie(w, &http.Cookie{Name: "sid", Value: v})
}

func (baseSession *BaseSession) Read(req *http.Request) {
	var err error
	var cookies *http.Cookie
	var v int

	cookies, err = req.Cookie("sid")
	if err != nil {
		fmt.Println("读取Session出错！")
		return
	}

	v, err = strconv.Atoi(cookies.Value)
	if _, ok := BaseMemorySession[v]; !ok {
		fmt.Printf("Session %d 不存在\n", v)
		return
	}
	baseSession.ExpireTime = BaseMemorySession[v].ExpireTime
	baseSession.KV = BaseMemorySession[v].KV
	baseSession.sid = BaseMemorySession[v].sid
}

func (baseSession BaseSession) Delete() {
	if _, ok := BaseMemorySession[baseSession.sid]; !ok {
		fmt.Printf("Session %d 不存在\n", baseSession.sid)
		return
	}
	delete(BaseMemorySession, baseSession.sid)
}
