package main

import (
	"github.com/urfave/cli"
	"os"
	"fmt"
)

// 创建
func CreateProject(c *cli.Context) error {
	//c.Args()
	//fmt.Println(c.Args()[0])
	if c.NArg() <= 0 {
		fmt.Println("The Create Command need a argument: projectName")
		return nil
	}

	path, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		return err
	}

	fullPath := path + "/" + c.Args()[0]
	// 创建文件夹
	err = os.MkdirAll(fullPath, 0755)
	if err != nil {
		fmt.Println(err)
		return err
	}
	initProject(fullPath)
	return nil
}
