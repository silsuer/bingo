package bingo

import (
	"sync"
	"strings"
	"os"
	"github.com/kylelemons/go-gypsy/yaml"
	"errors"
)

var Config = &ConfigMap{}

type ConfigMap struct {
	sync.RWMutex
	data      map[string]string   // map不是线程安全的，所以需要加锁
	dataSlice map[string][]string // 这里记录是数组的配置
}

// 配置文件
// 记录一个全局变量
// 写2个函数 setter getter

func (c *ConfigMap) Set(key string, value string) bool {
	// a.b.c  那么就要设置多维数组
	// 加写锁
	c.Lock()
	c.data[key] = value
	c.Unlock()
	return true
}

func (c *ConfigMap) SetSlice(key string, value []string) bool {
	c.Lock()
	c.dataSlice[key] = value
	c.Unlock()
	return true
}

func (c *ConfigMap) Get(key string) (string, error) {
	// 判断内存中是否存在，存在的话直接返回，否则去加载yaml文件
	var flag bool
	// 加读锁
	c.RLock()
	if _, ok := c.data[key]; ok {
		flag = true
	} else {
		flag = false
	}
	c.RUnlock()
	if flag {
		return c.data[key], nil
	} else {
		rootPath, _ := os.Getwd()
		configRootPath := Env.Get("CONFIG_ROOT")
		configPath := rootPath + "/" + configRootPath
		// 从yaml文件中读取，先判断是否存在这个文件 a.b.c.d  先判断文件是否存在，直到文件不存在，去判断内部数据
		fileSlice := strings.Split(key, ".") // 以 . 号分割
		for i := 0; i < len(fileSlice); i++ {
			if CheckFileIsExist(configPath + "/" + fileSlice[i]) { // 不存在文件夹
				// 如果存在，继续检查下一个
				configPath = configPath + "/" + fileSlice[i]
			} else {
				// 不存在文件夹，则先判断文件是否存在，则设现在的路径就是配置文件的路径
				// 如果是文件
				configPath = configPath + "/" + fileSlice[i] + ".yaml"
				f, err := os.Stat(configPath)
				if err != nil {
					return "", errors.New("no such file exists:" + configPath)
				}
				if !f.IsDir() { // 文件存在
					file, err := yaml.ReadFile(configPath)
					Check(err)
					result, err := file.Get(strings.Join(fileSlice[(i + 1):], "."))
					if err != nil {
						c.Set(key, result)
					}
					return result, err
				} else {
					return "", errors.New("no such file exists:" + configPath)
				}
			}
		}
		return "", nil
	}
}
