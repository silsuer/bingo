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

func sword() []cli.Command {
	return []cli.Command{
		{
			Name:  "make:controller",
			Usage: "Create a controller file from the template.",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name: "resource",
				},
			},
			Action: func(c *cli.Context) error {

				if c.Args().Get(0) == "" {
					log.Fatal("Need a argument: controller's name.  e.g. bingo sword make:controller example_controller")
				}

				//从templates文件夹中把控制器文件复制过来
				//控制器模板文件夹路径
				var controllerDir string
				if c.Bool("resource") == true {
					controllerDir = getTempDir() + "/resource_controller.gtpl"
				} else {
					controllerDir = getTempDir() + "/controller.gtpl"
				}

				if IfFileExist(controllerDir) == true {
					// 读取文件内容，创建文件
					f, err := os.Open(controllerDir)
					if err != nil {
						log.Fatal(err)
					}
					defer f.Close()
					bytes, err := ioutil.ReadAll(f)
					if err != nil {
						log.Fatal(err)
					}

					str := strings.Replace(string(bytes), "${name}", c.Args().Get(0), -1)

					// 得到要创建的env文件
					p, err := os.Getwd()
					if err != nil {
						log.Fatal(err)
					}
					dstDir := p + "/" + new(Environment).GetWithDefault("CONTROLLER_DIR", "http/controllers") + "/" + c.Args().Get(0) + ".go"

					// 创建文件并返回句柄
					df, err := CreateFile(dstDir, 0755)
					if err != nil {
						log.Fatal(err)
					}
					defer df.Close()

					io.WriteString(df, str)

					fmt.Printf("\n %c[0;48;32m%s%c[0m\n\n", 0x1B, "Controller created successfully :"+dstDir, 0x1B)
				} else {
					//不存在，则弹出提示
					log.Fatal("No such file named controller.gtpl in " + controllerDir)
				}

				return nil
			},
		},
	}
}
