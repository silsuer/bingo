package commands

import (
	"github.com/urfave/cli"
	"fmt"
)

var ${name} cli.Command

func init() {
	Test = cli.Command{
		Name: "",
		Action: func(c *cli.Context) error {
			return nil
		},
	}
}
