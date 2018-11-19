package controllers

import (
	"github.com/silsuer/bingo-router"
	"html/template"
)

func Index(c *bingo_router.Context) {
	t, err := template.ParseFiles("./resources/views/index.html")
	if err != nil {
		panic(err)
	}
	err = t.Execute(c.Writer, nil) // 渲染模板
	if err != nil {
		panic(err)
	}
}
