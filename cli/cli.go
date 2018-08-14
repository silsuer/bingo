package cli

import (
	"flag"
	"os"
	"fmt"
	"errors"
	"github.com/silsuer/bingo/bingo"
	"os/exec"
	"log"
	"io/ioutil"
	"io"
	"bufio"
)

type CLI struct{}

func (cli *CLI) Run() {
	//cli.validateArgs()
	initCmd := flag.NewFlagSet("init", flag.ExitOnError)
	runCmd := flag.NewFlagSet("run", flag.ExitOnError)

	swordCmd := flag.NewFlagSet("sword", flag.ExitOnError) // bingo sword 命令
	//swordConfig := swordCmd.String("name","list","the Commands name")
	var swordConfig []string
	switch os.Args[1] {
	case "init":
		err := initCmd.Parse(os.Args[2:]) // 把其余参数传入进去
		bingo.Check(err)
		break
	case "run":
		err := runCmd.Parse(os.Args[2:])
		bingo.Check(err)
		break
	case "sword":
		err := swordCmd.Parse(os.Args[2:])
		swordConfig = os.Args[2:]
		bingo.Check(err)
		break
	default:
		break
	}

	if initCmd.Parsed() {
		// 如果是init命令
		cli.InitProject()
	}

	if runCmd.Parsed() {
		cli.RunSite(os.Args[2:])
	}

	if swordCmd.Parsed() {
		cli.swordHandle(swordConfig)
	}
}

func (cli *CLI) RunSite(arg []string) {
	// bingo run watch   监听文件变动，有变动即平滑重启服务
	// bingo run daemon
	// bingo run production
	// 这里应该运行 go run start.go
	tmpSlice := []string{"run", "start.go"}
	cmd := exec.Command("go", append(tmpSlice, arg...)...)

	stdout, _ := cmd.StdoutPipe()
	defer stdout.Close()

	if err := cmd.Start(); err != nil {
		panic(err)
	}

	reader := bufio.NewReader(stdout)

	//实时循环读取输出流中的一行内容
	for {
		line, err2 := reader.ReadString('\n')
		if err2 != nil || io.EOF == err2 {
			break
		}
		fmt.Print(line)
	}

	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}
	opBytes, _ := ioutil.ReadAll(stdout)
	fmt.Print(string(opBytes))
}

func ThrowError(text string) {
	panic(errors.New(text))
}
