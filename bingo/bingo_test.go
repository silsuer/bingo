package bingo

import (
	"testing"
	"net/http"
	"io/ioutil"
	"strconv"
	"sync"
	"fmt"
)

var s sync.Once
// 压力测试

func Benchmark_Bingo_Run(b *testing.B) {
	// 测试
	b.StopTimer()

	// 要指定env路径

	// 启动一个服务器
	bin := &Bingo{}

	s.Do(func() {
		NewRoute().Get("/").Target(func(c *Context) {
			fmt.Fprint(c.Writer, "hello bingo!")
		}).Register()
		go bin.Run(":12345", []string{"dev", "path=/Users/silsuer/go/src/github.com/silsuer/bingo"})
	})

	url := "http://localhost:12345"

	b.StartTimer() // 开始计时

	//fmt.Println(111)
	for i := 0; i < 10000; i++ {
		//go func() {
			//fmt.Println("开始第 " + strconv.Itoa(i) + " 次请求")
			// 进行访问
			resp, err := http.Get(url)

			if err != nil {
				b.Error(err)
			}

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				b.Error(err)
			}
			b.Log("第" + strconv.Itoa(i) + "次请求:" + string(body))
			resp.Body.Close()
		//}()

	}
}
