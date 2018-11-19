package main

import (
	"github.com/silsuer/bingo"
	"os"
	"${path}/routes"
)

func main() {
	// 把env文件路径传入进去
	p, _ := os.Getwd()
	app := bingo.NewApp(p + "/.env.yml")
	app.Router.Mount(routes.Api())
	app.Cli.Run(os.Args)
}
