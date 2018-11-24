package controllers

import (
	"github.com/silsuer/bingo-router"
)

type ${name} struct {
	bingo_router.Controller
}

// get -> path
func (e *${name}) Index(c *bingo_router.Context) {

}


// get -> path/create
func (e *${name}) Create(c *bingo_router.Context) {

}

// post -> path
func (e *${name}) Store(c *bingo_router.Context) {

}

// get -> path/:id/edit
func (e *${name}) Edit(c *bingo_router.Context) {

}

// put/patch -> path/:id
func (e *${name}) Update(c *bingo_router.Context) {

}

// delete -> path/:id
func (e *${name}) Delete(c *bingo_router.Context) {

}