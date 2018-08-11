package cli

import (
	"os"
	"github.com/silsuer/bingo/bingo"
)

type MakeCommand struct {
	Command
	Name        string
	Description string
	Args        map[string]string
}

// 设置命令名
func (m *MakeCommand) SetName() {
	m.Name = "make:command"
}

// 设置命令所需参数
func (m *MakeCommand) SetArgs() {
	m.Args = make(map[string]string)
	m.Args["name"] = ""
}

// 设置命令描述
func (m *MakeCommand) SetDescription() {
	m.Description = "生成一个默认命令文件"
}

// 设置命令实现的方法
func (m *MakeCommand) Handle(input Input, output Output) {
	name := input.Args["name"]
	currentPath, _ := os.Getwd()
	// 生成命令这个命令的路径
	commandsPath := currentPath + "/" + bingo.Env.Get("CONSOLE_KERNEL_PATH") + "/Commands/" + name + ".go"
	// 建立文件 并写入信息
	if bingo.CheckFileIsExist(commandsPath) {
		output.Error("the file " + name + ".go has already exist :" + commandsPath)
	} else {
		if bingo.MakeFile(commandsPath, m.getContent(input)) {
			output.Info("Command make successfully :" + commandsPath)
		}
	}
}

// 获取默认文件的内容
func (m *MakeCommand) getContent(input Input) string {

	str := `package Commands

import (
	"github.com/silsuer/bingo/cli"
)

type ` + input.Args["name"] + ` struct {
	cli.Command
	Name        string
	Description string
	Args        map[string]string
}

// 设置命令名
func (c *` + input.Args["name"] + `) SetName() {
	c.Name = "command:name"
}

// 设置命令所需参数
func (c *` + input.Args["name"] + `) SetArgs() {
	c.Args = make(map[string]string)
	c.Args["name"] = ""
}

// 设置命令描述
func (c *` + input.Args["name"] + `) SetDescription() {
	c.Description = "the command description."
}

// 设置命令实现的方法
func (c *` + input.Args["name"] + `) Handle(input cli.Input, output cli.Output) {
	
}
`
	return str
}
