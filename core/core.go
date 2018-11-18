package core

import (
	"fmt"
	"github.com/facebookgo/grace/gracehttp"
	"github.com/silsuer/bingo-router"
	"github.com/silsuer/bingo/settings"
	"github.com/silsuer/bingo/utils"
	"github.com/urfave/cli"
	"net/http"
	"os"
	"strconv"
)

const (
	EnvVariable = "BINGOENVPID"
)

var title = `
 ____    ___   _   _    ____    ___    _
| __ )  |_ _| | \ | |  / ___|  / _ \  | |
|  _ \   | |  |  \| | | |  _  | | | | | |
| |_) |  | |  | |\  | | |_| | | |_| | |_|
|____/  |___| |_| \_|  \____|  \___/  (_)
`

var Singleton *Bingo
var DevDomainUrl string
var NetworkUrl string

func init() {
	DevDomainUrl = settings.Local + ":" + strconv.Itoa(settings.HttpPort)
	network, _ := utils.IntranetIP()
	if len(network) > 0 {
		NetworkUrl = "http://" + network[0] + ":" + strconv.Itoa(settings.HttpPort)
	} else {
		NetworkUrl = DevDomainUrl
	}

}

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
				fmt.Println("Start a http server, try 'bingo run dev' or 'bingo run production'")
				return nil
			},
			Subcommands: []cli.Command{
				{
					Name:  "dev",
					Usage: "Start a development server.",
					Action: func(c *cli.Context) error {
						// 打印一下控制台日志:
						printTitle()
						// 打印控制台
						printRunning()
						// 把当前pid写入环境变量中
						pid := os.Getpid()
						os.Setenv(EnvVariable, strconv.Itoa(pid))
						// 平滑开启一个测试服务器,正常开启
						gracehttp.Serve(&http.Server{Addr: ":" + strconv.Itoa(settings.HttpPort), Handler: Singleton.Router})
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

// 在控制台打印 bingo logo
func printTitle() {
	fmt.Printf("\n %c[0;48;32m%s%c[0m\n\n", 0x1B, title, 0x1B)
}

// 在控制台打印编译日志
func printRunning() {
	fmt.Printf("\n %c[0;48;32m%s%c[0m\n", 0x1B, " Started successfully.", 0x1B)
	fmt.Printf("  App running at:\n  - Local: %s (copied to clipboard)\n  - Network: %s\n", DevDomainUrl, NetworkUrl)
}

// 监听目录变化，监听到变化后，重启向进程发送重启信号
func watchDirAndStartServer() {

}
