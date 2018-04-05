package core

import (
	"bingo/app/controller"
	"bingo/app/middleware"
)

// 这里记录所有的应该注册的结构体
// 控制器map
var ControllerMap map[string]interface{}
// 中间件map
var MiddlewareMap map[string]interface{}

func init()  {
	ControllerMap = make(map[string]interface{})
	MiddlewareMap = make(map[string]interface{})
	// 给这两个map赋初始值 每次添加完一条路由或中间件，都要在此处把路由或者中间件注册到这里
	// 注册中间件
    MiddlewareMap["WebMiddleware"] =&middleware.WebMiddleware{}

	// 注册路由
	ControllerMap["Controller"] = &controller.Controller{}
}

