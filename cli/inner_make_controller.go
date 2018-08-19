package cli

import (
	"os"
	"github.com/silsuer/bingo/bingo"
	"strings"
)

type MakeController struct {
	Command
	Name        string
	Description string
	Args        map[string]string
}

// 设置命令名
func (m *MakeController) SetName() {
	m.Name = "make:controller"
}

// 设置命令所需参数
func (m *MakeController) SetArgs() {
	m.Args = make(map[string]string)
	m.Args["name"] = "" // 前面是参数名 后面是默认值
}

// 设置命令描述
func (m *MakeController) SetDescription() {
	m.Description = "创建控制器"
}

// 设置命令实现的方法
func (m *MakeController) Handle(input Input, output Output) {
	name := input.Args["name"]
	currentPath, _ := os.Getwd()
	// 生成命令这个命令的路径
	controllerPath := currentPath + "/" + bingo.Env.Get("CONTROLLERS_PATH") + "/" + name + ".go"

	// 获取文件名
	arr := strings.Split(name, "/")

	controllerName := arr[len(arr)-1] // 最后一个元素就是控制器名

	content := m.getContent(controllerName)
	// 生成文件
	res := bingo.MakeFile(controllerPath, content)
	if res {
		output.Info("the controller created successfully!  Path:" + controllerPath)
	} else {
		output.Error("the controller created failed! Path:" + controllerPath)
	}

}

func (m *MakeController) getContent(controllerName string) string {
	content := `package Controllers

import (
	"github.com/silsuer/bingo/bingo"
)

// 控制器
type ` + controllerName + ` struct {
	bingo.Controller
}

`
	return content
}
