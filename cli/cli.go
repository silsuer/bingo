package cli

import (
	"flag"
	"os"
	"fmt"
	"errors"
	"github.com/silsuer/bingo/bingo"
	"os/exec"
	"bytes"
	"log"
)

type CLI struct{}

func (cli *CLI) Run() {
	//cli.validateArgs()
	initCmd := flag.NewFlagSet("init", flag.ExitOnError)
	runCmd := flag.NewFlagSet("run", flag.ExitOnError)
	//runPort := runCmd.Int("port", 8088, "监听端口") // 默认是8088端口

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
		cli.RunSite()
	}

	if swordCmd.Parsed() {
		cli.swordHandle(swordConfig)
	}
}

func (cli *CLI) RunSite() {
	// 这里应该运行 go run start.go
	cmd := exec.Command("go", "run", "start.go")

	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Start(); err != nil {
		panic(err)
	}

	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}
	fmt.Println(out.String())
}

func ThrowError(text string) {
	panic(errors.New(text))
}
