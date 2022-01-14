package main

import (
	"github.com/geocodio/geocodio-cli/create"
	"github.com/geocodio/geocodio-cli/list"
	"github.com/geocodio/geocodio-cli/release"
	"github.com/geocodio/geocodio-cli/remove"
	"github.com/geocodio/geocodio-cli/status"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "Geocodio"
	app.Usage = "Geocode lists using the Geocodio API"
	app.Version = release.Version()
	app.EnableBashCompletion = true
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "hostname",
			Aliases: []string{"n"},
			Value:   "api.geocod.io",
			Usage:   "Geocodio hostname to use, change this for Geocodio+HIPAA or on-premise environments",
			EnvVars: []string{"GEOCODIO_HOSTNAME"},
		},
		&cli.StringFlag{
			Name:     "apikey",
			Aliases:  []string{"k"},
			Value:    "",
			Usage:    "Geocodio API Key to use. Generate a new one in the Geocodio Dashboard",
			EnvVars:  []string{"GEOCODIO_API_KEY"},
			Required: true,
		},
	}
	app.Commands = []*cli.Command{
		create.RegisterCommand(),
		list.RegisterCommand(),
		status.RegisterCommand(),
		remove.RegisterCommand(),
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
