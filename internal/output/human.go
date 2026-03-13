package output

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/geocodio/geocodio-cli/internal/api"
)

// Human formats output for human-readable display.
type Human struct {
	w         io.Writer
	useStyles bool
	opts      Options
}

// NewHuman creates a new Human formatter.
func NewHuman(w io.Writer, useStyles bool, opts Options) *Human {
	return &Human{w: w, useStyles: useStyles, opts: opts}
}

// style applies the given style if styling is enabled.
func (h *Human) style(s lipgloss.Style, text string) string {
	if h.useStyles {
		return s.Render(text)
	}
	return text
}

// FormatGeocode formats a geocode response.
func (h *Human) FormatGeocode(resp *api.GeocodeResponse) error {
	if len(resp.Results) == 0 {
		fmt.Fprintln(h.w, h.style(WarningStyle, "No results found"))
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
		fmt.Fprintf(h.w, "%s\n", h.style(HeaderStyle, fmt.Sprintf("Result %d of %d:", num, total)))
	}

	h.printField("Matched Address", r.FormattedAddress)
	h.printField("Coordinates", fmt.Sprintf("%.7f, %.7f", r.Location.Lat, r.Location.Lng))
	h.printField("Accuracy", fmt.Sprintf("%s (%.2f)", r.AccuracyType, r.Accuracy))

	if r.Source != "" {
		h.printField("Source", r.Source)
	}
	if h.opts.ShowAddressKey && r.StableAddressKey != "" {
		h.printField("Address Key", r.StableAddressKey)
	}
	if len(r.Destinations) > 0 {
		fmt.Fprintln(h.w)
		fmt.Fprintf(h.w, "  %s\n", h.style(HeaderStyle, "Distances:"))
		for _, d := range r.Destinations {
			destStr := d.Query
			if d.Geocode != nil {
				destStr = d.Geocode.FormattedAddress
			}
			h.printField("To", destStr)
			h.printField("Distance", fmt.Sprintf("%.1f miles (%.1f km)", d.DistanceMiles, d.DistanceKm))
			if d.DurationSeconds != nil {
				mins := float64(*d.DurationSeconds) / 60
				h.printField("Duration", fmt.Sprintf("%.0f minutes", mins))
			}
		}
	}
}

func (h *Human) printField(label, value string) {
	if h.useStyles {
		fmt.Fprintf(h.w, "  %s %s\n",
			LabelStyle.Width(18).Render(label),
			ValueStyle.Render(value))
	} else {
		fmt.Fprintf(h.w, "  %-18s %s\n", label, value)
	}
}

// FormatBatchGeocode formats a batch geocode response.
func (h *Human) FormatBatchGeocode(resp *api.BatchGeocodeResponse) error {
	for i, r := range resp.Results {
		if i > 0 {
			fmt.Fprintln(h.w)
			fmt.Fprintln(h.w, h.style(DividerStyle, strings.Repeat("-", 60)))
			fmt.Fprintln(h.w)
		}

		fmt.Fprintf(h.w, "%s %s\n\n", h.style(LabelStyle, "Query:"), r.Query)

		if len(r.Response.Results) == 0 {
			fmt.Fprintln(h.w, "  "+h.style(WarningStyle, "No results found"))
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

// FormatDistance formats a distance response.
func (h *Human) FormatDistance(resp *api.DistanceResponse) error {
	if len(resp.Destinations) == 0 {
		fmt.Fprintln(h.w, h.style(WarningStyle, "No results found"))
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

		h.printField("From", originAddr)

		destAddr := d.Query
		if d.Geocode != nil {
			destAddr = d.Geocode.FormattedAddress
		}
		h.printField("To", destAddr)
		h.printField("Distance", fmt.Sprintf("%.1f miles (%.1f km)", d.DistanceMiles, d.DistanceKm))
		if d.DurationSeconds != nil {
			mins := float64(*d.DurationSeconds) / 60
			h.printField("Duration", fmt.Sprintf("%.0f minutes", mins))
		}
	}
	return nil
}

// FormatDistanceMatrix formats a distance matrix response.
func (h *Human) FormatDistanceMatrix(resp *api.DistanceMatrixResponse) error {
	if len(resp.Results) == 0 {
		fmt.Fprintln(h.w, h.style(WarningStyle, "No results found"))
		return nil
	}

	for i, r := range resp.Results {
		if i > 0 {
			fmt.Fprintln(h.w)
			fmt.Fprintln(h.w, h.style(DividerStyle, strings.Repeat("-", 60)))
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
		fmt.Fprintf(h.w, "%s %s\n\n", h.style(LabelStyle, "Origin:"), originStr)

		for _, d := range r.Destinations {
			destStr := d.Query
			if destStr == "" && len(d.Location) == 2 {
				destStr = fmt.Sprintf("%.6f, %.6f", d.Location[0], d.Location[1])
			}
			h.printField("To", destStr)
			h.printField("Distance", fmt.Sprintf("%.1f miles (%.1f km)", d.DistanceMiles, d.DistanceKm))
			if d.DurationSeconds != nil {
				mins := float64(*d.DurationSeconds) / 60
				h.printField("Duration", fmt.Sprintf("%.0f minutes", mins))
			}
			fmt.Fprintln(h.w)
		}
	}
	return nil
}

// FormatDistanceJob formats a distance job response.
func (h *Human) FormatDistanceJob(resp *api.DistanceJobResponse) error {
	if resp.Data == nil {
		fmt.Fprintln(h.w, h.style(WarningStyle, "No job data found"))
		return nil
	}
	j := resp.Data
	h.printField("Job ID", j.Identifier)
	if j.Name != "" {
		h.printField("Name", j.Name)
	}
	h.printStatusField("Status", j.Status)
	if j.Progress > 0 {
		h.printField("Progress", fmt.Sprintf("%d%%", j.Progress))
	}
	if j.StatusMessage != "" {
		h.printField("Message", j.StatusMessage)
	}
	if j.CreatedAt != "" {
		h.printField("Created", j.CreatedAt)
	}
	return nil
}

func (h *Human) printStatusField(label, status string) {
	if h.useStyles {
		fmt.Fprintf(h.w, "  %s %s\n",
			LabelStyle.Width(18).Render(label),
			StatusStyle(status).Render(status))
	} else {
		fmt.Fprintf(h.w, "  %-18s %s\n", label, status)
	}
}

// FormatDistanceJobList formats a list of distance jobs.
func (h *Human) FormatDistanceJobList(resp *api.DistanceJobListResponse) error {
	if len(resp.Jobs) == 0 {
		fmt.Fprintln(h.w, h.style(WarningStyle, "No distance jobs found"))
		return nil
	}

	// Print header
	if h.useStyles {
		fmt.Fprintf(h.w, "%s%s%s%s%s\n",
			TableHeaderStyle.Width(30).Render("ID"),
			TableHeaderStyle.Width(20).Render("Name"),
			TableHeaderStyle.Width(15).Render("Status"),
			TableHeaderStyle.Width(10).Render("Progress"),
			TableHeaderStyle.Width(20).Render("Created"))
	} else {
		fmt.Fprintf(h.w, "%-30s %-20s %-15s %-10s %-20s\n", "ID", "Name", "Status", "Progress", "Created")
	}
	fmt.Fprintln(h.w, h.style(DividerStyle, strings.Repeat("-", 100)))

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
		if h.useStyles {
			fmt.Fprintf(h.w, "%s%s%s%s%s\n",
				ValueStyle.Width(30).Render(id),
				ValueStyle.Width(20).Render(name),
				StatusStyle(j.Status).Width(15).Render(j.Status),
				ValueStyle.Width(10).Render(progress),
				ValueStyle.Width(20).Render(j.CreatedAt))
		} else {
			fmt.Fprintf(h.w, "%-30s %-20s %-15s %-10s %-20s\n", id, name, j.Status, progress, j.CreatedAt)
		}
	}
	return nil
}

// FormatList formats a list response.
func (h *Human) FormatList(resp *api.ListResponse) error {
	h.printField("List ID", fmt.Sprintf("%d", resp.ID))
	if resp.File != nil && resp.File.Filename != "" {
		h.printField("Filename", resp.File.Filename)
	}
	if resp.Status != nil {
		h.printStatusField("Status", resp.Status.State)
		if resp.Status.Progress > 0 {
			h.printField("Progress", fmt.Sprintf("%.1f%%", resp.Status.Progress))
		}
	}
	if resp.Rows != nil {
		if resp.Rows.Total > 0 {
			h.printField("Total Rows", fmt.Sprintf("%d", resp.Rows.Total))
		}
		if resp.Rows.Processed > 0 {
			h.printField("Processed", fmt.Sprintf("%d", resp.Rows.Processed))
		}
		if resp.Rows.Failed > 0 {
			h.printField("Failed", fmt.Sprintf("%d", resp.Rows.Failed))
		}
	}
	if resp.ExpiresAt != "" {
		h.printField("Expires", resp.ExpiresAt)
	}
	return nil
}

// FormatListList formats a list of lists.
func (h *Human) FormatListList(resp *api.ListListResponse) error {
	if len(resp.Lists) == 0 {
		fmt.Fprintln(h.w, h.style(WarningStyle, "No lists found"))
		return nil
	}

	// Print header
	if h.useStyles {
		fmt.Fprintf(h.w, "%s%s%s%s\n",
			TableHeaderStyle.Width(10).Render("ID"),
			TableHeaderStyle.Width(30).Render("Filename"),
			TableHeaderStyle.Width(15).Render("Status"),
			TableHeaderStyle.Width(10).Render("Progress"))
	} else {
		fmt.Fprintf(h.w, "%-10s %-30s %-15s %-10s\n", "ID", "Filename", "Status", "Progress")
	}
	fmt.Fprintln(h.w, h.style(DividerStyle, strings.Repeat("-", 70)))

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
		if h.useStyles {
			fmt.Fprintf(h.w, "%s%s%s%s\n",
				ValueStyle.Width(10).Render(fmt.Sprintf("%d", l.ID)),
				ValueStyle.Width(30).Render(filename),
				StatusStyle(status).Width(15).Render(status),
				ValueStyle.Width(10).Render(progress))
		} else {
			fmt.Fprintf(h.w, "%-10d %-30s %-15s %-10s\n", l.ID, filename, status, progress)
		}
	}
	return nil
}

// FormatError formats an error.
func (h *Human) FormatError(err error) error {
	fmt.Fprintf(h.w, "%s %s\n", h.style(ErrorStyle, "Error:"), err.Error())
	return nil
}

// FormatMessage formats a message.
func (h *Human) FormatMessage(msg string) error {
	fmt.Fprintln(h.w, msg)
	return nil
}
