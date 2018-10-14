package bingo

import (
	"net/http"
	"github.com/gorilla/context"
	"fmt"
	"os"
	"github.com/fsnotify/fsnotify"
	"path/filepath"
	"time"
	"sync"
	"syscall"
	"strings"
)

// 重启队列
var restartSlice []int
var fileWatcher *fsnotify.Watcher
var wdDir string

var BasePath string

func init() {
	// 默认的BasePath是运行命令时所在的路径
	BasePath, _ = os.Getwd()
}

// bingo结构体，向外暴露一些属性和方法  实现了http方法
type Bingo struct{}

// 开启的端口号，传入参数，是否允许使用os.Args
func (b *Bingo) Run(port string) {

	args := []string{"daemon", "start"}

	if len(os.Args) != 0 {
		args = os.Args[1:]
	}
	b.setGlobalParamFromArgs(args)
	// 根据httprouter进行重写(根据Httprouter的原理，重新实现路由)
	// 这个时候要根据RouteList,对每一个方法解析出一个tree来
	router := New()
	// 开始把路由列表注册到tree中
	for _, v := range RouteList {
		router.Handle(v.Method, v.Path, v)
	}
	// 静态页面
	router.NotFound = func(writer http.ResponseWriter, request *http.Request) {
		http.ServeFile(writer, request, GetPublicPath()+request.URL.Path)
	}

	// 根据传入的参数，选择不同的开启服务器的方式，现在有几种：
	// watch 监听目录，一旦有更新，则重启服务器
	// daemon 以守护进程运行
	// 开启服务器
	b.startServer(args, port, context.ClearHandler(router))
}

// 从传入的参数中提取出根目录等参数并赋值
func (b *Bingo) setGlobalParamFromArgs(args []string) {
	for _, arg := range args {
		// a b path=/home/work/jx
		if strings.Contains(arg, "=") {
			p := strings.Split(arg, "=")
			if len(p) != 2 {
				continue
			}
			if p[0] == "path" {
				BasePath = p[1]
			}
		}
	}
}

func (b *Bingo) startServer(params []string, port string, handler http.Handler) {
	param := "dev"
	if len(params) > 0 {
		param = params[0]
	}

	switch param {
	case "dev":
		startDevServer(port, handler)
	case "watch":
		startWatchServer(port, handler)
	case "daemon":
		startDaemonServer(port, handler, params[1:])
	default:
		startDevServer(port, handler)
		//fmt.Println("undefined param:" + param)
	}
}

func startDevServer(port string, handler http.Handler) {
	// 平滑启动服务
	GracefulServe(port, handler)
}

// 监听目录变化，如果有变化，重启服务
// 守护进程开启服务，主进程阻塞不断扫描当前目录，有任何更新，向守护进程传递信号，守护进程重启服务
// 开启一个协程运行服务
// 监听目录变化，有变化运行 bingo run daemon restart
func startWatchServer(port string, handler http.Handler) {
	f, err := fsnotify.NewWatcher()
	if err != nil {
		panic(err)
	}
	defer f.Close()
	fileWatcher = f
	f.Add(BasePath)

	done := make(chan bool)

	go func() {
		procAttr := &syscall.ProcAttr{
			Env:   os.Environ(),
			Files: []uintptr{os.Stdin.Fd(), os.Stdout.Fd(), os.Stderr.Fd()},
		}
		_, err := syscall.ForkExec(os.Args[0], []string{os.Args[0], "daemon", "start"}, procAttr)
		if err != nil {
			fmt.Println(err)
		}
	}()

	go func() {
		for {
			select {
			case ev := <-f.Events:
				if ev.Op&fsnotify.Create == fsnotify.Create {
					fmt.Printf("\n %c[0;48;33m%s%c[0m", 0x1B, "["+time.Now().Format("2006-01-02 15:04:05")+"]created file:"+ev.Name, 0x1B)
				}
				if ev.Op&fsnotify.Remove == fsnotify.Remove {
					fmt.Printf("\n %c[0;48;31m%s%c[0m", 0x1B, "["+time.Now().Format("2006-01-02 15:04:05")+"]deleted file:"+ev.Name, 0x1B)
				}
				if ev.Op&fsnotify.Rename == fsnotify.Rename {
					fmt.Printf("\n %c[0;48;34m%s%c[0m", 0x1B, "["+time.Now().Format("2006-01-02 15:04:05")+"]renamed file:"+ev.Name, 0x1B)
				} else {
					fmt.Printf("\n %c[0;48;32m%s%c[0m", 0x1B, "["+time.Now().Format("2006-01-02 15:04:05")+"]modified file:"+ev.Name, 0x1B)
				}
				// 有变化，放入重启数组中
				restartSlice = append(restartSlice, 1)
			case err := <-f.Errors:
				fmt.Println("error:", err)
			}
		}
	}()

	// 准备重启守护进程
	go restartDaemonServer()

	<-done
}

// 重启守护进程
func restartDaemonServer() {
	listeningWatcherDir(wdDir)
	var mutex sync.Mutex
	for {
		// 如果重启切片中有数据，证明数据有变动，重新设置监听目录
		if len(restartSlice) > 0 {
			mutex.Lock()

			go func() {
				DaemonInit("restart")
			}()
			listeningWatcherDir(wdDir)
			restartSlice = []int{}
			mutex.Unlock()
		}
		// 睡1s
		time.Sleep(time.Second)
	}
}

func listeningWatcherDir(dir string) {
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		//dir, _ := os.Getwd()
		pidFile = BasePath + "/" + Env.Get("PID_FILE")

		fmt.Println(pidFile)
		fmt.Println(path)
		fileWatcher.Add(path)
		fileWatcher.Remove(pidFile)
		return nil
	})
}

// 以守护进程运行
func startDaemonServer(port string, handler http.Handler, args []string) {
	if len(args) > 0 {
		DaemonInit(args[0])
	} else {
		DaemonInit("start")
	}
}
