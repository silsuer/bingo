package main

import (
	"fmt"
	"github.com/urfave/cli"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func makeCommand(c *cli.Context) error {
	if c.Args().Get(0) == "" {
		log.Fatal("Need a argument: command's name.  e.g. bingo sword make:command example_command")
	}

	// tmpl文件路径
	tmplDir := getTempDir() + "/command.gtpl"
	if IfFileExist(tmplDir) { // 文件存在
		f, err := os.Open(tmplDir)
		defer f.Close()
		if err != nil {
			log.Fatal(err)
		}
		bytes, err := ioutil.ReadAll(f)
		if err != nil {
			log.Fatal(err)
		}
		var str string
		str = strings.Replace(string(bytes), "${name}", c.Args().Get(0), -1)

		// 得到要创建的env文件
		p, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		dstDir := p + "/" + new(Environment).GetWithDefault("COMMAND_DIR", "console/commands") + "/" + c.Args().Get(0) + ".go"

		df, err := CreateFile(dstDir, 0755)
		defer df.Close()
		if err != nil {
			log.Fatal(err)
		}
		io.WriteString(df, str)
		fmt.Printf("\n%c[0;48;32m%s%c[0m\n", 0x1B, "Command created successfully :"+dstDir, 0x1B)

	} else {
		log.Fatal("No such file named command.gtpl in " + tmplDir)

	}
	return nil
}
