package main

import (
	"github.com/kylelemons/go-gypsy/yaml"
	"log"
	"sync"
)

var Env *Environment
var s sync.Once

func getEnvInstance() *Environment {
	s.Do(func() {
		Env := new(Environment)
		Env.m = make(map[string]string)
	})
	return Env
}

type Environment struct {
	sync.RWMutex
	m map[string]string
}

func (e *Environment) Set(key, value string) {
	e.Lock()
	defer e.Unlock()

	e.m[key] = value
}

func (e *Environment) Get(key string) string {
	if v, exist := e.m[key]; exist {
		return v
	} else {
		// 从yaml中获取数据，如果存在 .env.yaml文件，则从里面查找
		path, exist := getEnvPath()
		if exist { // 存在，获取数据
			file, _ := yaml.ReadFile(path)
			res, err := file.Get(key)
			if err != nil {
				log.Fatal(err)
			}
			return res
		} else {
			// 不存在这个文件，返回空字符串
			return ""
		}

	}
	return ""
}

func (e *Environment) GetWithDefault(key, def string) string {
	if v, exist := e.m[key]; exist {
		return v
	} else {
		// 从yaml中获取数据，如果存在 .env.yaml文件，则从里面查找
		path, exist := getEnvPath()
		if exist { // 存在，获取数据
			file, _ := yaml.ReadFile(path)
			res, err := file.Get(key)
			if err != nil {
				return def
			}
			return res
		} else {
			// 不存在这个文件，返回空字符串
			return def
		}

	}
	return def
}
