package cli

import (
	"os"
	"github.com/silsuer/bingo/bingo"
)

type MakeOriginCommand struct {
	Command
	Name        string
	Description string
	Args        map[string]string
}

// 设置命令名
func (m *MakeOriginCommand) SetName() {
	m.Name = "make:origin:command"
}

// 设置命令所需参数
func (m *MakeOriginCommand) SetArgs() {
	m.Args = make(map[string]string)
	m.Args["name"] = ""
}

// 设置命令描述
func (m *MakeOriginCommand) SetDescription() {
	m.Description = "生成一个默认内置命令文件"
}

// 设置命令实现的方法
func (m *MakeOriginCommand) Handle(input Input, output Output) {
	name := input.Args["name"]
	goPath := os.Getenv("GOPATH")
	// 生成命令这个命令的路径
	commandsPath := goPath + "/src/github.com/silsuer/bingo/cli/" + name + ".go"
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
func (m *MakeOriginCommand) getContent(input Input) string {

	str := `package cli

type ` + input.Args["name"] + ` struct {
	Command
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
	c.Args["name"] = "" // 前面是参数名 后面是默认值
}

// 设置命令描述
func (c *` + input.Args["name"] + `) SetDescription() {
	c.Description = "the command description"
}

// 设置命令实现的方法
func (c *` + input.Args["name"] + `) Handle(input Input, output Output) {

}
`
	return str
}
