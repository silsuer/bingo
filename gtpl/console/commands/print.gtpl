package commands

import (
	"github.com/urfave/cli"
	"fmt"
)

var Test cli.Command

func init() {
	Test = cli.Command{
		Name: "test",
		Action: func(c *cli.Context) error {
			fmt.Println("this is a test command")
			return nil
		},
	}
}
