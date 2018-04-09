package console

import (
	"flag"
	"os"
	"bingo/core"
	"strconv"
	"fmt"
)

type CLI struct{}

func (cli *CLI) Run() {
	//cli.validateArgs()
	initCmd := flag.NewFlagSet("init", flag.ExitOnError)
	runCmd := flag.NewFlagSet("run",flag.ExitOnError)
	runPort:= runCmd.Int("port",8088,"监听端口")
	switch os.Args[1] {
	case "init":
		err := initCmd.Parse(os.Args[2:]) // 把其余参数传入进去
		Check(err)
		break
	case "run":
		err := runCmd.Parse(os.Args[2:])
		Check(err)
		break
	default:
		break
	}

	if initCmd.Parsed() {
		// 如果是init命令
       cli.InitProject()
	}

	if runCmd.Parsed() {
		cli.RunSite(runPort)
	}

}

func (cli *CLI) InitProject()  {
   // 在当前目录下创建项目 目录
   // 包括  env.yaml app databases

   // 创建env.yaml
   // 创建

}

func (cli *CLI) RunSite(port *int)  {
  // 运行网站 实例化一个
  bingo := new(core.Bingo)
  fmt.Println("Bingo Running......")
  bingo.Run(":"+ strconv.Itoa(*port))
}

func Check(err error) {
	if err != nil {
		panic(err)
	}
}
