package output

import (
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
)

func WarningAndExit(message string) error {
	return cli.Exit(color.YellowString(message), 0)
}
