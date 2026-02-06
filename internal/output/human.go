package output

import (
	"fmt"
	"io"
	"strings"

	"github.com/geocodio/geocodio-cli/internal/api"
)

type Human struct {
	w io.Writer
}

func NewHuman(w io.Writer) *Human {
	return &Human{w: w}
}

func (h *Human) FormatGeocode(resp *api.GeocodeResponse) error {
	if len(resp.Results) == 0 {
		fmt.Fprintln(h.w, "No results found")
		return nil
	}

	for i, r := range resp.Results {
		if i > 0 {
			fmt.Fprintln(h.w)
		}
		h.formatResult(&r, i+1, len(resp.Results))
	}
	return nil
}

func (h *Human) formatResult(r *api.GeocodeResult, num, total int) {
	if total > 1 {
		fmt.Fprintf(h.w, "Result %d of %d:\n", num, total)
	}

	fmt.Fprintf(h.w, "  %-18s %s\n", "Matched Address", r.FormattedAddress)
	fmt.Fprintf(h.w, "  %-18s %.7f, %.7f\n", "Coordinates", r.Location.Lat, r.Location.Lng)
	fmt.Fprintf(h.w, "  %-18s %s (%.2f)\n", "Accuracy", r.AccuracyType, r.Accuracy)

	if r.Source != "" {
		fmt.Fprintf(h.w, "  %-18s %s\n", "Source", r.Source)
	}
}

func (h *Human) FormatBatchGeocode(resp *api.BatchGeocodeResponse) error {
	for i, r := range resp.Results {
		if i > 0 {
			fmt.Fprintln(h.w)
			fmt.Fprintln(h.w, strings.Repeat("-", 60))
			fmt.Fprintln(h.w)
		}

		fmt.Fprintf(h.w, "Query: %s\n\n", r.Query)

		if len(r.Response.Results) == 0 {
			fmt.Fprintln(h.w, "  No results found")
			continue
		}

		for j, result := range r.Response.Results {
			if j > 0 {
				fmt.Fprintln(h.w)
			}
			h.formatResult(&result, j+1, len(r.Response.Results))
		}
	}
	return nil
}

func (h *Human) FormatDistance(resp *api.DistanceResponse) error {
	if len(resp.Destinations) == 0 {
		fmt.Fprintln(h.w, "No results found")
		return nil
	}

	originAddr := ""
	if resp.Origin != nil && resp.Origin.Geocode != nil {
		originAddr = resp.Origin.Geocode.FormattedAddress
	} else if resp.Origin != nil {
		originAddr = resp.Origin.Query
	}

	for i, d := range resp.Destinations {
		if i > 0 {
			fmt.Fprintln(h.w)
		}

		fmt.Fprintf(h.w, "  %-18s %s\n", "From", originAddr)

		destAddr := d.Query
		if d.Geocode != nil {
			destAddr = d.Geocode.FormattedAddress
		}
		fmt.Fprintf(h.w, "  %-18s %s\n", "To", destAddr)
		fmt.Fprintf(h.w, "  %-18s %.1f miles (%.1f km)\n", "Distance", d.DistanceMiles, d.DistanceKm)
		if d.DurationSeconds != nil {
			mins := float64(*d.DurationSeconds) / 60
			fmt.Fprintf(h.w, "  %-18s %.0f minutes\n", "Duration", mins)
		}
	}
	return nil
}

func (h *Human) FormatDistanceMatrix(resp *api.DistanceMatrixResponse) error {
	if len(resp.Results) == 0 {
		fmt.Fprintln(h.w, "No results found")
		return nil
	}

	for i, r := range resp.Results {
		if i > 0 {
			fmt.Fprintln(h.w)
			fmt.Fprintln(h.w, strings.Repeat("-", 60))
			fmt.Fprintln(h.w)
		}

		originStr := ""
		if r.Origin != nil {
			if r.Origin.Query != "" {
				originStr = r.Origin.Query
			} else if len(r.Origin.Location) == 2 {
				originStr = fmt.Sprintf("%.6f, %.6f", r.Origin.Location[0], r.Origin.Location[1])
			}
		}
		fmt.Fprintf(h.w, "Origin: %s\n\n", originStr)

		for _, d := range r.Destinations {
			destStr := d.Query
			if destStr == "" && len(d.Location) == 2 {
				destStr = fmt.Sprintf("%.6f, %.6f", d.Location[0], d.Location[1])
			}
			fmt.Fprintf(h.w, "  %-18s %s\n", "To", destStr)
			fmt.Fprintf(h.w, "  %-18s %.1f miles (%.1f km)\n", "Distance", d.DistanceMiles, d.DistanceKm)
			if d.DurationSeconds != nil {
				mins := float64(*d.DurationSeconds) / 60
				fmt.Fprintf(h.w, "  %-18s %.0f minutes\n", "Duration", mins)
			}
			fmt.Fprintln(h.w)
		}
	}
	return nil
}

func (h *Human) FormatDistanceJob(resp *api.DistanceJobResponse) error {
	if resp.Data == nil {
		fmt.Fprintln(h.w, "No job data found")
		return nil
	}
	j := resp.Data
	fmt.Fprintf(h.w, "  %-18s %s\n", "Job ID", j.Identifier)
	if j.Name != "" {
		fmt.Fprintf(h.w, "  %-18s %s\n", "Name", j.Name)
	}
	fmt.Fprintf(h.w, "  %-18s %s\n", "Status", j.Status)
	if j.Progress > 0 {
		fmt.Fprintf(h.w, "  %-18s %d%%\n", "Progress", j.Progress)
	}
	if j.StatusMessage != "" {
		fmt.Fprintf(h.w, "  %-18s %s\n", "Message", j.StatusMessage)
	}
	if j.CreatedAt != "" {
		fmt.Fprintf(h.w, "  %-18s %s\n", "Created", j.CreatedAt)
	}
	return nil
}

func (h *Human) FormatDistanceJobList(resp *api.DistanceJobListResponse) error {
	if len(resp.Jobs) == 0 {
		fmt.Fprintln(h.w, "No distance jobs found")
		return nil
	}

	fmt.Fprintf(h.w, "%-30s %-20s %-15s %-10s %-20s\n", "ID", "Name", "Status", "Progress", "Created")
	fmt.Fprintln(h.w, strings.Repeat("-", 100))

	for _, j := range resp.Jobs {
		progress := ""
		if j.Progress > 0 {
			progress = fmt.Sprintf("%d%%", j.Progress)
		}
		id := j.Identifier
		if len(id) > 28 {
			id = id[:25] + "..."
		}
		name := j.Name
		if len(name) > 18 {
			name = name[:15] + "..."
		}
		fmt.Fprintf(h.w, "%-30s %-20s %-15s %-10s %-20s\n", id, name, j.Status, progress, j.CreatedAt)
	}
	return nil
}

func (h *Human) FormatList(resp *api.ListResponse) error {
	fmt.Fprintf(h.w, "  %-18s %d\n", "List ID", resp.ID)
	if resp.File != nil && resp.File.Filename != "" {
		fmt.Fprintf(h.w, "  %-18s %s\n", "Filename", resp.File.Filename)
	}
	if resp.Status != nil {
		fmt.Fprintf(h.w, "  %-18s %s\n", "Status", resp.Status.State)
		if resp.Status.Progress > 0 {
			fmt.Fprintf(h.w, "  %-18s %.1f%%\n", "Progress", resp.Status.Progress)
		}
	}
	if resp.Rows != nil {
		if resp.Rows.Total > 0 {
			fmt.Fprintf(h.w, "  %-18s %d\n", "Total Rows", resp.Rows.Total)
		}
		if resp.Rows.Processed > 0 {
			fmt.Fprintf(h.w, "  %-18s %d\n", "Processed", resp.Rows.Processed)
		}
		if resp.Rows.Failed > 0 {
			fmt.Fprintf(h.w, "  %-18s %d\n", "Failed", resp.Rows.Failed)
		}
	}
	if resp.ExpiresAt != "" {
		fmt.Fprintf(h.w, "  %-18s %s\n", "Expires", resp.ExpiresAt)
	}
	return nil
}

func (h *Human) FormatListList(resp *api.ListListResponse) error {
	if len(resp.Lists) == 0 {
		fmt.Fprintln(h.w, "No lists found")
		return nil
	}

	fmt.Fprintf(h.w, "%-10s %-30s %-15s %-10s\n", "ID", "Filename", "Status", "Progress")
	fmt.Fprintln(h.w, strings.Repeat("-", 70))

	for _, l := range resp.Lists {
		progress := ""
		status := ""
		if l.Status != nil {
			status = l.Status.State
			if l.Status.Progress > 0 {
				progress = fmt.Sprintf("%.1f%%", l.Status.Progress)
			}
		}
		filename := ""
		if l.File != nil {
			filename = l.File.Filename
		}
		if len(filename) > 28 {
			filename = filename[:25] + "..."
		}
		fmt.Fprintf(h.w, "%-10d %-30s %-15s %-10s\n", l.ID, filename, status, progress)
	}
	return nil
}

func (h *Human) FormatError(err error) error {
	fmt.Fprintf(h.w, "Error: %s\n", err.Error())
	return nil
}

func (h *Human) FormatMessage(msg string) error {
	fmt.Fprintln(h.w, msg)
	return nil
}
