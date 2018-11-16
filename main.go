package bingo

import (
	"github.com/urfave/cli"
	"github.com/silsuer/bingo-router"
)

type Bingo struct {
	Cli    *cli.App
	Router *bingo_router.Router
}

// 创建一个app
func NewApp() *Bingo {
	b := &Bingo{
		Cli:    cli.NewApp(),
		Router: bingo_router.New(),
	}
	return b
}
