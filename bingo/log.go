package bingo

import (
	"log"
)

// log 包
// 1. 设定一个全局变量用来输出日志
// 2. 日志自动分割
// 3. 日志包含多个错误级别
// 4. 日志对象配置（日志文件的位置，自动分割规则等）
// 5. 协程池异步输出日志

var Logger Log // 全局log

type Log struct {
	log.Logger
	synchronize bool // 日志记录模式 0 同步记录 2 协程池记录
}

func (l *Log) set() {

}
