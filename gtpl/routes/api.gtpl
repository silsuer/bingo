package routes

import (
	"github.com/silsuer/bingo-router"
	"${path}/http/middlewares"
	"${path}/http/controllers"
)

func Api() *bingo_router.Route {
	return bingo_router.NewRoute().Middleware(middlewares.Log).Mount(func(b *bingo_router.Builder) {
		b.NewRoute().Get("/").Target(controllers.Index)
	})
}
