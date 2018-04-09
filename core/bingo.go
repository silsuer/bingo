package core

import (
	"net/http"
	"reflect"
	"fmt"
)

// bingo结构体，向外暴露一些属性和方法  实现了http方法
type Bingo struct{}

func (b *Bingo) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	flag := false
	// 每一个http请求都会走到这里，然后在这里，根据请求的URL，为其分配所需要调用的方法
	params := []reflect.Value{reflect.ValueOf(w), reflect.ValueOf(r)}
	for _, v := range RoutesList {
		// 检测中间件，根据中间件首先开启中间件，然后再注册其他路由
		// 检测路由，根据路由指向需要的数据·
		if r.URL.Path == v.path && r.Method == v.method {
              flag = true   // 寻找到了对应路由，无需使用静态服务器
			//TODO 调用一个公共中间件，在这个中间件中寻找路由以及调用中间件收尾等功能

			// 检测该路由中是否存在中间件，如果存在，顺序调用
			for _, m := range v.middleware {
				if mid, ok := MiddlewareMap[m]; ok { // 判断是否注册了这个中间件
					rmid := reflect.ValueOf(mid)
					params = rmid.MethodByName("Handle").Call(params) // 执行中间件，返回values数组
					// 判断中间件执行结果，是否还要继续往下走
					str := rmid.Elem().FieldByName("ResString").String()
					if str != "" {
						status := rmid.Elem().FieldByName("Status").Int()
						// 字符串不空，查看状态码，默认返回500错误
						if status == 0 {
							status = 500
						}
						w.WriteHeader(int(status))
						fmt.Fprint(w,str)
						return
					}
				}
			}
			// 检测成功，开始调用方法
			// 获取一个控制器包下的结构体
			if d, ok := ControllerMap[v.controller]; ok { // 存在  c为结构体，调用c上挂载的方法
				reflect.ValueOf(d).MethodByName(v.function).Call(params)
			}
			// 停止向后执行
			return
		}
	}

	// 如果路由列表中还是没有的话,去静态服务器中寻找
	if !flag {
         // 去静态目录中寻找
         http.ServeFile(w,r,GetPublicPath()+ r.URL.Path)
	}
	return
}
func (b *Bingo) Run(port string) {
	// 传入一个端口号，没有返回值，根据端口号开启http监听
	// 此处要进行资源初始化，加载所有路由、配置文件等等
	RouteInit()
	// 实例化router类，这个类去获取所有router目录下的json文件，然后根据json中的配置，加载数据
	// 实例化env文件和config文件夹下的所有数据，根据配置
	// 根据路由列表，开始定义路由，并且根据端口号，开启http服务器
	http.ListenAndServe(port, b)

	// TODO 监听平滑升级和重启
}
