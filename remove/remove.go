package remove

import (
	"fmt"
	"github.com/urfave/cli/v2"
)

func RegisterCommand() *cli.Command {
	var command *cli.Command
	command = new(cli.Command)
	command.Name = "remove"
	command.Aliases = []string{"rm"}
	command.Usage = "Delete an existing geocoding job"
	command.Action = remove

	return command
}

func remove(c *cli.Context) error {
	fmt.Println("Remove: ", c.Args().First())
	return nil
}
