package cli

import (
	"context"
	"fmt"
	"strings"

	"github.com/geocodio/geocodio-cli/internal/api"
	"github.com/geocodio/geocodio-cli/internal/output"
	"github.com/geocodio/geocodio-cli/internal/ui"
	"github.com/urfave/cli/v3"
)

// appendCountry appends the country to an address if it's not already present.
// The value is passed through as-is; the Geocodio API accepts a wide range of
// country names and formats, so no client-side validation or normalization is done.
func appendCountry(address, country string) string {
	if country == "" {
		return address
	}
	if strings.Contains(strings.ToLower(address), strings.ToLower(country)) {
		return address
	}
	return address + ", " + country
}

// appendCountryToAll appends the country to each address in the slice.
func appendCountryToAll(addresses []string, country string) []string {
	if country == "" {
		return addresses
	}
	result := make([]string, len(addresses))
	for i, addr := range addresses {
		result[i] = appendCountry(addr, country)
	}
	return result
}

// distanceFilterFlags returns the filtering and sorting flags shared by the
// distance and distance-matrix commands. They map to the API's max_results,
// max_distance, min_distance, max_duration, min_duration, order_by and
// sort_order parameters.
func distanceFilterFlags() []cli.Flag {
	return []cli.Flag{
		&cli.FloatFlag{
			Name:    "max-distance",
			Aliases: []string{"radius"},
			Usage:   "Only keep destinations within this distance (in --units)",
		},
		&cli.FloatFlag{
			Name:  "min-distance",
			Usage: "Only keep destinations at least this far away (in --units)",
		},
		&cli.IntFlag{
			Name:  "max-duration",
			Usage: "Only keep destinations within this travel time in seconds (driving mode only)",
		},
		&cli.IntFlag{
			Name:  "min-duration",
			Usage: "Only keep destinations at least this many seconds away (driving mode only)",
		},
		&cli.IntFlag{
			Name:  "max-results",
			Usage: "Only keep the N nearest destinations per origin",
		},
		&cli.StringFlag{
			Name:  "order-by",
			Usage: "Sort destinations by: distance or duration",
		},
		&cli.StringFlag{
			Name:  "sort-order",
			Usage: "Sort direction: asc or desc",
		},
	}
}

// parseDistanceFilters extracts the shared distance filters from CLI flags.
func parseDistanceFilters(cmd *cli.Command) api.DistanceFilters {
	return api.DistanceFilters{
		MaxResults:  int(cmd.Int("max-results")),
		MaxDistance: cmd.Float("max-distance"),
		MinDistance: cmd.Float("min-distance"),
		MaxDuration: int(cmd.Int("max-duration")),
		MinDuration: int(cmd.Int("min-duration")),
		OrderBy:     cmd.String("order-by"),
		SortOrder:   cmd.String("sort-order"),
	}
}

func distanceCmd() *cli.Command {
	return &cli.Command{
		Name:      "distance",
		Usage:     "Calculate distance from origin to destinations",
		ArgsUsage: "<origin> <destination> [destination...]",
		Flags: append([]cli.Flag{
			&cli.StringFlag{
				Name:    "mode",
				Aliases: []string{"m"},
				Usage:   "Routing mode: straightline (1 credit per calculation) or driving (2 credits per calculation)",
				Value:   defaultDistanceMode,
			},
			&cli.StringFlag{
				Name:    "units",
				Aliases: []string{"u"},
				Usage:   "Distance units: miles or km",
				Value:   "miles",
			},
			&cli.StringFlag{
				Name:    "country",
				Aliases: []string{"c"},
				Usage:   "Country to append to addresses (e.g. USA, Canada, United Kingdom)",
			},
		}, distanceFilterFlags()...),
		Action: distanceAction,
	}
}

func distanceAction(ctx context.Context, cmd *cli.Command) error {
	app, err := newApp(cmd, output.Options{Units: cmd.String("units")})
	if err != nil {
		return err
	}

	if cmd.NArg() < 2 {
		return fmt.Errorf("requires origin and at least one destination")
	}

	args := cmd.Args().Slice()
	country := cmd.String("country")
	origin := appendCountry(args[0], country)
	destinations := appendCountryToAll(args[1:], country)

	resp, err := app.client.Distance(ctx, &api.DistanceRequest{
		Origin:          origin,
		Destinations:    destinations,
		Mode:            cmd.String("mode"),
		Units:           cmd.String("units"),
		DistanceFilters: parseDistanceFilters(cmd),
	})
	if err != nil {
		return err
	}

	return app.formatter.FormatDistance(resp)
}

func distanceMatrixCmd() *cli.Command {
	return &cli.Command{
		Name:  "distance-matrix",
		Usage: "Calculate distances between multiple origins and destinations",
		Flags: append([]cli.Flag{
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
				Usage:   "Routing mode: straightline (1 credit per calculation) or driving (2 credits per calculation)",
				Value:   defaultDistanceMode,
			},
			&cli.StringFlag{
				Name:    "units",
				Aliases: []string{"u"},
				Usage:   "Distance units: miles or km",
				Value:   "miles",
			},
			&cli.StringFlag{
				Name:    "country",
				Aliases: []string{"c"},
				Usage:   "Country to append to addresses (e.g. USA, Canada, United Kingdom)",
			},
		}, distanceFilterFlags()...),
		Action: distanceMatrixAction,
	}
}

func distanceMatrixAction(ctx context.Context, cmd *cli.Command) error {
	app, err := newApp(cmd, output.Options{Units: cmd.String("units")})
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

	country := cmd.String("country")
	origins = appendCountryToAll(origins, country)
	destinations = appendCountryToAll(destinations, country)

	printCostNotice(app.stderr, len(origins), len(destinations), cmd.String("mode"))

	req := &api.DistanceMatrixRequest{
		Origins:         origins,
		Destinations:    destinations,
		Mode:            cmd.String("mode"),
		Units:           cmd.String("units"),
		DistanceFilters: parseDistanceFilters(cmd),
	}

	resp, err := ui.WithSpinner(app.stderr, "Calculating distance matrix...", func() (*api.DistanceMatrixResponse, error) {
		return app.client.DistanceMatrix(ctx, req)
	})
	if err != nil {
		return err
	}

	return app.formatter.FormatDistanceMatrix(resp)
}
