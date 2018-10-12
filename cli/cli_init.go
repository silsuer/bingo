package cli

import (
	"os"
	"github.com/silsuer/bingo/bingo"
	"io/ioutil"
	"fmt"
	"strings"
)

const envContent = `# Bingo Config File .......
# 静态文件夹路径
STATIC_FILE_DIR : public
# session驱动配置 map是以map的形式存储在内存中 db 存储在数据库中 redis存储在redis中 bolt，存储在bolt中
SESSION_DRIVER : kvstorage
# 用数据库存，会生成sessions表，用KVStorage存储，会生成sessions bucket
SESSION_DRIVER_NAME : sessions

# 配置文件夹根目录
CONFIG_ROOT : config
# 脚手架命令根目录
CONSOLE_KERNEL_PATH : app/Console

# Http Kernel目录
HTTP_KERNEL_PATH : app/Http

# 控制器目录
CONTROLLERS_PATH : app/Http/Controllers
# 中间件目录
MIDDLEWARES_PATH : app/Http/Middlewares

# 数据库配置
DB_DRIVER : MYSQL
DB_HOST : localhost
DB_NAME : bingo
DB_PORT : 3306
DB_USERNAME : root
DB_PASSWORD : root
DB_CHARSET: utf8

# 程序平滑重启的时候记录pid的文件
PID_FILE: start.pid



# bolt设置
# 数据库文件名
KVSTORAGE_DB_NAME : bingo.db
# kv存储时的默认bucket的名字
KVSTORAGE_BUCKET : bingo
              `

const utilsContent = `package utils

import "github.com/silsuer/bingo/bingo"

func Route() bingo.Route {
	return bingo.Route{}
}
`

const consoleKernelContent = `package main

import (
	Command "test/app/Console/Commands"
	"os"
	"github.com/silsuer/bingo/cli"
)

var Commands = []interface{}{
	&Command.ExampleCommand{},
}

func main() {
	console := cli.Console{}
	console.Exec(os.Args, Commands)
}

// 调度任务
func schedule() {

}
`

const consoleExampleContent = `package Commands

import (
	"github.com/silsuer/bingo/cli"
)

type ExampleCommand struct {
	cli.Command
	Name        string
	Description string
	Args        map[string]string
}

// 设置命令名
func (m *ExampleCommand) SetName() {
	m.Name = "command:name"
}

// 设置命令所需参数
func (m *ExampleCommand) SetArgs() {
	m.Args = make(map[string]string)
	m.Args["name"] = ""
}

// 设置命令描述
func (m *ExampleCommand) SetDescription() {
	m.Description = "the command description."
}

// 设置命令实现的方法
func (m *ExampleCommand) Handle(input cli.Input, output cli.Output) {

}
`

func (cli *CLI) InitProject() {
	// 在当前目录下创建项目 目录
	// 包括  env.yaml glide.yaml app databases console// app下放置一个路由，下一个基本的路由文件
	// 创建env.yaml
	makeEnvFile()
	//  创建start文件
	makeStartFile()

	// 创建utils文件
	makeUtilsFile()

	//创建路由文件
	makeRoutesFile()

	// 创建console文件夹
	makeConsoleFinder()

	// 创建Http文件夹
	makeHttpFinder()

	// 输出结果
	fmt.Println("Your Bingo Project Init Successfully!")

	asciiContent := `
 ____    ___   _   _    ____    ___    _
| __ )  |_ _| | \ | |  / ___|  / _ \  | |
|  _ \   | |  |  \| | | |  _  | | | | | |
| |_) |  | |  | |\  | | |_| | | |_| | |_|
|____/  |___| |_| \_|  \____|  \___/  (_)
`
	fmt.Printf("\n %c[0;48;32m%s%c[0m\n\n", 0x1B, asciiContent, 0x1B)
}

func makeHttpFinder() {
	// 创建Http/Kernel.go
	dir, err := os.Getwd()
	bingo.Check(err)
	kernelPath := dir + "/" + bingo.Env.Get("HTTP_KERNEL_PATH") + "/Kernel.go"

	// 写入数据
	content := `package Http

// 注册全局中间件 中间件组 路由中间件
func init() {
	//bingo.GlobalMiddlewareRegister(globalM)
	//bingo.GroupMiddlewareRegister(groupM)
	//bingo.RouteMiddlewareRegister(routeM)
}
//
//// 全局中间件
//var globalM = []bingo.MiddlewareInterface{
//	//&Middlewares.ExampleMiddleware{},
//}
//
//// 中间件组 组名->中间件组
//var groupM = map[string][]bingo.MiddlewareInterface{
//	"auth": {
//		//&Middlewares.ExampleMiddleware{},
//	},
//}
//
//// 路由中间件(给中间件命名)  名称->中间件
//var routeM = map[string]bingo.MiddlewareInterface{
//	//"m.name": &Middlewares.ExampleMiddleware{},
//}

`
	if bingo.FileIsExist(kernelPath) {
		fmt.Println("the http kernel.go file has alreay exist.")
	} else {
		res := bingo.MakeFile(kernelPath, content)
		if res {
			fmt.Println("the http kernel.go file created successfully!")
		} else {
			fmt.Println("the http kernel.go file created error!")
		}
	}

}

// 创建env文件
func makeEnvFile() {
	dir, err := os.Getwd()
	bingo.Check(err)
	//fmt.Println(dir)
	// 创建env
	envPath := dir + "/env.yaml"
	if bingo.CheckFileIsExist(envPath) {
		// 存在，抛出错误
		ThrowError("The env file has already exists")
	} else {
		// 不存在，创建
		err = ioutil.WriteFile(envPath, []byte(envContent), 0666)
		fmt.Println("env file created successfully!")
		bingo.Check(err)
	}
}

// 创建start文件
func makeStartFile() {

	// 获取当前目录
	currentPath := bingo.GetCurrentDirectory()
	currentDir := strings.Split(currentPath, "/")

	startContent := `package main

import (
	"github.com/silsuer/bingo/bingo"
	"` + currentDir[len(currentDir)-1] + `/routes"
	_ "` + currentDir[len(currentDir)-1] + `/app/Http"

)

func main() {
	b := bingo.Bingo{}
	routes.SetRoutes()
	b.Run(":12345")
}


`

	dir, err := os.Getwd()
	bingo.Check(err)
	startPath := dir + "/start.go"
	if bingo.CheckFileIsExist(startPath) {
		ThrowError("The start file has already exists")
	} else {
		err = ioutil.WriteFile(startPath, []byte(startContent), 0666)
		fmt.Println("start.go file created successfully!")
		bingo.Check(err)
	}
}

// 创建utils文件
func makeUtilsFile() {
	dir, err := os.Getwd()
	bingo.Check(err)
	utilsPath := dir + "/utils/utils.go"

	if bingo.MakeFile(utilsPath, utilsContent) {
		fmt.Println("utils.go file created successfully!")
	}

}

// 创建路由文件
func makeRoutesFile() {

	// 获取当前目录
	//currentPath := bingo.GetCurrentDirectory()
	//currentDir := strings.Split(currentPath, "/")

	routeContent := `package routes

import (
	"github.com/silsuer/bingo/bingo"
	"fmt"
)

func SetRoutes()  {
	
	bingo.NewRoute().Get("/").Target(func(c *bingo.Context) {
		fmt.Fprint(c.Writer, "Hello,Bingo!")
	}).Register()
}`

	dir, err := os.Getwd()
	bingo.Check(err)

	routePath := dir + "/routes/web.go"

	if bingo.MakeFile(routePath, routeContent) {
		fmt.Println("route file created successfully!")
	}
}

func makeConsoleFinder() {
	dir, err := os.Getwd()
	bingo.Check(err)
	// 创建kernel
	kernelPath := dir + "/app/Console/Kernel.go"
	// 创建example文件
	examplePath := dir + "/app/Console/Commands/ExampleCommand.go"

	if bingo.MakeFile(kernelPath, consoleKernelContent) && bingo.MakeFile(examplePath, consoleExampleContent) {
		fmt.Println("the console finder created successfully!")
	}
}
