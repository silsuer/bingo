package routes

import (
	"fmt"
	"github.com/silsuer/bingo-router"
	"github.com/silsuer/bingo/core"
	"os"
)

func Api() *bingo_router.Route {
	return bingo_router.NewRoute().Mount(func(b *bingo_router.Builder) {
		b.NewRoute().Get("/").Target(func(c *bingo_router.Context) {
			fmt.Fprintln(c.Writer, "hello bingo!")
			fmt.Fprintln(c.Writer, os.Getenv(core.EnvVariable))
		})
	})
}
