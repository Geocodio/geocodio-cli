package cli

import (
	"strings"

	"github.com/geocodio/geocodio-cli/internal/api"
	"github.com/urfave/cli/v3"
)

// destinationFlags returns the shared flags for inline distance calculations
// on geocode and reverse commands.
func destinationFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringSliceFlag{
			Name:    "destinations",
			Aliases: []string{"d"},
			Usage:   "Destination addresses or coordinates for distance calculation",
		},
		&cli.StringFlag{
			Name:    "distance-mode",
			Aliases: []string{"m"},
			Usage:   "Distance mode: driving or straightline",
		},
		&cli.StringFlag{
			Name:    "distance-units",
			Aliases: []string{"u"},
			Usage:   "Distance units: miles or km",
		},
		&cli.IntFlag{
			Name:  "distance-max-results",
			Usage: "Maximum number of destinations to return per result",
		},
		&cli.FloatFlag{
			Name:  "distance-max-distance",
			Usage: "Maximum distance filter (in specified units)",
		},
		&cli.IntFlag{
			Name:  "distance-max-duration",
			Usage: "Maximum duration filter in seconds (driving mode only)",
		},
		&cli.FloatFlag{
			Name:  "distance-min-distance",
			Usage: "Minimum distance filter (in specified units)",
		},
		&cli.IntFlag{
			Name:  "distance-min-duration",
			Usage: "Minimum duration filter in seconds (driving mode only)",
		},
		&cli.StringFlag{
			Name:  "distance-order-by",
			Usage: "Field to sort destinations by",
		},
		&cli.StringFlag{
			Name:  "distance-sort-order",
			Usage: "Sort order for destinations",
		},
	}
}

// parseDestinationParams extracts destination parameters from CLI flags.
func parseDestinationParams(cmd *cli.Command) api.DestinationParams {
	dests := cmd.StringSlice("destinations")

	// Also support comma-separated destinations in a single flag value
	var expanded []string
	for _, d := range dests {
		if strings.Contains(d, ";") {
			expanded = append(expanded, strings.Split(d, ";")...)
		} else {
			expanded = append(expanded, d)
		}
	}

	return api.DestinationParams{
		Destinations: expanded,
		Mode:         cmd.String("distance-mode"),
		Units:        cmd.String("distance-units"),
		MaxResults:   int(cmd.Int("distance-max-results")),
		MaxDistance:  cmd.Float("distance-max-distance"),
		MaxDuration:  int(cmd.Int("distance-max-duration")),
		MinDistance:  cmd.Float("distance-min-distance"),
		MinDuration:  int(cmd.Int("distance-min-duration")),
		OrderBy:      cmd.String("distance-order-by"),
		SortOrder:    cmd.String("distance-sort-order"),
	}
}
