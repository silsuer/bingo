package bingo

import (
	"net/http"
	"fmt"
	"github.com/gorilla/context"
)

// bingo结构体，向外暴露一些属性和方法  实现了http方法
type Bingo struct{}

func (b *Bingo) Run(port string) {
	// 根据httprouter进行重写(根据Httprouter的原理，重新实现路由)
	// 这个时候要根据RouteList,对每一个方法解析出一个tree来
	router := New()
    // 开始把路由列表注册到tree中
    for _,v:= range RouteList {
    	router.Handle(v.Method,v.Path,v)
	}
	fmt.Println("Bingo Running......")
	// 静态页面
	router.NotFound = func(writer http.ResponseWriter, request *http.Request) {
		http.ServeFile(writer,request,GetPublicPath()+request.URL.Path)
	}
	// 开启服务器
	err := http.ListenAndServe(port, context.ClearHandler(router))
	if err != nil {
		fmt.Println(err)
	}
	// TODO 监听平滑升级和重启
}
