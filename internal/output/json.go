package output

import (
	"encoding/json"
	"io"

	"github.com/geocodio/geocodio-cli/internal/api"
)

type JSON struct {
	w   io.Writer
	enc *json.Encoder
}

func NewJSON(w io.Writer) *JSON {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return &JSON{w: w, enc: enc}
}

func (j *JSON) FormatGeocode(resp *api.GeocodeResponse) error {
	return j.enc.Encode(resp)
}

func (j *JSON) FormatBatchGeocode(resp *api.BatchGeocodeResponse) error {
	return j.enc.Encode(resp)
}

func (j *JSON) FormatDistance(resp *api.DistanceResponse) error {
	return j.enc.Encode(resp)
}

func (j *JSON) FormatDistanceMatrix(resp *api.DistanceMatrixResponse) error {
	return j.enc.Encode(resp)
}

func (j *JSON) FormatDistanceJob(resp *api.DistanceJobResponse) error {
	return j.enc.Encode(resp)
}

func (j *JSON) FormatDistanceJobList(resp *api.DistanceJobListResponse) error {
	return j.enc.Encode(resp)
}

func (j *JSON) FormatList(resp *api.ListResponse) error {
	return j.enc.Encode(resp)
}

func (j *JSON) FormatListList(resp *api.ListListResponse) error {
	return j.enc.Encode(resp)
}

func (j *JSON) FormatError(err error) error {
	return j.enc.Encode(map[string]string{"error": err.Error()})
}

func (j *JSON) FormatMessage(msg string) error {
	return j.enc.Encode(map[string]string{"message": msg})
}
