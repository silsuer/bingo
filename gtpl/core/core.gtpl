package core

import (
	"github.com/urfave/cli"
	"github.com/silsuer/bingo-router"
	"github.com/kylelemons/go-gypsy/yaml"
	"${path}/config"
	"github.com/facebookgo/grace/gracehttp"
	"net/http"
	"${path}/utils"
	"${path}/console"
	"fmt"
)

// 测试服务器路径和内网服务器路径
var (
	DomainUrl  string
	NetworkUrl string
)

// 核心模块，创建Bingo
// 构建Bingo结构体
type Bingo struct {
	Cli     *cli.App
	Router  *bingo_router.Router
	envPath string
	envFile *yaml.File
	Env     *utils.Environment
}

// 传入env文件路径
func NewApp(configPath string) *Bingo {
	b := &Bingo{
		Cli:     cli.NewApp(),
		Router:  bingo_router.New(),
		envPath: configPath,
	}
	f, err := yaml.ReadFile(configPath)
	if err != nil {
		panic(err)
	}
	b.envFile = f
	b.Env = utils.GetInstance(b.envPath)
	return b.Init()
}

func (app *Bingo) Init() *Bingo {
	// 初始化结构体
	// 拼接路径
	DomainUrl = app.Env.GetWithDefault("DOMAIN", "http://localhost") + ":" + app.Env.GetWithDefault("HTTP_PORT", "8080")
	network, _ := utils.IntranetIP()
	if len(network) > 0 {
		NetworkUrl = "http://" + network[0] + ":" + app.Env.GetWithDefault("HTTP_PORT", "8080")
	} else {
		NetworkUrl = DomainUrl
	}

	// 挂载基础命令
	app.Cli.Name = "bingo"
	app.Cli.Usage = "Bingo Cli"
	app.Cli.Version = config.VERSION
	app.Cli.Action = func(c *cli.Context) error {
		// 执行命令
		return nil
	}
	app.Cli.Commands = []cli.Command{
		{
			Name: "run",
			Action: func(c *cli.Context) error {
				printRunning()
				return gracehttp.Serve(&http.Server{Addr: ":" + app.Env.GetWithDefault("HTTP_PORT", "8080"), Handler: app.Router})
			},
		},
		{
			Name: "sword", // 添加sword
			Action: func(c *cli.Context) error {
				return nil
			},
			Subcommands: console.Schedule(),
		},
	}

	//console.Schedule()
	return app
}


func printRunning() {
	fmt.Printf("\n %c[0;48;32m%s%c[0m\n", 0x1B, " Started successfully.", 0x1B)
	fmt.Printf("  App running at:\n  - Local: %s (copied to clipboard)\n  - Network: %s\n", DomainUrl, NetworkUrl)
}
