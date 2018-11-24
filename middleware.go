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

func makeMiddleware(c *cli.Context) error {
	if c.Args().Get(0) == "" {
		log.Fatal("Need a argument: middleware's name.  e.g. bingo sword make:middleware example_middleware")
	}
	tmpDir := getTempDir() + "/middleware.gtpl"
	if IfFileExist(tmpDir) {
		// 文件存在
		f, err := os.Open(tmpDir)
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
		dstDir := p + "/" + new(Environment).GetWithDefault("MIDDLEWARE_DIR", "http/middlewares") + "/" + c.Args().Get(0) + ".go"

		df, err := CreateFile(dstDir, 0755)
		defer df.Close()
		if err != nil {
			log.Fatal(err)
		}
		io.WriteString(df, str)
		fmt.Printf("\n%c[0;48;32m%s%c[0m\n", 0x1B, "Middleware created successfully :"+dstDir, 0x1B)


	} else {
		log.Fatal("No such file named middleware.gtpl in " + tmpDir)
	}
	return nil
}
