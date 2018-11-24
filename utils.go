package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func getEnvPath() (string, bool) {
	// 返回值 string代表路径，bool代表是否存在
	p, _ := os.Getwd()
	path := p + "/.env.yaml" // 拼接路径

	// 如果存在
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return path, false
		}
	}
	return path, true
}

func getTempDir() string {
	return os.Getenv("GOPATH") + "/src/github.com/silsuer/bingo/templates"
}

func IfFileExist(dir string) bool {
	_, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		} else {
			log.Fatal(err)
		}
	}
	return true
}

// 获取替换时的要替换的包名
func getReplacePackage() string {
	path, _ := os.Getwd()
	goSrcPath := os.Getenv("GOPATH") + "/src/"
	if !strings.Contains(path, goSrcPath) { // 包含这个路径
		fmt.Printf("You must use dys command in the $GOPATH directory.\n")
		fmt.Printf("\tNow the current path: " + path + "\n")
		fmt.Println("\tNow the $GOPATH: " + goSrcPath + "\n")
		return ""
	}

	p := path[len(goSrcPath):]
	return p
}

// 创建文件，递归创建文件目录，并且返回文件句柄
func CreateFile(path string, permision os.FileMode) (*os.File, error) {
	dirPath := filepath.Dir(path)
	err := os.MkdirAll(dirPath, permision)
	if err != nil {
		return nil, err
	}
	return os.Create(path)
}
