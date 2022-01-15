package app

import (
	"bytes"
	"github.com/urfave/cli/v2"
	"os"
)

func RunAppForTesting(args []string) (error, string) {
	w := new(bytes.Buffer)

	program := os.Args[0:1]
	args = append(program, args...)

	app := BuildApp()
	app.ExitErrHandler = func(context *cli.Context, err error) {
		// Do nothing
	}
	app.Writer = w

	err := app.Run(args)

	return err, w.String()
}
