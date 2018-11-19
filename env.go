package bingo

import (
	"sync"
	"github.com/kylelemons/go-gypsy/yaml"
)

var env *Environment
var once sync.Once

type Environment struct {
	sync.RWMutex
	values map[string]string
	path   string
	file   *yaml.File
}

// 创建一个单例
func GetInstance(path string) *Environment {
	once.Do(func() {
		env = new(Environment)
		env.path = path
		env.values = make(map[string]string)
		f, err := yaml.ReadFile(path)
		if err != nil {
			panic(err)
		}
		env.file = f
	})
	return env
}

func (env *Environment) Set(key, value string) bool {
	env.Lock()
	defer env.Unlock()
	env.values[key] = value
	return true
}

func (env *Environment) Get(key string) string {
	// 先判断map里有没有，如果有直接返回，如果没有，从文件中取出
	if v, ok := env.values[key]; ok {
		return v
	}
	v, err := env.file.Get(key)
	if err != nil {
		panic(err)
	}
	return v
}
