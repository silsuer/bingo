package main

import (
	"github.com/urfave/cli"
)

func sword() []cli.Command {
	return []cli.Command{
		{
			Name:  "make:controller",
			Usage: "Create a controller file from the bingo template.",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name: "resource",
				},
			},
			Action: makeController,
		},
		{
			Name:   "make:command",
			Usage:  "Create a command file from the bingo template.",
			Action: makeCommand,
		},
	}
}
