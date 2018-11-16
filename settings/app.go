package settings

const (
	VERSION     = "1.0.0"  // 版本号
	Production  = iota
	Testing
	Development
)

var Env = Development // 当前环境

var Debug = true // 调试模式

var HttpPort = 8080 // 开启服务器的时候监听的端口号
