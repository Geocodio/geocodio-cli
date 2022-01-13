package create

import (
	"fmt"
	"github.com/urfave/cli/v2"
)

func RegisterCommand() *cli.Command {
	var command *cli.Command
	command = new(cli.Command)
	command.Name = "create"
	command.Usage = "Geocode a new spreadsheet"
	command.Action = geocode

	return command
}

func geocode(c *cli.Context) error {
	fmt.Println("Geocode: ", c.Args().First())
	return nil
}
