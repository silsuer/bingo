package cli

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
	m.Description = "这是哥测试命令"
}

// 设置命令实现的方法
func (m *MakeController) Handle(input Input, output Output) {

}
