package main

import (
   "fmt"
   "bingo/core"
)



func main() {
 // new 一个bingo对象，然后bingo.Run()即可
 // 指定静态目录，默认指向index.html ，加载路由文件
 // 加载env文件
   fmt.Println("由于需要加载配置文件，所以请务必在项目根目录下执行该命令,即go run start.go...")
   fmt.Println("否则会导致路径错误！")
   bingo := new(core.Bingo)
   bingo.Run(":12345")
}