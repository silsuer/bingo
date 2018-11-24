package main

import (
	"github.com/urfave/cli"
	"os"
)

func main() {

	// bingo run dev
	// bingo sword make...
	app := cli.NewApp()
	app.Commands = commands()

	app.Run(os.Args)
}
