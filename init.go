package main

import (
	"os"
	"strings"
	"fmt"
	"path/filepath"
	"io/ioutil"
	"io"
)

// 初始化项目
func initProject(path string) {
	// 得到当前在src下的路径
	// 得到当前路径
	//path, err := os.Getwd()
	//if err != nil {
	//	panic(err)
	//}
	goSrcPath := os.Getenv("GOPATH") + "/src/"
	if !strings.Contains(path, goSrcPath) { // 包含这个路径
		fmt.Printf("You must use dys command in the $GOPATH directory.\n")
		fmt.Printf("\tNow the current path: " + path + "\n")
		fmt.Println("\tNow the $GOPATH: " + goSrcPath + "\n")
		return
	}

	p := path[len(goSrcPath):]
	// 从当前路径下查找bingo_template
	bingoTemplatePath := os.Getenv("GOPATH") + `/src/github.com/silsuer/bingo/gtpl`
	// 如果是dir
	stat, err := os.Stat(bingoTemplatePath)

	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println(`Can not find the template directory in ` + bingoTemplatePath)
			fmt.Println(`Please use go get -u github.com/silsuer/bingo...`)
			return
		} else {
			panic(err)
		}
	}
	if !stat.IsDir() {
		fmt.Println(`No such directory in ` + bingoTemplatePath)
		return
	}

	copyDir(bingoTemplatePath, path, p)
	// 打印欢迎信息
	welecome()
}

func welecome() {
	info := `
__        _______ _     ____ ___  __  __ _____     ____ ___ _   _  ____  ___  _ 
\ \      / / ____| |   / ___/ _ \|  \/  | ____|   | __ )_ _| \ | |/ ___|/ _ \| |
 \ \ /\ / /|  _| | |  | |  | | | | |\/| |  _|     |  _ \| ||  \| | |  _| | | | |
  \ V  V / | |___| |__| |__| |_| | |  | | |___    | |_) | || |\  | |_| | |_| |_|
   \_/\_/  |_____|_____\____\___/|_|  |_|_____|   |____/___|_| \_|\____|\___/(_)
  
  Now you can start a development server with: make dev 
`
	fmt.Printf("\n %c[0;48;32m%s%c[0m\n\n", 0x1B, info, 0x1B)
}

func copyDir(src string, dest string, variable string) {
	src_original := src
	err := filepath.Walk(src, func(src string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() {
			//			fmt.Println(f.Name())
			//			copyDir(f.Name(), dest+"/"+f.Name())
		} else {
			dest_new := strings.Replace(src, src_original, dest, -1)
			CopyFile(src, dest_new, variable)
		}
		//println(path)
		return nil
	})
	if err != nil {
		fmt.Printf("filepath.Walk() returned %v\n", err)
	}
}

//egodic directories
func getFilelist(path string) {
	err := filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() {
			return nil
		}
		println(path)
		return nil
	})
	if err != nil {
		fmt.Printf("filepath.Walk() returned %v\n", err)
	}
}
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

//copy file
func CopyFile(src, dst string, variable string) (w int, err error) {
	srcFile, err := os.Open(src)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer srcFile.Close()
	dst_slices := strings.Split(dst, string(os.PathSeparator))
	dst_slices_len := len(dst_slices)
	dest_dir := ""
	for i := 0; i < dst_slices_len-1; i++ {
		dest_dir = dest_dir + dst_slices[i] + string(os.PathSeparator)
	}
	b, err := PathExists(dest_dir)
	if b == false {
		err := os.MkdirAll(dest_dir, os.ModePerm) //在当前目录下生成md目录
		if err != nil {
			fmt.Println(err)
		}
	}

	ext := filepath.Ext(dst)
	if ext == ".gtpl" {
		//oldName := filepath.Base(dst)
		//newName := oldName[:len(filepath.Ext(dst))]
		newNameSlice := strings.Split(filepath.Base(dst), ".")
		newName := strings.Join(newNameSlice[:len(newNameSlice)-1], ".")
		dst = filepath.Dir(dst) + "/" + newName + ".go"
	}
	dstFile, err := os.Create(dst)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	bytes, err := ioutil.ReadAll(srcFile)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// 获取srcFile中的所有数据，并替换变量
	defer dstFile.Close()
	var str string
	if fName := filepath.Base(dst); fName == "Makefile" || fName == ".gitignore" {
		// 如果是makefile，则把其中的 ${name} 替换成 当前目录
		// 否则再进行判断后缀名的操作
		// 如果是.gitignore,则也进行替换
		v := strings.Split(variable, "/")
		vv := v[len(v)-1]
		str = strings.Replace(string(bytes), "${name}", vv, -1)
	} else {
		str = strings.Replace(string(bytes), "${path}", variable, -1)
	}

	// 这个文件写入dstFile中
	return io.WriteString(dstFile, str)
}
