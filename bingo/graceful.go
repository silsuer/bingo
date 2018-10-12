package bingo

import (
	"net/http"
	"os"
	"net"
	"sync"
	"time"
	"syscall"
	"io/ioutil"
	"strconv"
	"os/signal"
	"log"
	"fmt"
)

var (
	TimeDeadLine = 10 * time.Second
	srv          *server
	appName      string
	pidFile      string
	pidVal       int
)

type ConnectionManager struct {
	sync.WaitGroup
	Counter   int
	mux       sync.Mutex
	idleConns map[string]net.Conn  // 面向网络的连接
}

//improvement http.Server
type server struct {
	http.Server
	listener *listener
	cm       *ConnectionManager
}

//用来重载net.Listener的方法
type listener struct {
	net.Listener
	server *server
}

func (cm *ConnectionManager) add(delta int) {
	cm.Counter += delta
	cm.WaitGroup.Add(delta)
}

func (cm *ConnectionManager) done() {
	cm.Counter--
	cm.WaitGroup.Done()
}

// 设置超时时间，超过这个时间，连接将关闭
func (cm *ConnectionManager) close(t time.Duration) {
	cm.mux.Lock()
	dt := time.Now().Add(t)
	for _, c := range cm.idleConns {
		c.SetDeadline(dt)
	}
	cm.idleConns = nil
	cm.mux.Unlock()
	cm.WaitGroup.Wait()
	return
}

func (cm *ConnectionManager) rmIdleConns(key string) {
	cm.mux.Lock()
	delete(cm.idleConns, key)
	cm.mux.Unlock()
}

func (cm *ConnectionManager) addIdleConns(key string, conn net.Conn) {
	cm.mux.Lock()
	cm.idleConns[key] = conn
	cm.mux.Unlock()
}

// 平滑重启与关闭进程
//使用addr和handler来启动一个支持graceful的服务
func GracefulServe(addr string, handler http.Handler) error {

	// 生成一个server对象
	s := http.Server{
		Addr:    addr,
		Handler: handler,
	}
	// 调用方法
	return Graceful(s)
}

func newConnectionManager() *ConnectionManager {
	cm := &ConnectionManager{}
	cm.WaitGroup = sync.WaitGroup{}
	cm.idleConns = make(map[string]net.Conn)
	return cm
}

//处理http.Server，使支持graceful stop/restart
func Graceful(s http.Server) error {
	// 设置一个环境变量
	os.Setenv("__GRACEFUL", "true")
	// 创建一个自定义的server
	srv = &server{
		cm:     newConnectionManager(),
		Server: s,
	}

	// 设置server的状态
	srv.ConnState = func(conn net.Conn, state http.ConnState) {
		switch state {
		case http.StateNew:
			srv.cm.add(1)
		case http.StateActive:
			srv.cm.rmIdleConns(conn.LocalAddr().String())
		case http.StateIdle:
			srv.cm.addIdleConns(conn.LocalAddr().String(), conn)
		case http.StateHijacked, http.StateClosed:
			srv.cm.done()
		}
	}
	l, err := srv.getListener()
	if err == nil {
		err = srv.Server.Serve(l)
	} else {
		fmt.Println(err)
	}
	return err
}

//获取listener
func (this *server) getListener() (*listener, error) {
	var l net.Listener
	var err error
	if os.Getenv("_GRACEFUL_RESTART") == "true" { //grace restart出来的进程，从FD FILE获取
		f := os.NewFile(3, "")
		l, err = net.FileListener(f)
		syscall.Kill(syscall.Getppid(), syscall.SIGTERM) //发信号给父进程，让父进程停止服务
	} else { //初始启动，监听addr
		l, err = net.Listen("tcp", this.Addr)
	}
	if err == nil {
		this.listener = &listener{
			Listener: l,
			server:   this,
		}
	}
	return this.listener, err
}

//fork一个新的进程
func (this *server) fork() error {
	os.Setenv("_GRACEFUL_RESTART", "true")
	lFd, err := this.listener.File()
	if err != nil {
		return err
	}
	execSpec := &syscall.ProcAttr{
		Env:   os.Environ(),
		Files: []uintptr{os.Stdin.Fd(), os.Stdout.Fd(), os.Stderr.Fd(), lFd},
	}
	pid, err := syscall.ForkExec(os.Args[0], os.Args, execSpec)
	if err != nil {
		return err
	}
	savePid(pid)
	fmt.Printf("\n %c[0;48;32m%s%c[0m", 0x1B, "["+strconv.Itoa(pid)+"] "+appName+" forked ok", 0x1B)
	return nil
}

//关闭服务
func (this *server) shutdown() {
	this.SetKeepAlivesEnabled(false)
	this.cm.close(TimeDeadLine)
	this.listener.Close()
	log.Printf("[%d] %s stopped.", os.Getpid(), appName)
}

//检查pidFile是否存在以及文件里的pid是否存活
func isRunning() bool {
	if mf, err := os.Open(pidFile); err == nil {
		pid, _ := ioutil.ReadAll(mf)
		pidVal, _ = strconv.Atoi(string(pid))
	}
	running := false
	if pidVal > 0 {
		if err := syscall.Kill(pidVal, 0); err == nil { //发一个信号为0到指定进程ID，如果没有错误发生，表示进程存活
			running = true
		}
	}
	return running
}

//保存pid
func savePid(pid int) error {
	file, err := os.OpenFile(pidFile, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()
	file.WriteString(strconv.Itoa(pid))
	return nil
}

//捕获系统信号
func handleSignals() {
	signals := make(chan os.Signal)
	signal.Notify(signals, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGINT)
	var err error
	for {
		sig := <-signals
		switch sig {
		case syscall.SIGHUP: //重启
			if srv != nil {
				err = srv.fork()
			} else { //only deamon时不支持kill -HUP,因为可能监听地址会占用
				log.Printf("[%d] %s stopped.", os.Getpid(), appName)
				os.Remove(pidFile)
				os.Exit(2)
			}
			if err != nil {
				log.Fatalln(err)
			}
		case syscall.SIGINT:
			fallthrough
		case syscall.SIGTERM:
			log.Printf("[%d] %s stop graceful", os.Getpid(), appName)
			if srv != nil {
				srv.shutdown()
			} else {
				log.Printf("[%d] %s stopped.", os.Getpid(), appName)
			}
			os.Exit(1)
		}
	}
}

func DaemonInit(cmd string) {
	// 得到存放pid文件的路径
	//pdiFile := BasePath +"/"+Env.Get("PID_FILE")
	//dir, _ := os.Getwd()
	//pidFile = dir + "/" + Env.Get("PID_FILE")
	//if os.Getenv("__Daemon") != "true" { //master
	//	cmd := "start" //缺省为start
	//	if l := len(os.Args); l > 2 {
	//		cmd = os.Args[l-1]
	//	}
		switch cmd {
		case "start":
			if isRunning() {
				fmt.Printf("\n %c[0;48;34m%s%c[0m", 0x1B, "["+strconv.Itoa(pidVal)+"] Bingo is running", 0x1B)
			} else { //fork daemon进程
				if err := forkDaemon(); err != nil {
					fmt.Println(err)
				}
			}
		case "restart": //重启:
			if !isRunning() {
				fmt.Printf("\n %c[0;48;31m%s%c[0m", 0x1B, "[Warning]bingo not running", 0x1B)
				restart(pidVal)
			} else {
				fmt.Printf("\n %c[0;48;34m%s%c[0m", 0x1B, "["+strconv.Itoa(pidVal)+"] Bingo restart now", 0x1B)
				restart(pidVal)
			}
		case "stop": //停止
			if !isRunning() {
				fmt.Printf("\n %c[0;48;31m%s%c[0m", 0x1B, "[Warning]bingo not running", 0x1B)
			} else {
				syscall.Kill(pidVal, syscall.SIGTERM) //kill
			}
		case "-h":
			fmt.Println("Usage: " + appName + " start|restart|stop")
		default:   //其它不识别的参数
			return //返回至调用方
		}
		//主进程退出
		os.Exit(0)
	//}
	go handleSignals()
}

//forkDaemon,当checkPid为true时，检查是否有存活的，有则不执行
func forkDaemon() error {
	args := os.Args
	os.Setenv("__Daemon", "true")
	procAttr := &syscall.ProcAttr{
		Env:   os.Environ(),
		Files: []uintptr{os.Stdin.Fd(), os.Stdout.Fd(), os.Stderr.Fd()},
	}
	pid, err := syscall.ForkExec(args[0], []string{args[0], "dev"}, procAttr)
	if err != nil {
		panic(err)
	}
	savePid(pid)
	fmt.Printf("\n %c[0;48;32m%s%c[0m", 0x1B, "["+strconv.Itoa(pid)+"] Bingo running...", 0x1B)
	fmt.Println()
	return nil
}

//重启(先发送kill -HUP到运行进程，手工重启daemon ...当有运行的进程时，daemon不启动)
func restart(pid int) {

	//dir, _ := os.Getwd()
	//pidFile = dir + "/" + Env.Get("PID_FILE")
	//if mf, err := os.Open(pidFile); err == nil {
	//	pid, _ := ioutil.ReadAll(mf)
	//	pidVal, _ = strconv.Atoi(string(pid))
	//}
	//// 判断当前系统中
	//p, err := os.FindProcess(pid)
	//if err != nil {
	//	// 不存在这个pid
	//	panic(err)
	//}
	//// 存在进程，杀掉
	//fmt.Println(p.Pid)
	//if p.Pid > 0 {
	syscall.Kill(pid, syscall.SIGHUP) //kill -HUP, daemon only时，会直接退出
	//
	//}
	forkDaemon()
	//fork := make(chan bool, 1)
	//go func() { //循环，查看pidFile是否存在，不存在或值已改变，发送消息
	//	for {
	//		f, err := os.Open(pidFile)
	//		if err != nil || os.IsNotExist(err) { //文件已不存在
	//			fork <- true
	//			break
	//		} else {
	//			pidVal, _ := ioutil.ReadAll(f)
	//			fmt.Println(string(pidVal))
	//			fmt.Println(strconv.Itoa(pid))
	//			if strconv.Itoa(pid) != string(pidVal) {
	//				fork <- false
	//				break
	//			}
	//		}
	//		time.Sleep(500 * time.Millisecond)
	//	}
	//}()
	////处理结果
	//select {
	//case r := <-fork:
	//	fmt.Println(r)
	//	if r == true {
	//		forkDaemon()
	//	}
	//case <-time.After(time.Second * 5):
	//	log.Fatalln("restart timeout")
	//}

}

//获取sock文件句柄
func (this *listener) File() (uintptr, error) {
	f, err := this.Listener.(*net.TCPListener).File()
	if err != nil {
		return 0, err
	}
	return f.Fd(), nil
}
