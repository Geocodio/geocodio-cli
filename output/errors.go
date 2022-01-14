package output

import (
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
)

func ErrorAndExit(err error) error {
	message := fmt.Sprintf("%s %s", color.RedString("❌ Error:"), err)
	return cli.Exit(message, 1)
}

func ErrorStringAndExit(message string) error {
	return ErrorAndExit(errors.New(message))
}
