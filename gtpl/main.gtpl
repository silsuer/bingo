package main

import (
	"os"
	"${path}/routes"
	"${path}/core"
)

func main() {
	// 把env文件路径传入进去
	p, _ := os.Getwd()
	app := core.NewApp(p + "/.env.yml")
	// 挂载路由
	app.Router.Mount(routes.Api())
	// 挂载命令
	app.Cli.Run(os.Args)
}
