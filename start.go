package main

import (
   "bingo/core/console"
)



func main() {
   //fmt.Println("由于需要加载配置文件，所以请务必在项目根目录下执行该命令,即go run start.go...")
   //fmt.Println("否则会导致路径错误！")
   //bingo := new(core.Bingo)
   //bingo.Run(":12345")
   cli := console.CLI{}
   cli.Run()
}