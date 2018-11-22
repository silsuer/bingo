package main

import (
	"github.com/urfave/cli"
	"fmt"
	"os"
	"github.com/fsnotify/fsnotify"
	"path/filepath"
	"log"
	"github.com/xcltapestry/xclpkg/clcolor"
	"time"
	"os/exec"
	"sync"
	"bufio"
	"io"
	"strings"
)

const title = `
 ____    ___   _   _    ____    ___    _
| __ )  |_ _| | \ | |  / ___|  / _ \  | |
|  _ \   | |  |  \| | | |  _  | | | | | |
| |_) |  | |  | |\  | | |_| | | |_| | |_|
|____/  |___| |_| \_|  \____|  \___/  (_)
`

var command *exec.Cmd

var a sync.WaitGroup
// 在控制台打印 bingo logo
func printTitle() {
	fmt.Printf("\n %c[0;48;32m%s%c[0m\n\n", 0x1B, title, 0x1B)
}

// 开启一个dev服务器
// 并监听当前所有目录，发现有变化，则重启服务
func Dev(c *cli.Context) error {
	printTitle()
	// 执行make dev
	// 并监听目录
	// make dev
	// 定义一个函数，用来执行make dev
	// 定义一个函数，用来监听当前目录

	makeDev()
	a.Wait()
	return nil
}

func makeDev() {
	go watchDir()
	cmd := exec.Command("make", "dev")
	cmd.Stdout = os.Stdout // 控制台输出命令
	cmd.Stderr = os.Stdout // 如果有错误，也使用控制台进行输出
	command = cmd          // 赋值给全局变量
	if err := cmd.Start(); err != nil {
		fmt.Printf("The command `make dev` was wrong: %s", err)
	}

	if err := cmd.Wait(); err != nil {
		fmt.Printf("The command `make dev` was wrong: %s", err)
	}

}

// 监听当前目录下的数据，有变化就重启，并过滤掉.bingoignore中的数据
var watcher *fsnotify.Watcher

//监控目录
func watchDir() {
	a.Add(1)
	//先判断是否存在 .bingoignore
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		return
	}

	w, err := fsnotify.NewWatcher()
	watcher = w
	if err != nil {
		fmt.Println(err)
		return
	}

	// 默认监控下面所有文件
	err = watcher.Add("./")
	if err != nil {
		log.Fatal(err)
	}

	ignorePath := dir + string(os.PathSeparator) + ".bingoignore"

	//通过Walk来遍历目录下的所有子目录
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		//这里判断是否为目录，只需监控目录即可
		//目录下的文件也在监控范围内，不需要我们一个一个加
		if info.IsDir() {
			path, err := filepath.Abs(path)
			if err != nil {
				return err
			}
			err = watcher.Add(path)
			if err != nil {
				return err
			}
		}
		return nil
	})

	_, err = os.Stat(ignorePath)
	if err != nil {
		if !os.IsNotExist(err) {
			fmt.Println(err)
			return
		}
	}

	// 存在，则读取，并从watcher中移出这些文件或目录
	file, err := os.Open(ignorePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	r := bufio.NewReader(file)
	for {
		line, err := r.ReadString('\n') //以'\n'为结束符读入一行

		if err != nil || io.EOF == err {
			break
		}
		// 根据读入的数据进行移除
		p, _ := os.Getwd()
		// 拼接路径，传入一个方法中
		removeWatcher(p + "/" + strings.TrimSpace(line))
	}

	go func() {
		for {
			select {
			case ev := <-watcher.Events:
				{
					if ev.Op&fsnotify.Create == fsnotify.Create {
						fmt.Println(clcolor.Green("Created: " + ev.Name))
						timer()
						//这里获取新创建文件的信息，如果是目录，则加入监控中
						fi, err := os.Stat(ev.Name)
						if err == nil && fi.IsDir() {
							watcher.Add(ev.Name)
						}
					}
					if ev.Op&fsnotify.Write == fsnotify.Write {
						timer()
						fmt.Println(clcolor.Cyan("Writing: " + ev.Name))
					}
					if ev.Op&fsnotify.Remove == fsnotify.Remove {
						timer()
						fmt.Println(clcolor.Red("Deleted: " + ev.Name))
						//如果删除文件是目录，则移除监控
						fi, err := os.Stat(ev.Name)
						if err == nil && fi.IsDir() {
							watcher.Remove(ev.Name)
						}
					}
					if ev.Op&fsnotify.Rename == fsnotify.Rename {
						timer()
						fmt.Println(clcolor.Blue("Renamed: " + ev.Name))
						//如果重命名文件是目录，则移除监控
						//注意这里无法使用os.Stat来判断是否是目录了
						//因为重命名后，go已经无法找到原文件来获取信息了
						//所以这里就简单粗爆的直接remove好了
						watcher.Remove(ev.Name)
					}
					if ev.Op&fsnotify.Chmod == fsnotify.Chmod {
						// 修改权限
						//timer()
						//fmt.Printf(clcolor.Magenta("\nModifyPerm: " + ev.Name))
					}
				}
			case err := <-watcher.Errors:
				{
					fmt.Println("error : ", err)
					return
				}
			}
		}
	}()

	select {}
}

// 一个全局定时器
var t *time.Timer

func timer() {
	// 停止当前定时器，并重新定义定时器
	if t == nil {
		t = time.AfterFunc(1*time.Second, func() {
			// 重启
			// 删除
			// 向程序发送退出信号，并重新执行make dev
			if command != nil {
				command.Process.Kill()
				makeDev()
			}
		})
	} else {
		t.Reset(1 * time.Second)
	}

}

// 从监控器中移除目录
func removeWatcher(dir string) {
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			//removeWatcher(path)
			p, _ := filepath.Abs(path)
			watcher.Remove(p)
		} else {
			watcher.Remove(path)
		}
		return nil
	})
}
