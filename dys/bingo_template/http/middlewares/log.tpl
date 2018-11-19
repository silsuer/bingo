package middlewares

import (
	"github.com/silsuer/bingo-router"
	"${path}/utils"
	"fmt"
)

func Log(c *bingo_router.Context, next func(c *bingo_router.Context)) {
	defer func() {
		err := recover()
		if err != nil {
			utils.Log.Fatal(fmt.Sprint(c.Request, err))
		}
	}()
	go func() {
		// 收尾时要进行操作
		// 记录日志,默认会记录所有的请求与响应，响应时上下文应该提供一个函数，向其中打印数据
		utils.Log.Info(fmt.Sprint(c.Request))
	}()
	next(c)
}
