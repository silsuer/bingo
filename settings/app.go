package settings

import (
	"os"
)

const (
	VERSION    = "1.0.0" // 版本号
	Production = iota
	Testing
	Development
)

var Env = Development // 当前环境

var Debug = true // 调试模式

var Local = "http://localhost" // 调试模式下的数据

var Domain = "http://localhost" // 当前访问地址

var Root string

var HttpPort = 8080 // 开启服务器的时候监听的端口号

// 设置一些默认值
func init() {
	Root, _ = os.Getwd()
}
