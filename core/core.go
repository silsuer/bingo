package core

import (
	"github.com/urfave/cli"
	"github.com/silsuer/bingo-router"
	"fmt"
	"github.com/silsuer/bingo/settings"
)

type Bingo struct {
	Cli    *cli.App
	Router *bingo_router.Router
}

// 创建一个app
func NewApp() *Bingo {
	b := &Bingo{
		Cli:    cli.NewApp(),
		Router: bingo_router.New(),
	}

	// 给cli添加一个bingo命令，然后在上面
	b.Cli.Name = "bingo"
	b.Cli.Usage = "bingo cli"
	b.Cli.Version = settings.VERSION
	b.Cli.Commands = []cli.Command{
		{
			Name: "run",
			Action: func(c *cli.Context) error {
				// 打印函数
				fmt.Println("Start a http server, try 'bingo run dev' ,'bingo run watch' or 'bingo run production'")
				return nil
			},
			Subcommands: []cli.Command{
				{
					Name:  "dev",
					Usage: "Start a development server.",
					Action: func(c *cli.Context) error {

						// 开启一个测试服务器

						fmt.Println("dev")
						return nil
					},
				},
				{
					Name:  "watch",
					Usage: "Listen the current directory and auto restart dev server.",
					Action: func(c *cli.Context) error {

						// 监听工程目录，有变化时自动重启

						return nil
					},
				},
				{
					Name:  "production",
					Usage: "Start a production environment server.",
					Action: func(c *cli.Context) error {
						// 开启生产环境服务器，要检测各种变量，检测数据等
						return nil
					},
				},
			},
		},
	}

	return b
}
