package status

import (
	"fmt"
	"github.com/urfave/cli/v2"
)

func RegisterCommand() *cli.Command {
	var command *cli.Command
	command = new(cli.Command)
	command.Name = "status"
	command.Usage = "Query the status for a specific geocoding job"
	command.Action = status

	return command
}

func status(c *cli.Context) error {
	fmt.Println("Status: ", c.Args().First())
	return nil
}
