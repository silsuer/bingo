package main

import (
	"github.com/urfave/cli"
	"os"
)

func main() {

	// bingo run dev
	// bingo sword make...
	app := cli.NewApp()
	app.Commands = []cli.Command{
		{
			Name:  "init",
			Usage: "Init a bingo project.",
			Action: func(c *cli.Context) error {
				path, _ := os.Getwd()
				initProject(path)
				return nil
			},
		},
		{
			Name:  "create",
			Usage: "Create a New Project Directory.",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "projectName",
					Value: "bingo_project",
					Usage: "The new project name.",
				},
			},
			Action: CreateProject,
		},
	}

	app.Run(os.Args)
}
