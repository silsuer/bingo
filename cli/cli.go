package cli

import (
	"flag"
	"os"
	"strconv"
	"fmt"
	"errors"
	"io/ioutil"
	"github.com/silsuer/bingo/bingo"
)

const envContent = `# Bingo Config File .......
# 静态文件夹路径
STATIC_FILE_DIR : public
# session驱动配置 map是以map的形式存储在内存中 db 存储在数据库中 redis存储在redis中 bolt，存储在bolt中
SESSION_DRIVER : kvstorage
# 用数据库存，会生成sessions表，用KVStorage存储，会生成sessions bucket
SESSION_DRIVER_NAME : sessions

# 数据库配置
DB_DRIVER : MYSQL
DB_HOST : localhost
DB_NAME : bingo
DB_PORT : 3306
DB_USERNAME : root
DB_PASSWORD : root
DB_CHARSET: utf8

# bolt设置
# 数据库文件名
KVSTORAGE_DB_NAME : bingo.db
# kv存储时的默认bucket的名字
KVSTORAGE_BUCKET : bingo
           
           `

// start.go中的数据
const startContent = `

package main

import (
	"github.com/silsuer/bingo/bingo"
	"fmt"
)

var Welcome = []bingo.Route{
	{
		Path:"/",
		Method:bingo.GET,
		Target: func(c *bingo.Context) {
		    fmt.Fprintln(c.Writer,"<h1>Welcome to Bingo!</h1>")	
		},
	},
}

func main() {
	b := bingo.Bingo{}
	bingo.RegisterRoute(Welcome)
	b.Run(":12345")
}

           `

type CLI struct{}

func (cli *CLI) Run() {
	//cli.validateArgs()
	initCmd := flag.NewFlagSet("init", flag.ExitOnError)
	runCmd := flag.NewFlagSet("run", flag.ExitOnError)
	runPort := runCmd.Int("port", 8088, "监听端口") // 默认是8088端口

	swordCmd := flag.NewFlagSet("sword", flag.ExitOnError) // bingo sword 命令
	//swordConfig := swordCmd.String("name","list","the Commands name")
	var swordConfig []string
	switch os.Args[1] {
	case "init":
		err := initCmd.Parse(os.Args[2:]) // 把其余参数传入进去
		bingo.Check(err)
		break
	case "run":
		err := runCmd.Parse(os.Args[2:])
		bingo.Check(err)
		break
	case "sword":
		err := swordCmd.Parse(os.Args[2:])
		swordConfig = os.Args[2:]
		bingo.Check(err)
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

	if swordCmd.Parsed() {
       cli.swordHandle(swordConfig)
	}
}


func (cli *CLI) InitProject() {
	// 在当前目录下创建项目 目录
	// 包括  env.yaml glide.yaml app databases // app下放置一个路由，下一个基本的路由文件
	// 创建env.yaml
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
		bingo.Check(err)
	}
	//  创建start文件
	startPath := dir + "/start.go"
	if bingo.CheckFileIsExist(startPath) {
		ThrowError("The start file has already exists")
	} else {
		err = ioutil.WriteFile(startPath, []byte(startContent), 0666)
		bingo.Check(err)
	}
	// 输出结果
	fmt.Println("Your Bingo Project Init Successfully!")
}

func (cli *CLI) RunSite(port *int) {
	// 运行网站 实例化一个
	b := new(bingo.Bingo)
	fmt.Println("Bingo Running......")
	b.Run(":" + strconv.Itoa(*port))
}

func ThrowError(text string) {
	panic(errors.New(text))
}
