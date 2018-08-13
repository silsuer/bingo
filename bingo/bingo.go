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
)

// 重启队列
var restartSlice []int
var fileWatcher *fsnotify.Watcher
var wdDir string

// bingo结构体，向外暴露一些属性和方法  实现了http方法
type Bingo struct{}

func (b *Bingo) Run(port string, args []string) {
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

	// TODO 监听平滑升级和重启
}

func (b *Bingo) startServer(params []string, port string, handler http.Handler) {
	param := "dev"
	if len(params) > 0 {
		param = params[0]
	}

	switch param {
	case "dev":
		//startDevServer(port, handler)
		startDevServer(port, handler)
	case "watch":
		startWatchServer(port, handler)
	case "daemon":
		startDaemonServer(port, handler)
	default:
		fmt.Println("undefined param:" + param)
	}
}

func startDevServer(port string, handler http.Handler) {
	// 平滑启动服务
	GracefulServe(port, handler)
}

func startWatchServer(port string, handler http.Handler) {
	// 监听目录变化，如果有变化，重启服务
	// 守护进程开启服务，主进程阻塞不断扫描当前目录，有任何更新，向守护进程传递信号，守护进程重启服务
	// 开启一个协程运行服务
	// 监听目录变化，有变化运行 bingo run daemon restart
	f, err := fsnotify.NewWatcher()
	if err != nil {
		panic(err)
	}
	defer f.Close()
	dir, _ := os.Getwd()
	wdDir = dir
	fileWatcher = f
	f.Add(dir)

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
					fmt.Println("created file : ", ev.Name)
				}
				if ev.Op&fsnotify.Remove == fsnotify.Remove {
					fmt.Println("deleted file : ", ev.Name)
				}
				if ev.Op&fsnotify.Rename == fsnotify.Rename {
					fmt.Println("rename file : ", ev.Name)
				} else {
					fmt.Println("modified file : ", ev.Name)
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
	procAttr := &syscall.ProcAttr{
		Env:   os.Environ(),
		Files: []uintptr{os.Stdin.Fd(), os.Stdout.Fd(), os.Stderr.Fd()},
	}
	listeningWatcherDir(wdDir)
	var mutex sync.Mutex
	for {
		// 如果重启切片中有数据，证明数据有变动，重新设置监听目录
		if len(restartSlice) > 0 {
			mutex.Lock()
			_, err := syscall.ForkExec(os.Args[0], []string{os.Args[0], "daemon", "restart"}, procAttr)
			if err != nil {
				fmt.Println(err)
			}

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
		dir, _ := os.Getwd()
		pidFile = dir + "/" + Env.Get("PID_FILE")
		fileWatcher.Add(path)
		fileWatcher.Remove(pidFile)
		return nil
	})
}

// 以守护进程运行
func startDaemonServer(port string, handler http.Handler) {
	DaemonInit()
}
