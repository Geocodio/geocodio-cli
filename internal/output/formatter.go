package output

import (
	"io"

	"github.com/geocodio/geocodio-cli/internal/api"
)

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

func New(w io.Writer, jsonOutput bool) Formatter {
	if jsonOutput {
		return NewJSON(w)
	}
	return NewHuman(w)
}
