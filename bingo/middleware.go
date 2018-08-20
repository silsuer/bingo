package bingo

type MiddlewareInterface interface {
	Handle(c *Context) *Context
	Init()
	SyncStatus() bool
}

// 中间件父类
type Middleware struct {
	Sync bool // 是否是同步中间件 true为同步 false为异步
	//handle func(c *Context)
}

func (m *Middleware) Handle(c *Context) *Context {
	return c
}

// 初始化中间件，比如设置同步/异步执行
func (m *Middleware) Init() {
	m.Sync = true
}

// 返回是否是同步中间件
func (m *Middleware) SyncStatus() bool {
	if m.Sync == true {
		return true
	} else {
		return false
	}
}

var GlobalMiddlewares []MiddlewareInterface
var GroupMiddlewares map[string][]MiddlewareInterface
var RouteMiddlewares map[string]MiddlewareInterface

// 注册全局路由
func GlobalMiddlewareRegister(middlewares []MiddlewareInterface) {
	// 根据传入的中间件数组，将其赋值给全局中间件变量，在启动的时候，遍历这个全局变量，将其分别挂在到每个路由上
	GlobalMiddlewares = middlewares
}

// 注册路由组
func GroupMiddlewareRegister(middlewares map[string][]MiddlewareInterface) {
	// 赋值给全局路由组中间件比那里，在启动的时候分别挂载
	GroupMiddlewares = make(map[string][]MiddlewareInterface)
	GroupMiddlewares = middlewares
}

// 注册单个路由
func RouteMiddlewareRegister(middlewares map[string]MiddlewareInterface) {
	RouteMiddlewares = make(map[string]MiddlewareInterface)
	RouteMiddlewares = middlewares
}
