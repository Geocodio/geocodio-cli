package output

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"

	"github.com/geocodio/geocodio-cli/internal/api"
)

type Agent struct {
	w    io.Writer
	opts Options
}

func NewAgent(w io.Writer, opts Options) *Agent {
	return &Agent{w: w, opts: opts}
}

func (a *Agent) FormatGeocode(resp *api.GeocodeResponse) error {
	if len(resp.Results) == 0 {
		fmt.Fprintln(a.w, "No results found.")
		return nil
	}

	fmt.Fprintln(a.w, "## Geocode Result")
	fmt.Fprintln(a.w)

	for i, r := range resp.Results {
		if len(resp.Results) > 1 {
			fmt.Fprintf(a.w, "### Result %d of %d\n\n", i+1, len(resp.Results))
		}
		a.writeResultTable(&r)
		if i < len(resp.Results)-1 {
			fmt.Fprintln(a.w)
		}
	}
	return nil
}

func (a *Agent) writeResultTable(r *api.GeocodeResult) {
	fmt.Fprintln(a.w, "| Field | Value |")
	fmt.Fprintln(a.w, "|-------|-------|")
	fmt.Fprintf(a.w, "| Matched Address | %s |\n", r.FormattedAddress)
	fmt.Fprintf(a.w, "| Coordinates | %.7f, %.7f |\n", r.Location.Lat, r.Location.Lng)
	fmt.Fprintf(a.w, "| Accuracy | %s (%.2f) |\n", r.AccuracyType, r.Accuracy)
	if r.Source != "" {
		fmt.Fprintf(a.w, "| Source | %s |\n", r.Source)
	}
	if a.opts.ShowAddressKey && r.StableAddressKey != "" {
		fmt.Fprintf(a.w, "| Address Key | %s |\n", r.StableAddressKey)
	}
	if r.Fields != nil && len(*r.Fields) > 0 {
		fmt.Fprintln(a.w)
		fmt.Fprintln(a.w, "**Fields:**")
		fmt.Fprintln(a.w)
		fieldNames := make([]string, 0, len(*r.Fields))
		for name := range *r.Fields {
			fieldNames = append(fieldNames, name)
		}
		sort.Strings(fieldNames)
		for _, name := range fieldNames {
			val := (*r.Fields)[name]
			fmt.Fprintf(a.w, "**%s:**\n\n", name)
			data, err := json.MarshalIndent(val, "", "  ")
			if err != nil {
				fmt.Fprintf(a.w, "%v\n\n", val)
			} else {
				fmt.Fprintf(a.w, "```json\n%s\n```\n\n", string(data))
			}
		}
	}
	if len(r.Destinations) > 0 {
		fmt.Fprintln(a.w)
		fmt.Fprintln(a.w, "**Distances:**")
		fmt.Fprintln(a.w)
		fmt.Fprintln(a.w, "| Destination | Distance | Duration |")
		fmt.Fprintln(a.w, "|-------------|---------:|----------:|")
		for _, d := range r.Destinations {
			destStr := d.Query
			if d.Geocode != nil {
				destStr = d.Geocode.FormattedAddress
			}
			duration := "-"
			if d.DurationSeconds != nil {
				mins := float64(*d.DurationSeconds) / 60
				duration = fmt.Sprintf("%.0f min", mins)
			}
			fmt.Fprintf(a.w, "| %s | %.1f mi / %.1f km | %s |\n",
				destStr, d.DistanceMiles, d.DistanceKm, duration)
		}
	}
}

func (a *Agent) FormatBatchGeocode(resp *api.BatchGeocodeResponse) error {
	fmt.Fprintln(a.w, "## Batch Geocode Results")
	fmt.Fprintln(a.w)

	for i, r := range resp.Results {
		fmt.Fprintf(a.w, "### Query: %s\n\n", r.Query)

		if len(r.Response.Results) == 0 {
			fmt.Fprintln(a.w, "No results found.")
		} else {
			for j, result := range r.Response.Results {
				if len(r.Response.Results) > 1 {
					fmt.Fprintf(a.w, "#### Result %d of %d\n\n", j+1, len(r.Response.Results))
				}
				a.writeResultTable(&result)
			}
		}

		if i < len(resp.Results)-1 {
			fmt.Fprintln(a.w)
			fmt.Fprintln(a.w, "---")
			fmt.Fprintln(a.w)
		}
	}
	return nil
}

func (a *Agent) FormatDistance(resp *api.DistanceResponse) error {
	if len(resp.Destinations) == 0 {
		fmt.Fprintln(a.w, "No results found.")
		return nil
	}

	fmt.Fprintln(a.w, "## Distance Results")
	fmt.Fprintln(a.w)

	originAddr := ""
	if resp.Origin != nil && resp.Origin.Geocode != nil {
		originAddr = resp.Origin.Geocode.FormattedAddress
	} else if resp.Origin != nil {
		originAddr = resp.Origin.Query
	}

	fmt.Fprintln(a.w, "| From | To | Distance | Duration |")
	fmt.Fprintln(a.w, "|------|----|---------:|----------:|")

	for _, d := range resp.Destinations {
		destAddr := d.Query
		if d.Geocode != nil {
			destAddr = d.Geocode.FormattedAddress
		}
		duration := "-"
		if d.DurationSeconds != nil {
			mins := float64(*d.DurationSeconds) / 60
			duration = fmt.Sprintf("%.0f min", mins)
		}
		fmt.Fprintf(a.w, "| %s | %s | %.1f mi / %.1f km | %s |\n",
			originAddr, destAddr, d.DistanceMiles, d.DistanceKm, duration)
	}
	return nil
}

func (a *Agent) FormatDistanceMatrix(resp *api.DistanceMatrixResponse) error {
	if len(resp.Results) == 0 {
		fmt.Fprintln(a.w, "No results found.")
		return nil
	}

	fmt.Fprintln(a.w, "## Distance Matrix Results")
	fmt.Fprintln(a.w)

	for i, r := range resp.Results {
		originStr := ""
		if r.Origin != nil {
			if r.Origin.Query != "" {
				originStr = r.Origin.Query
			} else if len(r.Origin.Location) == 2 {
				originStr = fmt.Sprintf("%.6f, %.6f", r.Origin.Location[0], r.Origin.Location[1])
			}
		}

		fmt.Fprintf(a.w, "### Origin: %s\n\n", originStr)
		fmt.Fprintln(a.w, "| Destination | Distance | Duration |")
		fmt.Fprintln(a.w, "|-------------|---------:|----------:|")

		for _, d := range r.Destinations {
			destStr := d.Query
			if destStr == "" && len(d.Location) == 2 {
				destStr = fmt.Sprintf("%.6f, %.6f", d.Location[0], d.Location[1])
			}
			duration := "-"
			if d.DurationSeconds != nil {
				mins := float64(*d.DurationSeconds) / 60
				duration = fmt.Sprintf("%.0f min", mins)
			}
			fmt.Fprintf(a.w, "| %s | %.1f mi / %.1f km | %s |\n",
				destStr, d.DistanceMiles, d.DistanceKm, duration)
		}

		if i < len(resp.Results)-1 {
			fmt.Fprintln(a.w)
		}
	}
	return nil
}

func (a *Agent) FormatDistanceJob(resp *api.DistanceJobResponse) error {
	if resp.Data == nil {
		fmt.Fprintln(a.w, "No job data found.")
		return nil
	}

	fmt.Fprintln(a.w, "## Distance Job")
	fmt.Fprintln(a.w)
	fmt.Fprintln(a.w, "| Field | Value |")
	fmt.Fprintln(a.w, "|-------|-------|")

	j := resp.Data
	fmt.Fprintf(a.w, "| Job ID | %s |\n", j.Identifier)
	if j.Name != "" {
		fmt.Fprintf(a.w, "| Name | %s |\n", j.Name)
	}
	fmt.Fprintf(a.w, "| Status | %s |\n", j.Status)
	if j.Progress > 0 {
		fmt.Fprintf(a.w, "| Progress | %d%% |\n", j.Progress)
	}
	if j.StatusMessage != "" {
		fmt.Fprintf(a.w, "| Message | %s |\n", j.StatusMessage)
	}
	if j.CreatedAt != "" {
		fmt.Fprintf(a.w, "| Created | %s |\n", j.CreatedAt)
	}
	return nil
}

func (a *Agent) FormatDistanceJobList(resp *api.DistanceJobListResponse) error {
	if len(resp.Jobs) == 0 {
		fmt.Fprintln(a.w, "No distance jobs found.")
		return nil
	}

	fmt.Fprintln(a.w, "## Distance Jobs")
	fmt.Fprintln(a.w)
	fmt.Fprintln(a.w, "| ID | Name | Status | Progress | Created |")
	fmt.Fprintln(a.w, "|----|------|--------|----------|---------|")

	for _, j := range resp.Jobs {
		progress := "-"
		if j.Progress > 0 {
			progress = fmt.Sprintf("%d%%", j.Progress)
		}
		name := j.Name
		if name == "" {
			name = "-"
		}
		fmt.Fprintf(a.w, "| %s | %s | %s | %s | %s |\n",
			j.Identifier, name, j.Status, progress, j.CreatedAt)
	}
	return nil
}

func (a *Agent) FormatList(resp *api.ListResponse) error {
	fmt.Fprintln(a.w, "## List")
	fmt.Fprintln(a.w)
	fmt.Fprintln(a.w, "| Field | Value |")
	fmt.Fprintln(a.w, "|-------|-------|")

	fmt.Fprintf(a.w, "| List ID | %d |\n", resp.ID)
	if resp.File != nil && resp.File.Filename != "" {
		fmt.Fprintf(a.w, "| Filename | %s |\n", resp.File.Filename)
	}
	if resp.Status != nil {
		fmt.Fprintf(a.w, "| Status | %s |\n", resp.Status.State)
		if resp.Status.Progress > 0 {
			fmt.Fprintf(a.w, "| Progress | %.1f%% |\n", resp.Status.Progress)
		}
	}
	if resp.Rows != nil {
		if resp.Rows.Total > 0 {
			fmt.Fprintf(a.w, "| Total Rows | %d |\n", resp.Rows.Total)
		}
		if resp.Rows.Processed > 0 {
			fmt.Fprintf(a.w, "| Processed | %d |\n", resp.Rows.Processed)
		}
		if resp.Rows.Failed > 0 {
			fmt.Fprintf(a.w, "| Failed | %d |\n", resp.Rows.Failed)
		}
	}
	if resp.ExpiresAt != "" {
		fmt.Fprintf(a.w, "| Expires | %s |\n", resp.ExpiresAt)
	}
	return nil
}

func (a *Agent) FormatListList(resp *api.ListListResponse) error {
	if len(resp.Lists) == 0 {
		fmt.Fprintln(a.w, "No lists found.")
		return nil
	}

	fmt.Fprintln(a.w, "## Lists")
	fmt.Fprintln(a.w)
	fmt.Fprintln(a.w, "| ID | Filename | Status | Progress |")
	fmt.Fprintln(a.w, "|----|----------|--------|----------|")

	for _, l := range resp.Lists {
		progress := "-"
		status := "-"
		if l.Status != nil {
			status = l.Status.State
			if l.Status.Progress > 0 {
				progress = fmt.Sprintf("%.1f%%", l.Status.Progress)
			}
		}
		filename := "-"
		if l.File != nil && l.File.Filename != "" {
			filename = l.File.Filename
		}
		fmt.Fprintf(a.w, "| %d | %s | %s | %s |\n", l.ID, filename, status, progress)
	}
	return nil
}

func (a *Agent) FormatError(err error) error {
	fmt.Fprintf(a.w, "**Error:** %s\n", err.Error())
	return nil
}

func (a *Agent) FormatMessage(msg string) error {
	fmt.Fprintln(a.w, msg)
	return nil
}
