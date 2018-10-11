package bingo

import (
	"github.com/kylelemons/go-gypsy/yaml"
	"sync"
)

type env struct {
	sync.RWMutex
	submap map[string]string
}

var Env *env

func init() {
	// 读取env.yaml中的所有数据，初始化Env
	Env = &env{submap: make(map[string]string)}
}

func (e *env) Set(k string, v string) bool {
	e.Lock()
	defer e.Unlock()

	e.submap[k] = v
	return true
}

func (e *env) Get(k string) string {
	//输入键，返回值
	if v, ok := e.submap[k]; ok {
		// 如果存在，直接返回返回值
		return v
	} else {
		// 如果不存在，读取env文件，并且把数据存入Env中
		envConfig, err := yaml.ReadFile(BasePath + "/env.yaml")
		val, err := envConfig.Get(k)
		// 写入内存
		Check(err)
		e.submap[k] = val
		return val
	}
}
