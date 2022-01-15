package main

import (
	cli "github.com/geocodio/geocodio-cli/app"
	"github.com/geocodio/geocodio-cli/output"
	"os"
)

func main() {
	app := cli.BuildApp()

	err := app.Run(os.Args)
	if err != nil {
		output.ErrorAndExit(err)
	}
}

