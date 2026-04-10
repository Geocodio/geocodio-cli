package cli

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/geocodio/geocodio-cli/internal/api"
	"github.com/geocodio/geocodio-cli/internal/output"
	"github.com/geocodio/geocodio-cli/internal/ui"
	"github.com/urfave/cli/v3"
)

func geocodeCmd() *cli.Command {
	return &cli.Command{
		Name:                      "geocode",
		Usage:                     "Geocode an address",
		ArgsUsage:                 "[address]",
		DisableSliceFlagSeparator: true,
		Flags: append([]cli.Flag{
			&cli.StringFlag{
				Name:    "batch",
				Aliases: []string{"b"},
				Usage:   "Batch geocode from file (one address per line)",
			},
			&cli.StringFlag{
				Name:    "fields",
				Aliases: []string{"f"},
				Usage:   "Data append fields (comma-separated)",
			},
			&cli.IntFlag{
				Name:    "limit",
				Aliases: []string{"l"},
				Usage:   "Maximum number of results",
			},
			&cli.StringFlag{
				Name:    "country",
				Aliases: []string{"c"},
				Usage:   "Country hint (US or CA)",
			},
			&cli.BoolFlag{
				Name:  "show-address-key",
				Usage: "Show stable address key in output",
			},
		}, destinationFlags()...),
		Action: geocodeAction,
	}
}

func geocodeAction(ctx context.Context, cmd *cli.Command) error {
	app, err := newApp(cmd, output.Options{ShowAddressKey: cmd.Bool("show-address-key"), Units: cmd.String("distance-units")})
	if err != nil {
		return err
	}

	batchFile := cmd.String("batch")
	if batchFile != "" {
		return geocodeBatch(ctx, cmd, app, batchFile)
	}

	if cmd.NArg() < 1 {
		return fmt.Errorf("address required")
	}

	address := strings.Join(cmd.Args().Slice(), " ")
	return geocodeSingle(ctx, cmd, app, address)
}

func geocodeSingle(ctx context.Context, cmd *cli.Command, app *App, address string) error {
	req := &api.GeocodeRequest{
		Address:           address,
		Limit:             int(cmd.Int("limit")),
		Country:           cmd.String("country"),
		DestinationParams: parseDestinationParams(cmd),
	}

	if fields := cmd.String("fields"); fields != "" {
		req.Fields = strings.Split(fields, ",")
	}

	resp, err := app.client.Geocode(ctx, req)
	if err != nil {
		return err
	}

	return app.formatter.FormatGeocode(resp)
}

func geocodeBatch(ctx context.Context, cmd *cli.Command, app *App, filename string) error {
	addresses, err := readLines(filename)
	if err != nil {
		return fmt.Errorf("reading batch file: %w", err)
	}

	if len(addresses) == 0 {
		return fmt.Errorf("batch file is empty")
	}

	if len(addresses) > 10000 {
		return fmt.Errorf("batch size exceeds maximum of 10,000 addresses")
	}

	req := &api.BatchGeocodeRequest{
		Addresses:         addresses,
		Limit:             int(cmd.Int("limit")),
		Country:           cmd.String("country"),
		DestinationParams: parseDestinationParams(cmd),
	}

	if fields := cmd.String("fields"); fields != "" {
		req.Fields = strings.Split(fields, ",")
	}

	resp, err := ui.WithSpinner(app.stderr, "Geocoding addresses...", func() (*api.BatchGeocodeResponse, error) {
		return app.client.BatchGeocode(ctx, req)
	})
	if err != nil {
		return err
	}

	return app.formatter.FormatBatchGeocode(resp)
}

func readLines(filename string) ([]string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			lines = append(lines, line)
		}
	}

	return lines, scanner.Err()
}
