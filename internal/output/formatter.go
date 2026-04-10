package output

import (
	"io"

	"github.com/geocodio/geocodio-cli/internal/api"
)

// OutputMode represents the output format mode.
type OutputMode int

const (
	// OutputModeHuman is human-readable formatted output with optional styling.
	OutputModeHuman OutputMode = iota
	// OutputModeJSON is raw JSON output.
	OutputModeJSON
	// OutputModeAgent is markdown-formatted output suitable for LLM consumption.
	OutputModeAgent
)

// Formatter defines the interface for formatting API responses.
type Formatter interface {
	FormatGeocode(resp *api.GeocodeResponse) error
	FormatBatchGeocode(resp *api.BatchGeocodeResponse) error
	FormatDistance(resp *api.DistanceResponse) error
	FormatDistanceMatrix(resp *api.DistanceMatrixResponse) error
	FormatDistanceJob(resp *api.DistanceJobResponse) error
	FormatDistanceJobList(resp *api.DistanceJobListResponse) error
	FormatList(resp *api.ListResponse) error
	FormatListList(resp *api.ListListResponse) error
	FormatError(err error) error
	FormatMessage(msg string) error
}

// Options configures output formatting behavior.
type Options struct {
	ShowAddressKey bool
	Units          string // "miles" or "km"
}

// New creates a new Formatter based on the specified mode and styling preference.
func New(w io.Writer, mode OutputMode, useStyles bool, opts ...Options) Formatter {
	var o Options
	if len(opts) > 0 {
		o = opts[0]
	}

	switch mode {
	case OutputModeJSON:
		return NewJSON(w)
	case OutputModeAgent:
		return NewAgent(w, o)
	default:
		return NewHuman(w, useStyles, o)
	}
}
