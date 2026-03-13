package api

// DestinationParams holds parameters for inline distance calculations
// on geocode and reverse endpoints.
type DestinationParams struct {
	Destinations []string
	Mode         string // "driving" or "straightline"
	Units        string // "miles" or "km"
	MaxResults   int
	MaxDistance   float64
	MaxDuration  int
	MinDistance   float64
	MinDuration  int
	OrderBy      string
	SortOrder    string
}

// GeocodeRequest represents a request to the geocode endpoint.
type GeocodeRequest struct {
	Address string
	Fields  []string
	Limit   int
	Country string
	DestinationParams
}

// GeocodeResponse represents a response from the geocode endpoint.
type GeocodeResponse struct {
	Input   *GeocodeInput   `json:"input"`
	Results []GeocodeResult `json:"results"`
}

// GeocodeInput represents the parsed input address.
type GeocodeInput struct {
	AddressComponents *AddressComponents `json:"address_components,omitempty"`
	FormattedAddress  string             `json:"formatted_address,omitempty"`
}

// AddressComponents represents the components of a parsed address.
type AddressComponents struct {
	Number          string `json:"number,omitempty"`
	Predirectional  string `json:"predirectional,omitempty"`
	Prefix          string `json:"prefix,omitempty"`
	Street          string `json:"street,omitempty"`
	Suffix          string `json:"suffix,omitempty"`
	Postdirectional string `json:"postdirectional,omitempty"`
	SecondaryUnit   string `json:"secondaryunit,omitempty"`
	SecondaryNumber string `json:"secondarynumber,omitempty"`
	City            string `json:"city,omitempty"`
	County          string `json:"county,omitempty"`
	State           string `json:"state,omitempty"`
	Zip             string `json:"zip,omitempty"`
	Country         string `json:"country,omitempty"`
}

// GeocodeResult represents a single geocoding result.
type GeocodeResult struct {
	AddressComponents *AddressComponents     `json:"address_components,omitempty"`
	FormattedAddress  string                 `json:"formatted_address"`
	Location          Location               `json:"location"`
	Accuracy          float64                `json:"accuracy"`
	AccuracyType      string                 `json:"accuracy_type"`
	Source            string                 `json:"source,omitempty"`
	StableAddressKey  string                 `json:"stable_address_key,omitempty"`
	Fields            *Fields                `json:"fields,omitempty"`
	Destinations      []DistanceDestination  `json:"destinations,omitempty"`
}

// Location represents geographic coordinates.
type Location struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

// Fields represents optional data appends.
// Uses map for flexibility since the API returns varying structures.
type Fields map[string]interface{}

// BatchGeocodeRequest represents a batch geocoding request.
type BatchGeocodeRequest struct {
	Addresses []string
	Fields    []string
	Limit     int
}

// BatchGeocodeResponse represents a batch geocoding response.
type BatchGeocodeResponse struct {
	Results []BatchGeocodeResult `json:"results"`
}

// BatchGeocodeResult represents a single result in a batch response.
type BatchGeocodeResult struct {
	Query    string          `json:"query"`
	Response GeocodeResponse `json:"response"`
}

// ReverseGeocodeRequest represents a reverse geocode request.
type ReverseGeocodeRequest struct {
	Lat            float64
	Lng            float64
	Fields         []string
	Limit          int
	SkipGeocoding  bool
	DestinationParams
}

// BatchReverseGeocodeRequest represents a batch reverse geocode request.
type BatchReverseGeocodeRequest struct {
	Coordinates []Location
	Fields      []string
	Limit       int
}

// BatchReverseGeocodeResponse represents a batch reverse geocode response.
type BatchReverseGeocodeResponse struct {
	Results []BatchReverseGeocodeResult `json:"results"`
}

// BatchReverseGeocodeResult represents a single result in a batch reverse response.
type BatchReverseGeocodeResult struct {
	Query    string          `json:"query"`
	Response GeocodeResponse `json:"response"`
}

// DistanceRequest represents a distance calculation request.
type DistanceRequest struct {
	Origins      []string
	Destinations []string
	Mode         string // "driving" or "straightline"
	Units        string // "miles" or "km"
}

// DistanceResponse represents a distance calculation response.
type DistanceResponse struct {
	Origin       *DistanceLocation     `json:"origin,omitempty"`
	Mode         string                `json:"mode,omitempty"`
	Destinations []DistanceDestination `json:"destinations,omitempty"`
}

type DistanceLocation struct {
	Query    string         `json:"query,omitempty"`
	Location []float64      `json:"location,omitempty"`
	Geocode  *GeocodeResult `json:"geocode,omitempty"`
}

type DistanceDestination struct {
	Query           string         `json:"query,omitempty"`
	Location        []float64      `json:"location,omitempty"`
	ID              *string        `json:"id,omitempty"`
	DistanceMiles   float64        `json:"distance_miles,omitempty"`
	DistanceKm      float64        `json:"distance_km,omitempty"`
	DurationSeconds *int           `json:"duration_seconds,omitempty"`
	Geocode         *GeocodeResult `json:"geocode,omitempty"`
}

// DistanceMatrixResponse represents a distance-matrix calculation response.
type DistanceMatrixResponse struct {
	Mode    string                 `json:"mode,omitempty"`
	Results []DistanceMatrixResult `json:"results,omitempty"`
}

// DistanceMatrixResult represents a single origin's distances to all destinations.
type DistanceMatrixResult struct {
	Origin       *DistanceLocation     `json:"origin,omitempty"`
	Destinations []DistanceDestination `json:"destinations,omitempty"`
}

// DistanceJobCreateRequest represents a request to create a distance job.
type DistanceJobCreateRequest struct {
	Name         string   `json:"name"`
	Origins      []string `json:"origins"`
	Destinations []string `json:"destinations"`
	Mode         string   `json:"mode,omitempty"`
	Units        string   `json:"units,omitempty"`
}

// DistanceJobResponse represents a distance job response (wrapper).
type DistanceJobResponse struct {
	Data *DistanceJob `json:"data,omitempty"`
}

// DistanceJob represents the actual distance job data.
type DistanceJob struct {
	Identifier     string `json:"identifier"`
	Name           string `json:"name,omitempty"`
	Status         string `json:"status"`
	Progress       int    `json:"progress,omitempty"`
	StatusMessage  string `json:"status_message,omitempty"`
	CreatedAt      string `json:"created_at,omitempty"`
	DownloadURL    string `json:"download_url,omitempty"`
	OriginsCount   int    `json:"origins_count,omitempty"`
	DestsCount     int    `json:"destinations_count,omitempty"`
	TotalCalcs     int    `json:"total_calculations,omitempty"`
	CalcsCompleted int    `json:"calculations_completed,omitempty"`
}

// DistanceJobListResponse represents a list of distance jobs.
type DistanceJobListResponse struct {
	Jobs []DistanceJob `json:"data"`
}

// ListUploadRequest represents a request to upload a list/spreadsheet.
type ListUploadRequest struct {
	Filename  string
	Data      []byte
	Direction string // "forward" or "reverse"
	Format    string // Column format template
	Callback  string // Optional callback URL
}

// ListResponse represents a list/spreadsheet response.
type ListResponse struct {
	ID            int            `json:"id"`
	File          *ListFile      `json:"file,omitempty"`
	Status        *ListStatus    `json:"status,omitempty"`
	Rows          *ListRowCounts `json:"rows,omitempty"`
	ExpiresAt     string         `json:"expires_at,omitempty"`
	DownloadURL   string         `json:"download_url,omitempty"`
	GeocodeFields []string       `json:"fields,omitempty"`
}

type ListFile struct {
	Filename          string   `json:"filename,omitempty"`
	Headers           []string `json:"headers,omitempty"`
	EstimatedRowCount int      `json:"estimated_rows_count,omitempty"`
}

type ListStatus struct {
	State               string  `json:"state,omitempty"`
	Progress            float64 `json:"progress,omitempty"`
	Message             string  `json:"message,omitempty"`
	TimeLeft            *int    `json:"time_left_seconds,omitempty"`
	TimeLeftDescription string  `json:"time_left_description,omitempty"`
}

// ListRowCounts contains row count information for a list.
type ListRowCounts struct {
	Total     int `json:"total,omitempty"`
	Processed int `json:"processed,omitempty"`
	Failed    int `json:"failed,omitempty"`
}

// ListListResponse represents a list of uploaded lists.
type ListListResponse struct {
	Lists      []ListResponse `json:"lists"`
	Page       int            `json:"page,omitempty"`
	TotalPages int            `json:"total_pages,omitempty"`
}
