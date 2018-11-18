package main

import (
	"github.com/silsuer/bingo/core"
	"github.com/silsuer/bingo/routes"
	"os"
)

var App *core.Bingo

func init() {
	App = core.NewApp()
}

func main() {
	// 挂载路由
	App.Router.Mount(routes.Api())
	core.Singleton = App // 添加单例
	App.Cli.Run(os.Args)

}
