package cli

import (
	"fmt"
	"github.com/silsuer/bingo/bingo"
	"os"
	"os/exec"
	"bytes"
	"log"
)

// 专门处理sword方法

// 处理sword命令
func (cli *CLI) swordHandle(args []string) {
	// 解析这个参数，将数据传入外部
	//fmt.Println(args)
	//获取env中的kernel路径
	//根据kernel去 go shell执行 go run xxx/kernel.go make:controller AdminController
	consoleKernelPath := bingo.Env.Get("CONSOLE_KERNEL_PATH")
	// 获取当前目录
	dir, _ := os.Getwd()

	consoleKernelAbsolutePath := dir + "/" + consoleKernelPath + "/Kernel.go"

	// 使用go shell 调用 go run xxx/Kernel.go arg1 arg2 arg3
	var tmpSlice = []string{"go", "run", consoleKernelAbsolutePath}

	args = append(tmpSlice, args...)

	// 先检查这个命令是否属于内部命令
	// arg第一个就是命令
	console := Console{}
	console.Exec(args[2:], InnerCommands)

	//[run /Users/silsuer/go/src/test/app/Console/Kernel.go aaa bbb ccc]
	//cmd := exec.Command("sh", "-c", strings.Join(args, " "))
	cmd := exec.Command("go",args[1:]...)
	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Start(); err != nil {
		panic(err)
	}

	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}

	fmt.Print(out.String())
}
