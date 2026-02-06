package cli

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
)

func distanceCmd() *cli.Command {
	return &cli.Command{
		Name:      "distance",
		Usage:     "Calculate distance from origin to destinations",
		ArgsUsage: "<origin> <destination> [destination...]",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "mode",
				Aliases: []string{"m"},
				Usage:   "Routing mode: driving or straightline",
				Value:   "driving",
			},
			&cli.StringFlag{
				Name:    "units",
				Aliases: []string{"u"},
				Usage:   "Distance units: miles or km",
				Value:   "miles",
			},
		},
		Action: distanceAction,
	}
}

func distanceAction(ctx context.Context, cmd *cli.Command) error {
	app, err := newApp(cmd)
	if err != nil {
		return err
	}

	if cmd.NArg() < 2 {
		return fmt.Errorf("requires origin and at least one destination")
	}

	args := cmd.Args().Slice()
	origin := args[0]
	destinations := args[1:]

	resp, err := app.client.Distance(ctx, origin, destinations, cmd.String("mode"), cmd.String("units"))
	if err != nil {
		return err
	}

	return app.formatter.FormatDistance(resp)
}

func distanceMatrixCmd() *cli.Command {
	return &cli.Command{
		Name:  "distance-matrix",
		Usage: "Calculate distances between multiple origins and destinations",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "origins",
				Aliases:  []string{"o"},
				Usage:    "File containing origins (one per line)",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "destinations",
				Aliases:  []string{"d"},
				Usage:    "File containing destinations (one per line)",
				Required: true,
			},
			&cli.StringFlag{
				Name:    "mode",
				Aliases: []string{"m"},
				Usage:   "Routing mode: driving or straightline",
				Value:   "driving",
			},
			&cli.StringFlag{
				Name:    "units",
				Aliases: []string{"u"},
				Usage:   "Distance units: miles or km",
				Value:   "miles",
			},
		},
		Action: distanceMatrixAction,
	}
}

func distanceMatrixAction(ctx context.Context, cmd *cli.Command) error {
	app, err := newApp(cmd)
	if err != nil {
		return err
	}

	origins, err := readLines(cmd.String("origins"))
	if err != nil {
		return fmt.Errorf("reading origins file: %w", err)
	}

	destinations, err := readLines(cmd.String("destinations"))
	if err != nil {
		return fmt.Errorf("reading destinations file: %w", err)
	}

	if len(origins) == 0 {
		return fmt.Errorf("origins file is empty")
	}

	if len(destinations) == 0 {
		return fmt.Errorf("destinations file is empty")
	}

	resp, err := app.client.DistanceMatrix(ctx, origins, destinations, cmd.String("mode"), cmd.String("units"))
	if err != nil {
		return err
	}

	return app.formatter.FormatDistanceMatrix(resp)
}
