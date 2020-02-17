package BaseRouter

import (
	"GolangWebFramwork/BaseController"
)

//RouterTable 路由表结构
type RouterTable struct {
	PathToController map[string]BaseController.Controller
}

//Router 路由器结构
type Router struct {
	RouterTable
}

//Router.Register 路由注册
func (r *Router) Register(path string, c BaseController.Controller) {
	r.PathToController[path] = c
}

//Init 生成并初试化一个路由器
func Init() *Router {
	router := &Router{}
	router.PathToController = make(map[string]BaseController.Controller)
	return router
}
