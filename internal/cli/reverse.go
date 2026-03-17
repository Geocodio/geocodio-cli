package cli

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/geocodio/geocodio-cli/internal/api"
	"github.com/geocodio/geocodio-cli/internal/output"
	"github.com/geocodio/geocodio-cli/internal/ui"
	"github.com/urfave/cli/v3"
)

func reverseCmd() *cli.Command {
	return &cli.Command{
		Name:      "reverse",
		Usage:     "Reverse geocode coordinates",
		ArgsUsage: "<lat,lng>",
		Flags: append([]cli.Flag{
			&cli.StringFlag{
				Name:    "batch",
				Aliases: []string{"b"},
				Usage:   "Batch reverse geocode from file (one coordinate pair per line)",
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
			&cli.BoolFlag{
				Name:  "skip-geocoding",
				Usage: "Skip reverse geocoding and only return field appends for the coordinates",
			},
			&cli.BoolFlag{
				Name:  "show-address-key",
				Usage: "Show stable address key in output",
			},
		}, destinationFlags()...),
		Action: reverseAction,
	}
}

func reverseAction(ctx context.Context, cmd *cli.Command) error {
	app, err := newApp(cmd, output.Options{ShowAddressKey: cmd.Bool("show-address-key")})
	if err != nil {
		return err
	}

	batchFile := cmd.String("batch")
	if batchFile != "" {
		return reverseBatch(ctx, cmd, app, batchFile)
	}

	if cmd.NArg() < 1 {
		return fmt.Errorf("coordinates required (format: lat,lng)")
	}

	return reverseSingle(ctx, cmd, app, cmd.Args().First())
}

func reverseSingle(ctx context.Context, cmd *cli.Command, app *App, coords string) error {
	lat, lng, err := parseCoordinates(coords)
	if err != nil {
		return err
	}

	req := &api.ReverseGeocodeRequest{
		Lat:               lat,
		Lng:               lng,
		Limit:             int(cmd.Int("limit")),
		SkipGeocoding:     cmd.Bool("skip-geocoding"),
		DestinationParams: parseDestinationParams(cmd),
	}

	if fields := cmd.String("fields"); fields != "" {
		req.Fields = strings.Split(fields, ",")
	}

	resp, err := app.client.ReverseGeocode(ctx, req)
	if err != nil {
		return err
	}

	return app.formatter.FormatGeocode(resp)
}

func reverseBatch(ctx context.Context, cmd *cli.Command, app *App, filename string) error {
	lines, err := readLines(filename)
	if err != nil {
		return fmt.Errorf("reading batch file: %w", err)
	}

	if len(lines) == 0 {
		return fmt.Errorf("batch file is empty")
	}

	if len(lines) > 10000 {
		return fmt.Errorf("batch size exceeds maximum of 10,000 coordinates")
	}

	coords := make([]api.Location, 0, len(lines))
	for i, line := range lines {
		lat, lng, err := parseCoordinates(line)
		if err != nil {
			return fmt.Errorf("line %d: %w", i+1, err)
		}
		coords = append(coords, api.Location{Lat: lat, Lng: lng})
	}

	req := &api.BatchReverseGeocodeRequest{
		Coordinates:       coords,
		Limit:             int(cmd.Int("limit")),
		DestinationParams: parseDestinationParams(cmd),
	}

	if fields := cmd.String("fields"); fields != "" {
		req.Fields = strings.Split(fields, ",")
	}

	resp, err := ui.WithSpinner(app.stderr, "Reverse geocoding coordinates...", func() (*api.BatchReverseGeocodeResponse, error) {
		return app.client.BatchReverseGeocode(ctx, req)
	})
	if err != nil {
		return err
	}

	batchResp := &api.BatchGeocodeResponse{
		Results: make([]api.BatchGeocodeResult, len(resp.Results)),
	}
	for i, r := range resp.Results {
		batchResp.Results[i] = api.BatchGeocodeResult(r)
	}

	return app.formatter.FormatBatchGeocode(batchResp)
}

func parseCoordinates(s string) (float64, float64, error) {
	parts := strings.Split(strings.TrimSpace(s), ",")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid coordinate format: expected lat,lng")
	}

	lat, err := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid latitude: %w", err)
	}

	lng, err := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid longitude: %w", err)
	}

	if lat < -90 || lat > 90 {
		return 0, 0, fmt.Errorf("latitude must be between -90 and 90")
	}

	if lng < -180 || lng > 180 {
		return 0, 0, fmt.Errorf("longitude must be between -180 and 180")
	}

	return lat, lng, nil
}
