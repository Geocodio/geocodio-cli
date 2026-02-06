package cli

import (
	"context"
	"io"
	"os"

	"github.com/geocodio/geocodio-cli/internal/api"
	"github.com/geocodio/geocodio-cli/internal/config"
	"github.com/geocodio/geocodio-cli/internal/output"
	"github.com/urfave/cli/v3"
)

var Version = "dev"

type App struct {
	cfg       *config.Config
	client    *api.Client
	formatter output.Formatter
	stdout    io.Writer
	stderr    io.Writer
}

func NewApp() *cli.Command {
	return &cli.Command{
		Name:    "geocodio",
		Usage:   "Geocodio API command line interface",
		Version: Version,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "api-key",
				Usage:   "Geocodio API key (or set GEOCODIO_API_KEY)",
				Sources: cli.EnvVars("GEOCODIO_API_KEY"),
			},
			&cli.StringFlag{
				Name:  "base-url",
				Usage: "Override API base URL",
			},
			&cli.BoolFlag{
				Name:  "json",
				Usage: "Output as JSON",
			},
			&cli.BoolFlag{
				Name:  "debug",
				Usage: "Enable debug output",
			},
		},
		Commands: []*cli.Command{
			geocodeCmd(),
			reverseCmd(),
			distanceCmd(),
			distanceMatrixCmd(),
			distanceJobsCmd(),
			listsCmd(),
		},
	}
}

func newApp(cmd *cli.Command) (*App, error) {
	cfg := config.New(
		cmd.String("api-key"),
		cmd.String("base-url"),
		cmd.Bool("debug"),
	)

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	stdout := os.Stdout
	stderr := os.Stderr

	var clientOpts []api.ClientOption
	clientOpts = append(clientOpts, api.WithUserAgent("geocodio-cli/"+Version))
	if cfg.Debug {
		clientOpts = append(clientOpts, api.WithDebug(stderr))
	}

	client := api.NewClient(cfg.BaseURL, cfg.APIKey, clientOpts...)
	formatter := output.New(stdout, cmd.Bool("json"))

	return &App{
		cfg:       cfg,
		client:    client,
		formatter: formatter,
		stdout:    stdout,
		stderr:    stderr,
	}, nil
}

func Run(ctx context.Context, args []string) error {
	return NewApp().Run(ctx, args)
}
