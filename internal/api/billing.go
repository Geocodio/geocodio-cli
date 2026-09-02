package api

import (
	"net/http"
	"strconv"
)

// Billing holds the per-request billing counters the Geocodio API returns as
// response headers. Surfacing them lets callers reconcile what they were
// actually billed against what they estimated before submitting the request.
type Billing struct {
	Lookups              *int `json:"lookups,omitempty"`
	DistanceCalculations *int `json:"distance_calculations,omitempty"`
}

// billingSetter is implemented by response types that carry billing counters.
// The client fills them in from the response headers after decoding.
type billingSetter interface {
	setBilling(*Billing)
}

// parseBilling reads the billing headers from a response. It returns nil when
// the response carries none of them, so the field stays omitted from output.
func parseBilling(h http.Header) *Billing {
	lookups := headerInt(h, "X-Billable-Lookups-Count")
	calculations := headerInt(h, "X-Billable-Distance-Calculations")

	if lookups == nil && calculations == nil {
		return nil
	}

	return &Billing{
		Lookups:              lookups,
		DistanceCalculations: calculations,
	}
}

// headerInt parses an integer header value, returning nil if it is absent or
// not a number.
func headerInt(h http.Header, name string) *int {
	raw := h.Get(name)
	if raw == "" {
		return nil
	}
	n, err := strconv.Atoi(raw)
	if err != nil {
		return nil
	}
	return &n
}

func (r *GeocodeResponse) setBilling(b *Billing)             { r.Billing = b }
func (r *BatchGeocodeResponse) setBilling(b *Billing)        { r.Billing = b }
func (r *BatchReverseGeocodeResponse) setBilling(b *Billing) { r.Billing = b }
func (r *DistanceResponse) setBilling(b *Billing)            { r.Billing = b }
func (r *DistanceMatrixResponse) setBilling(b *Billing)      { r.Billing = b }
