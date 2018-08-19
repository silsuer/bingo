package cli

import (
	"os"
	"github.com/silsuer/bingo/bingo"
	"strings"
	"fmt"
)

type MakeMiddleware struct {
	Command
	Name        string
	Description string
	Args        map[string]string
}

// 设置命令名
func (c *MakeMiddleware) SetName() {
	c.Name = "make:middleware"
}

// 设置命令所需参数
func (c *MakeMiddleware) SetArgs() {
	c.Args = make(map[string]string)
	c.Args["name"] = "" // 前面是参数名 后面是默认值
}

// 设置命令描述
func (c *MakeMiddleware) SetDescription() {
	c.Description = "create a middleware template"
}

// 设置命令实现的方法
func (c *MakeMiddleware) Handle(input Input, output Output) {
	name := input.Args["name"]
	currentPath, _ := os.Getwd()
	// 生成命令这个命令的路径
	middlewarePath := currentPath + "/" + bingo.Env.Get("MIDDLEWARES_PATH") + "/" + name + ".go"

	fmt.Println(middlewarePath)
	// 获取文件名
	arr := strings.Split(name, "/")

	middlewareName := arr[len(arr)-1] // 最后一个元素就是控制器名

	content := c.getContent(middlewareName)
	// 生成文件
	res := bingo.MakeFile(middlewarePath, content)
	if res {
		output.Info("the middleware created successfully!  Path:" + middlewarePath)
	} else {
		output.Error("the middleware created failed! Path:" + middlewarePath)
	}
}

func (c *MakeMiddleware) getContent(name string) string {
	content := `package Middlewares

import (
	"github.com/silsuer/bingo/bingo"
)

// 中间件结构体,可以重写Init方法设置是否是同步中间件
type ` + name + ` struct {
	bingo.Middleware
}

func (e *` + name + `) Handle(c *bingo.Context) *bingo.Context {
	return c
}
`
	return content
}
