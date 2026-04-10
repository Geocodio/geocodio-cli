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
// Only accepts "USA" or "Canada" (case-insensitive). Other values are ignored.
func appendCountry(address, country string) string {
	if country == "" {
		return address
	}
	normalized := normalizeCountry(country)
	if normalized == "" {
		return address
	}
	if strings.Contains(strings.ToLower(address), strings.ToLower(normalized)) {
		return address
	}
	return address + ", " + normalized
}

// normalizeCountry returns the accepted country string or empty if invalid.
func normalizeCountry(country string) string {
	switch strings.ToLower(country) {
	case "usa":
		return "USA"
	case "canada":
		return "Canada"
	default:
		return ""
	}
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
			&cli.StringFlag{
				Name:    "country",
				Aliases: []string{"c"},
				Usage:   "Country to append to addresses (e.g. Canada)",
			},
		},
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
	if country != "" && normalizeCountry(country) == "" {
		fmt.Fprintf(app.stderr, "Warning: %q is not a valid country. Accepted values are USA or Canada. Flag will be ignored.\n", country)
	}
	origin := appendCountry(args[0], country)
	destinations := appendCountryToAll(args[1:], country)

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
			&cli.StringFlag{
				Name:    "country",
				Aliases: []string{"c"},
				Usage:   "Country to append to addresses (e.g. Canada)",
			},
		},
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
	if country != "" && normalizeCountry(country) == "" {
		fmt.Fprintf(app.stderr, "Warning: %q is not a valid country. Accepted values are USA or Canada. Flag will be ignored.\n", country)
	}
	origins = appendCountryToAll(origins, country)
	destinations = appendCountryToAll(destinations, country)
	mode := cmd.String("mode")
	units := cmd.String("units")

	resp, err := ui.WithSpinner(app.stderr, "Calculating distance matrix...", func() (*api.DistanceMatrixResponse, error) {
		return app.client.DistanceMatrix(ctx, origins, destinations, mode, units)
	})
	if err != nil {
		return err
	}

	return app.formatter.FormatDistanceMatrix(resp)
}
