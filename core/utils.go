package core

import (
	"os"
	"path/filepath"
	"strings"
	"io/ioutil"
	"errors"
	"bufio"
)

//type Json struct {
//	key   string
//	value string
//}
//
//type JsonMap struct {
//	Jsons []Json
//}

// 这里主要写自定义函数

func FileGetContents(filepath string) (string, error) {

	// 先检查文件是否存在
	status := FileIsExist(filepath)
	if status == false {
		return "", errors.New(filepath + "文件不存在!")
	}

	content, err := ioutil.ReadFile(filepath) // 直接读取到内存里
	Check(err)
	return string(content), nil
}

func FileGetJson(filepath string) string {

	//// 获取json文件的内容，并且返回json数据
	//content, err := FileGetContents(filepath)
	//Check(err)
	////var j JsonMap
	////jsons,err:= json.Unmarshal([]byte(content), &j)
	//js,err :=
	//Check(err)
	//return jsons
	return "111"
}

func Check(e error) {
	if e != nil {
		panic(e)
	}
}

/**
 * 判断文件是否存在  存在返回 true 不存在返回false
 */
func FileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

func GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	Check(err)
	return strings.Replace(dir, "\\", "/", -1)
}

// env文件,传入一个配置项名称，返回这个配置项的值
func Env(name string) (string,error) {

	// 获取env 路径
	dir, err := os.Getwd()
	dir = strings.Replace(dir, "\\", "/", -1) // 获取文件执行路径
	dir = dir + "/.env"
	// 打开env文件
	EnvFile, err := os.Open(dir)
	Check(err)
	defer EnvFile.Close() // 最后要关闭
	buf := bufio.NewReader(EnvFile)
	for {
		line, _, err := buf.ReadLine()
		Check(err)
		// 处理拿到的这一行数据
		newLine := strings.TrimSpace(string(line))
		// 如果第一个字符是# 的话，证明是注释，跳过这一行
		if string([]rune(newLine)[:1]) == "#" {
			continue
		}
		// 否则用=分割数据，拿到
        env := strings.Split(newLine,"=")
        if name==env[0] {
        	return env[1],nil
		}
	}

	return "",errors.New("env中不存在"+name +"配置项")
}

// 获取路由文件夹路径
func GetRoutesPath() string {

	dir, err := os.Getwd()
	dir = strings.Replace(dir, "\\", "/", -1)
	Check(err)
	//// 获取env文件的所有值
	//c := FileGetJson(dir + "/.env")
	//return c
	route_path,err := Env("ROUTES_PATH")
	Check(err)
	path := dir + "/" + route_path
	return path
}

// 获取静态文件夹路径
func GetPublicPath() string {
	dir, err := os.Getwd()
	dir = strings.Replace(dir, "\\", "/", -1)
	Check(err)
	return dir + "/public/"
}
