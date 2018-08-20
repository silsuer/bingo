// Copyright 2013 Julien Schmidt. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be found
// in the LICENSE file.

// Package httprouter is a trie based high performance HTTP request router.
//
// A trivial example is:
//
//  package main
//
//  import (
//      "fmt"
//      "github.com/julienschmidt/httprouter"
//      "net/http"
//      "log"
//  )
//
//  func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
//      fmt.Fprint(w, "Welcome!\n")
//  }
//
//  func Hello(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
//      fmt.Fprintf(w, "hello, %s!\n", ps.ByName("name"))
//  }
//
//  func main() {
//      router := httprouter.New()
//      router.GET("/", Index)
//      router.GET("/hello/:name", Hello)
//
//      log.Fatal(http.ListenAndServe(":8080", router))
//  }
//
// The router matches incoming requests by the request method and the path.
// If a handle is registered for this path and method, the router delegates the
// request to that function.
// For the methods GET, POST, PUT, PATCH and DELETE shortcut functions exist to
// register handles, for all other methods router.Handle can be used.
//
// The registered path, against which the router matches incoming requests, can
// contain two types of parameters:
//  Syntax    Type
//  :name     named parameter
//  *name     catch-all parameter
//
// Named parameters are dynamic path segments. They match anything until the
// next '/' or the path end:
//  Path: /blog/:category/:post
//
//  Requests:
//   /blog/go/request-routers            match: category="go", post="request-routers"
//   /blog/go/request-routers/           no match, but the router would redirect
//   /blog/go/                           no match
//   /blog/go/request-routers/comments   no match
//
// Catch-all parameters match anything until the path end, including the
// directory index (the '/' before the catch-all). Since they match anything
// until the end, catch-all paramerters must always be the final path element.
//  Path: /files/*filepath
//
//  Requests:
//   /files/                             match: filepath="/"
//   /files/LICENSE                      match: filepath="/LICENSE"
//   /files/templates/article.html       match: filepath="/templates/article.html"
//   /files                              no match, but the router would redirect
//
// The value of parameters is saved as a slice of the Param struct, consisting
// each of a key and a value. The slice is passed to the Handle func as a third
// parameter.
// There are two ways to retrieve the value of a parameter:
//  // by the name of the parameter
//  user := ps.ByName("user") // defined by :user or *user
//
//  // by the index of the parameter. This way you can also get the name (key)
//  thirdKey   := ps[2].Key   // the name of the 3rd parameter
//  thirdValue := ps[2].Value // the value of the 3rd parameter
package bingo

import (
	"net/http"
)

// Handle is a function that can be registered to a route to handle HTTP
// requests. Like http.HandlerFunc, but has a third parameter for the values of
// wildcards (variables).
// 这里重写了httprouter中的方法
type Handle func(http.ResponseWriter, *http.Request, Params)

// Param is a single URL parameter, consisting of a key and a value.
type Param struct {
	Key   string
	Value string
}

// Params is a Param-slice, as returned by the router.
// The slice is ordered, the first URL parameter is also the first slice value.
// It is therefore safe to read values by the index.
type Params []Param

// ByName returns the value of the first Param which key matches the given name.
// If no matching Param is found, an empty string is returned.
func (ps Params) ByName(name string) string {
	for i := range ps {
		if ps[i].Key == name {
			return ps[i].Value
		}
	}
	return ""
}

// Router is a http.Handler which can be used to dispatch requests to different
// handler functions via configurable routes
type Router struct {
	trees map[string]*node

	// Enables automatic redirection if the current route can't be matched but a
	// handler for the path with (without) the trailing slash exists.
	// For example if /foo/ is requested but a route only exists for /foo, the
	// client is redirected to /foo with http status code 301 for GET requests
	// and 307 for all other request methods.
	RedirectTrailingSlash bool

	// If enabled, the router tries to fix the current request path, if no
	// handle is registered for it.
	// First superfluous path elements like ../ or // are removed.
	// Afterwards the router does a case-insensitive lookup of the cleaned path.
	// If a handle can be found for this route, the router makes a redirection
	// to the corrected path with status code 301 for GET requests and 307 for
	// all other request methods.
	// For example /FOO and /..//Foo could be redirected to /foo.
	// RedirectTrailingSlash is independent of this option.
	RedirectFixedPath bool

	// If enabled, the router checks if another method is allowed for the
	// current route, if the current request can not be routed.
	// If this is the case, the request is answered with 'Method Not Allowed'
	// and HTTP status code 405.
	// If no other Method is allowed, the request is delegated to the NotFound
	// handler.
	HandleMethodNotAllowed bool

	// Configurable http.HandlerFunc which is called when no matching route is
	// found. If it is not set, http.NotFound is used.
	NotFound http.HandlerFunc

	// Configurable http.HandlerFunc which is called when a request
	// cannot be routed and HandleMethodNotAllowed is true.
	// If it is not set, http.Error with http.StatusMethodNotAllowed is used.
	MethodNotAllowed http.HandlerFunc

	// Function to handle panics recovered from http handlers.
	// It should be used to generate a error page and return the http error code
	// 500 (Internal Server Error).
	// The handler can be used to keep your server from crashing because of
	// unrecovered panics.
	PanicHandler func(http.ResponseWriter, *http.Request, interface{})
}

// Make sure the Router conforms with the http.Handler interface
var _ http.Handler = New()

// New returns a new initialized Router.
// Path auto-correction, including trailing slashes, is enabled by default.
func New() *Router {
	return &Router{
		RedirectTrailingSlash:  true,
		RedirectFixedPath:      true,
		HandleMethodNotAllowed: true,
	}
}

// GET is a shortcut for router.Handle("GET", path, handle)
func (r *Router) GET(path string, route Route) {
	r.Handle("GET", path, route)
}

// HEAD is a shortcut for router.Handle("HEAD", path, handle)
func (r *Router) HEAD(path string, route Route) {
	r.Handle("HEAD", path, route)
}

// OPTIONS is a shortcut for router.Handle("OPTIONS", path, handle)
func (r *Router) OPTIONS(path string, route Route) {
	r.Handle("OPTIONS", path, route)
}

// POST is a shortcut for router.Handle("POST", path, handle)
func (r *Router) POST(path string, route Route) {
	r.Handle("POST", path, route)
}

// PUT is a shortcut for router.Handle("PUT", path, handle)
func (r *Router) PUT(path string, route Route) {
	r.Handle("PUT", path, route)
}

// PATCH is a shortcut for router.Handle("PATCH", path, handle)
func (r *Router) PATCH(path string, route Route) {
	r.Handle("PATCH", path, route)
}

// DELETE is a shortcut for router.Handle("DELETE", path, handle)
func (r *Router) DELETE(path string, route Route) {
	r.Handle("DELETE", path, route)
}

// Handle registers a new request handle with the given path and method.
//
// For GET, POST, PUT, PATCH and DELETE requests the respective shortcut
// functions can be used.
//
// This function is intended for bulk loading and to allow the usage of less
// frequently used, non-standardized or custom methods (e.g. for internal
// communication with a proxy).
func (r *Router) Handle(method, path string, route Route) {
	if path[0] != '/' {
		panic("path must begin with '/' in path '" + path + "'")
	}

	if r.trees == nil {
		r.trees = make(map[string]*node)
	}

	root := r.trees[method]
	if root == nil {
		root = new(node)
		r.trees[method] = root
	}

	root.addRoute(path, route)
}

// Handler is an adapter which allows the usage of an http.Handler as a
// request handle.
func (r *Router) Handler(method, path string, handler http.Handler) {
	route := Route{}
	route.TargetMethod = func(context *Context) {
		handler.ServeHTTP(context.Writer, context.Request)
	}
	r.Handle(method, path, route)
}

// HandlerFunc is an adapter which allows the usage of an http.HandlerFunc as a
// request handle.
func (r *Router) HandlerFunc(method, path string, handler http.HandlerFunc) {
	r.Handler(method, path, handler)
}

// ServeFiles serves files from the given file system root.
// The path must end with "/*filepath", files are then served from the local
// path /defined/root/dir/*filepath.
// For example if root is "/etc" and *filepath is "passwd", the local file
// "/etc/passwd" would be served.
// Internally a http.FileServer is used, therefore http.NotFound is used instead
// of the Router's NotFound handler.
// To use the operating system's file system implementation,
// use http.Dir:
//     router.ServeFiles("/src/*filepath", http.Dir("/var/www"))
func (r *Router) ServeFiles(path string, root http.FileSystem) {
	if len(path) < 10 || path[len(path)-10:] != "/*filepath" {
		panic("path must end with /*filepath in path '" + path + "'")
	}

	fileServer := http.FileServer(root)

	route := Route{}
	route.TargetMethod = func(context *Context) {
		context.Request.URL.Path = context.Params.ByName("filepath")
		fileServer.ServeHTTP(context.Writer, context.Request)
	}

	r.GET(path, route)
}

func (r *Router) recv(w http.ResponseWriter, req *http.Request) {
	if rcv := recover(); rcv != nil {
		r.PanicHandler(w, req, rcv)
	}
}

// Lookup allows the manual lookup of a method + path combo.
// This is e.g. useful to build a framework around this router.
// If the path was found, it returns the handle function and the path parameter
// values. Otherwise the third return value indicates whether a redirection to
// the same path with an extra / without the trailing slash should be performed.
func (r *Router) Lookup(method, path string) (Route, Params, bool) {
	if root := r.trees[method]; root != nil {
		return root.getValue(path)
	}
	return Route{}, nil, false
}

// ServeHTTP makes the router implement the http.Handler interface.
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if r.PanicHandler != nil {
		defer r.recv(w, req)
	}

	// 在查找之前，要先看看是否存在中间件
	// 注册路由的时候，应该把中间件也放在此处

	// 开始去查找注册的路由函数
	if root := r.trees[req.Method]; root != nil {
		path := req.URL.Path

		if route, ps, tsr := root.getValue(path); route.TargetMethod != nil {
			// 封装上下文
			session, _ := globalSession.Get(req, "bingoSess") // 从cookie读取数据并且返回对应session

			//fmt.Println(err)
			context := &Context{w, req, ps, Session{session: session, writer: w, req: req}}
			// 判断路由是否有中间件列表，如果有，就执行
			if len(route.Middleware) != 0 {
				for _, m := range route.Middleware {
					if m.SyncStatus() == true { // 同步中间件
						context = m.Handle(context) // 顺序执行中间件，得到的返回结果重新注入到上下文中
					} else {
						go m.Handle(context) // 开启一个协程异步执行中间件
					}
				}
			}
			// 执行目标函数
			route.TargetMethod(context)
			return
		} else if req.Method != "CONNECT" && path != "/" {
			code := 301 // Permanent redirect, request with GET method
			if req.Method != "GET" {
				// Temporary redirect, request with same method
				// As of Go 1.3, Go does not support status code 308.
				code = 307
			}

			if tsr && r.RedirectTrailingSlash {
				if len(path) > 1 && path[len(path)-1] == '/' {
					req.URL.Path = path[:len(path)-1]
				} else {
					req.URL.Path = path + "/"
				}
				http.Redirect(w, req, req.URL.String(), code)
				return
			}

			// Try to fix the request path
			if r.RedirectFixedPath {
				fixedPath, found := root.findCaseInsensitivePath(
					CleanPath(path),
					r.RedirectTrailingSlash,
				)
				if found {
					req.URL.Path = string(fixedPath)
					http.Redirect(w, req, req.URL.String(), code)
					return
				}
			}
		}
	}

	// Handle 405
	if r.HandleMethodNotAllowed {
		for method := range r.trees {
			// Skip the requested method - we already tried this one
			if method == req.Method {
				continue
			}

			route, _, _ := r.trees[method].getValue(req.URL.Path)
			if route.Target != nil {
				if r.MethodNotAllowed != nil {
					r.MethodNotAllowed(w, req)
				} else {
					http.Error(w,
						http.StatusText(http.StatusMethodNotAllowed),
						http.StatusMethodNotAllowed,
					)
				}
				return
			}
		}
	}

	// Handle 404
	if r.NotFound != nil {
		r.NotFound(w, req)
	} else {
		http.NotFound(w, req)
	}
}

const GET = "GET"
const POST = "POST"
const DELETE = "DELETE"
const PUT = "PUT"
const PATCH = "PATCH"

//const PUTCH  = "PUTCH"
//const put

// 上下文结构体
type Context struct {
	Writer  http.ResponseWriter // 响应
	Request *http.Request       // 请求
	Params  Params              //参数
	Session Session             // 保存session
}

//var store = sessions.NewCookieStore（[] byte（“something-very-secret”））

type TargetHandle func(c *Context)

type MiddlewareHandle func(c *Context) *Context // 中间件需要把上下文返回回来，用来传入TargetHandle中

type Route struct {
	Prefix       string           // 路由前缀
	Path         string           // 路径
	TargetMethod func(c *Context) // 要执行的方法
	//Controller *ControllerInterface  // 路由对应的控制器
	Method     string                // 访问类型 是get post 或者其他
	Alias      string                // 路由的别名，并没有什么卵用的样子.......
	Middleware []MiddlewareInterface // 中间件名称，在执行TargetHandle之前执行的方法
	//Middleware []reflect.Value // 中间件名称，在执行TargetHandle之前执行的方法
	MGroup []string // 中间件组，这里记录中间件组的名称，用于在初始化的时候添加中间件
}

var RouteList []Route // 全体列表

// 注册路由
func RegisterRoute(r []Route) {
	for _, v := range r {
		RouteList = append(RouteList, v) // 把要注册的路由放置到路由列表中
	}
}

// 添加路由时需要，设置为Get方法
func (r Route) Get(path string) Route {
	//return r.Request(GET, path, target)
	r.Path = path
	r.Method = GET
	return r
}

// 这里传入一个回调
func (r Route) Target(target TargetHandle) Route {
	return r.Request(r.Method, r.Path, target)
}

// 添加路由时需要，设置为Post方法
func (r Route) Post(path string, target TargetHandle) Route {
	//return r.Request(POST, path, target)
	r.Path = path
	r.Method = POST
	return r
}

// 添加路由时需要，设置为put方法
func (r Route) Put(path string, target TargetHandle) Route {
	//return r.Request(PUT, path, target)
	r.Path = path
	r.Method = PUT
	return r
}

// 添加路由时需要，设置为patch方法
func (r Route) Patch(path string, target TargetHandle) Route {
	//return r.Request(PATCH, path, target)
	r.Path = path
	r.Method = PATCH
	return r
}

// 添加路由时需要，设置为delete方法
func (r Route) Delete(path string, target TargetHandle) Route {
	//return r.Request(DELETE, path, target)
	r.Path = path
	r.Method = DELETE
	return r
}

func (r Route) Request(method string, path string, target TargetHandle) Route {
	r.Method = method
	r.Path = path
	r.TargetMethod = target
	return r
}

func (r Route) MiddlewareGroup(groupName []string) Route {
	r.MGroup = groupName
	return r
}

// 通过传入在Kernel中定义好的路由名称，来设置路由
func (r Route) MiddlewareWithName(names []string) Route {
	for _, name := range names {
		if m, ok := RouteMiddlewares[name]; ok {
			r.addMiddleware(m)
		}
	}
	return r
}

// 通过直接传入中间件数组来设置路由
func (r Route) Middlewares(ms []MiddlewareInterface) Route {
	for _, m := range ms {
		r.addMiddleware(m)
	}
	return r
}

func (r Route) Register() Route {
	// 将全局中间件挂载上去
	for _, m := range GlobalMiddlewares {
		r.addMiddleware(m) // 添加全局中间件
	}

	// 将中间件组挂载上去
	for _, mName := range r.MGroup {
		if mg, ok := GroupMiddlewares[mName]; ok { // 存在这个路由组
			for _, m := range mg { // 得到路由
				r.addMiddleware(m)
			}
		}
	}
	// 将单个中间件挂载上去
	RouteList = append(RouteList, r)
	return r
}

func (r *Route) addMiddleware(mid MiddlewareInterface) {
	flag := false // 当前路由上不存在这个中间件
	for _, m := range r.Middleware {
		if m == mid {
			flag = true // 存在中间件
		}
	}

	if flag == false { // 如果存在，则跳过，否则添加
		mid.Init()
		r.Middleware = append(r.Middleware, mid)
	}
}

func NewRoute() Route {
	return Route{}
}

// 路由组
type RouteGroup struct {
	PrefixString string                // 路由前缀
	MGName       []string              // 中间件组名
	MGroup       []MiddlewareInterface // 中间件集合
}

// 新建路由组
func NewRouteGroup() *RouteGroup {
	return &RouteGroup{}
}

func (r *RouteGroup) MiddlewareGroup(args ...string) *RouteGroup {
	for _, mName := range args {
		if ms, ok := GroupMiddlewares[mName]; ok {
			r.MGName = append(r.MGName, mName) // 录入中间件组名
			for _, m := range ms {
				m.Init()
				r.MGroup = append(r.MGroup, m)
			}
		}
	}
	return r
}

// 传入路由组名称，根据路由组的名字设定中间件
// NewRouteGroup().MiddlewareGroup('aa','bb','cc').Group(func()).Register()
func (r *RouteGroup) Group(call func(routes []Route) []Route) *RouteGroup {
	rs := make([]Route, 0)
	rs = call(rs)

	// 遍历路由并注册
	for _, route := range rs {
		route.Path = r.PrefixString + route.Path
		route.MiddlewareGroup(r.MGName).Register()
	}

	// 中间件组
	return r
}

// 为路由组中的所有路由加入前缀
func (r *RouteGroup) Prefix(prefix string) *RouteGroup {
	r.PrefixString = prefix
	return r
}
