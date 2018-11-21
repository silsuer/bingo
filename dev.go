package main

import (
	"github.com/urfave/cli"
	"os/exec"
	"fmt"
)

// 开启一个dev服务器
// 并监听当前所有目录，发现有变化，则重启服务
func Dev(c *cli.Context) error {
	// 执行make dev
	// 并监听目录
	// make dev
    fmt.Println("111")
	// 定义一个函数，用来执行make dev
	// 定义一个函数，用来监听当前目录
    makeDev()
	return nil
}

func makeDev() {
	cmd := exec.Command("make", "dev")
	cmd.Run()  // 后台运行
	cmd.Wait()

}
