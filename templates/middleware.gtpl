package middlewares

import (
	"github.com/silsuer/bingo-router"
)


func ${name}(c *bingo_router.Context, next func(c *bingo_router.Context)) {
   next(c)
}
