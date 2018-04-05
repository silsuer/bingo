package controller

import (
	"net/http"
	"fmt"
)

// 控制器，处理数据的逻辑
type Controller struct{}

// 所有的方法都挂载到这里控制器下
func (c *Controller) Index(w http.ResponseWriter,r *http.Request)  {
	fmt.Fprintln(w,"Helloworld")
	return
}

// 每一个方法名，必须大写
func (c *Controller) Asd(w http.ResponseWriter,r *http.Request)  {
	fmt.Fprintln(w,"这是asd方法")
	return
}