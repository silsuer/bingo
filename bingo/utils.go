package bingo

import (
	"os"
	"path/filepath"
	"strings"
	"io/ioutil"
	"errors"
	"bingo/drivers/db/mysql"
)



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

// 获取路由文件夹路径
func GetRoutesPath() string {

	dir, err := os.Getwd()
	dir = strings.Replace(dir, "\\", "/", -1)
	Check(err)
	//// 获取env文件的所有值
	//c := FileGetJson(dir + "/.env")
	//return c
	//route_path,err := Env("ROUTES_PATH")
	route_path := Env.Get("ROUTES_PATH")
	Check(err)
	path := dir + "/" + route_path
	return path
}

// 获取静态文件夹路径
func GetPublicPath() string {
	dir, err := os.Getwd()
	dir = strings.Replace(dir, "\\", "/", -1)
	Check(err)
	return dir+"/" + Env.Get("STATIC_FILE_DIR")
}

func DB() interface{} {
	// 返回一个驱动的操作类
	// 传入db配置  env的driver，实例化不同的类,单例模式，获取唯一驱动
	if Driver == nil {
		DriverInit()
	}
	con := Driver.GetConnection()  // 获取数据库连接
	return con
}

func MySqlDB() *mysql.Mysql  {
	return DB().(*mysql.Mysql)
}


/**
 * 判断文件是否存在  存在返回 true 不存在返回false
 */
func CheckFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}


