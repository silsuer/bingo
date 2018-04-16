package bingo

import (
	"flag"
	"os"
	"strconv"
	"fmt"
	"errors"
	"io/ioutil"
)

const envContent = `# Bingo Config File .......
# 静态文件夹路径
STATIC_FILE_DIR : public

# 数据库配置
DB_DRIVER : MYSQL
DB_HOST : localhost
DB_NAME : bingo
DB_PORT : 3306
DB_USERNAME : root
DB_PASSWORD : root
DB_CHARSET: utf8
           `

// start.go中的数据
const startContent = `
package main

import (
	"bingo/bingo"
	"net/http"
	"fmt"
)

var Welcome = []bingo.Route{
	{
		Path:"/",
		Method:bingo.GET,
		Target: func(writer http.ResponseWriter, request *http.Request, params bingo.Params) {
			fmt.Fprint(writer,"<h1>Welcome to Bingo!</h1>")
		},
	},
}

func main() {
	b := bingo.Bingo{}
	bingo.RegistRoute(Welcome)
	b.Run(":12345")
}

           `

type CLI struct{}

func (cli *CLI) Run() {
	//cli.validateArgs()
	initCmd := flag.NewFlagSet("init", flag.ExitOnError)
	runCmd := flag.NewFlagSet("run", flag.ExitOnError)
	runPort := runCmd.Int("port", 8088, "监听端口")  // 默认是8088端口
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

func (cli *CLI) InitProject() {
	// 在当前目录下创建项目 目录
	// 包括  env.yaml glide.yaml app databases // app下放置一个路由，下一个基本的路由文件
	// 创建env.yaml
	dir, err := os.Getwd()
	Check(err)
	//fmt.Println(dir)
	// 创建env
	envPath := dir + "/env.yaml"
	if CheckFileIsExist(envPath) {
		// 存在，抛出错误
		ThrowError("The env file has already exists")
	} else {
		// 不存在，创建
		err = ioutil.WriteFile(envPath, []byte(envContent), 0666)
		Check(err)
	}
	//  创建start文件
	startPath:= dir+"/start.go"
	if CheckFileIsExist(startPath) {
		ThrowError("The start file has already exists")
	}else{
		err = ioutil.WriteFile(startPath,[]byte(startContent),0666)
		Check(err)
	}
	// 输出结果
	fmt.Println("Your Bingo Project Init Successfully!")
}

func (cli *CLI) RunSite(port *int) {
	// 运行网站 实例化一个
	bingo := new(Bingo)
	fmt.Println("Bingo Running......")
	bingo.Run(":" + strconv.Itoa(*port))
}

func ThrowError(text string) {
	panic(errors.New(text))
}
