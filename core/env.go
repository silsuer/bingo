package core

import (
	"github.com/kylelemons/go-gypsy/yaml"
	"os"
	"fmt"
)

type env struct {
	submap map[string]string
}

var Env *env

func init() {
	// 读取env.yaml中的所有数据，初始化Env
	fmt.Println("正在加载env.yaml文件")
	Env = &env{make(map[string]string)}
}

func (e *env) Set(k string, v string) bool {
	return true
}

func (e *env) Get(k string) string {
	//输入键，返回值
	if v, ok := e.submap[k]; ok {
		// 如果存在，直接返回返回值
		return v
	}else{
		// 如果不存在，读取env文件，并且把数据存入Env中
		path,err := os.Getwd()
		Check(err)
		env_config,err := yaml.ReadFile(path+"/env.yaml")
		fmt.Println(k)
		val,err := env_config.Get(k)
		fmt.Println(val)
		// 写入内存
		e.submap[k] = val
		fmt.Println(e)
		return val
	}
}
