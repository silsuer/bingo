package main

import (
	"github.com/urfave/cli"
	"os"
)

// 这个文件用来存储所有命令
func commands() []cli.Command {
	return []cli.Command{
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
		{
			Name:  "run",
			Usage: "Start a server",
			Subcommands: []cli.Command{
				{
					Name:   "dev",
					Usage:  "Start a development server. Listen the current directory.",
					Action: Dev,
				},
			},
		},
		{
			Name:        "sword",
			Subcommands: sword(),
		},
	}
}
