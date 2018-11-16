package main

import (
	"github.com/silsuer/bingo/core"
	"os"
)

var App *core.Bingo

func init() {
	App = core.NewApp()
}

func main() {
	// 创建一个cli
	// 创建一个应用
	// 挂载路由
	App.Cli.Run(os.Args)
	//App.Router.Mount(routes.Api())
	//http.ListenAndServe(":8080", App.Router)
}
