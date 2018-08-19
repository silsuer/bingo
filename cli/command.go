package cli

import "strings"
import (
	"reflect"
	"fmt"
	"time"
)

//import "fmt"

// 控制台的kernel都继承这个对象，将自动筛选指定的命令
type Console struct{}

var InnerCommands = []interface{}{
	&MakeController{},
	&MakeCommand{},
	&MakeOriginCommand{},
	&MakeMiddleware{},
}

// 命令运行
// 传入命令行参数
// 第一个参数是命令名，后面跟着的是参数名
// 接着传入命令组
// 首先拆出命令名
// 然后遍历命令数组，断言是否实现了Command接口，然后根据参数，构建input对象和output对象
// 调用命令的handle方法，传入输入输出对象
func (console *Console) Exec(args []string, commands []interface{}) {

	input := console.initInput(args)
	//fmt.Println(args)
	for _, command := range commands {

		// 先做检查,查找对应的命令名
		commandValue := reflect.ValueOf(command)
		// 初始化命令结构体
		initCommand(&commandValue)

		// 映射期望参数与实际输入参数（验证参数输入是否正确）
		target := checkParams(command, &input)
		// 不是这个命令，跳过这个命令
		if target == false {
			continue
		}
		// 获得输入和输出并准备作为参数传入Handle方法中
		var params = []reflect.Value{reflect.ValueOf(input), reflect.ValueOf(Output{})}
		commandValue.MethodByName("Handle").Call(params)
	}
}

func checkParams(command interface{}, input *Input) bool {
	// 检查是否有输入问题
	// 先检查是否是这个命令
	if reflect.ValueOf(command).Elem().FieldByName("Name").Interface().(string) != input.Name {
		return false
	}

	// 如果需要两个参数 aaa="" bbb="",只输入了一个bbb，那么要将aaa的默认参数传进去
	// 判断是否有参数
	//fmt.Println(reflect.ValueOf(command).NumField())
	args := reflect.ValueOf(command).Elem().FieldByName("Args").Interface().(map[string]string)

	// 输入的值不需要管，只需要查那些需要输入但是没有在输入中的值
	for arg, value := range args {
		if _, ok := input.Args[arg]; !ok {
			input.Args[arg] = value
		}
	}
	return true
}

func initCommand(commandValue *reflect.Value) {
	var initParam = []reflect.Value{}
	commandValue.MethodByName("SetName").Call(initParam)
	commandValue.MethodByName("SetArgs").Call(initParam)
	commandValue.MethodByName("SetDescription").Call(initParam)
}

type Input struct {
	Name string            // 命令名
	Args map[string]string // 参数名->参数值
}

type Output struct{}

// 在控制台中输出普通信息
func (o *Output) Info(content string) {
	fmt.Printf("\n %c[0;48;32m%s%c[0m\n\n", 0x1B, "["+time.Now().Format("2006-01-02 15:04:05")+"]"+content, 0x1B)
}

// 在控制台中输出错误信息
func (o *Output) Error(content string) {
	fmt.Printf("\n %c[0;48;31m%s%c[0m\n\n", 0x1B, "["+time.Now().Format("2006-01-02 15:04:05")+"]"+content, 0x1B)
}

func (console *Console) initInput(args []string) Input {
	// 第一组就是命令名
	// 之后的分别拆分，存入map
	var argMap = make(map[string]string)
	name := "default"
	if len(args) > 1 {
		name = args[1:2][0] // 得到命令名
		for _, arg := range args[2:] {
			//以=分割
			a := strings.Split(arg, "=")
			getArg(&argMap, a) // 获取参数
		}
	}
	return Input{Name: name, Args: argMap}
}

func getArg(argMap *map[string]string, arg []string) {
	argName := strings.Trim(arg[0], "--")
	argValue := arg[1]
	(*argMap)[argName] = argValue
}

// 命令原型
type OriginCommand interface {
	SetName()
	SetDescription()
	SetArgs()
	Handle(input Input, output Output)
}

// 所有命令都藏了一个这个结构体
type Command struct {
	Name        string
	Description string
	Args        map[string]string
}

// 命令执行方法
func (c *Command) Handle(input Input, output Output) {

}

// 设置命令名，这里将设置默认命令名
func (c *Command) SetName() {

}

// 设置描述
func (c *Command) SetDescription() {

}

func (c *Command) SetArgs() {

}
