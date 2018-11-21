package console

import (
	"github.com/urfave/cli"
	"${path}/console/commands"
)

func Schedule() []cli.Command {
	return []cli.Command{
		commands.Test,
	}
}
