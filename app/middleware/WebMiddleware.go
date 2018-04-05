package middleware

import (
	"net/http"
)

// 定义一个结构体
type WebMiddleware struct {
	ResString    string   // 响应的字符串
	Status       int      // 响应的状态码
}

// 在结构体上挂载一个handle方法
func (m *WebMiddleware) Handle(w http.ResponseWriter,r *http.Request) (http.ResponseWriter,*http.Request) {
	// 首先定义一个为空的响应，这个响应默认状态码是200，如果不是200，会返回一个错误
	// 返回的数据
	//fmt.Fprintln(w,"<h1>调用了中间件！</h1>")
	//m.ResString = "调用错误"
    // 如果响应没有问题，就返回这两个数据，如果有问题，就重新赋值这个中间件的两个字段,默认抛出500错误
	return w,r   // 返回这些数据，就是继续往下走，否则返回
}