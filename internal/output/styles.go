package output

import "github.com/charmbracelet/lipgloss"

// Color constants for consistent theming.
var (
	colorSuccess = lipgloss.Color("2")  // Green
	colorWarning = lipgloss.Color("3")  // Yellow
	colorError   = lipgloss.Color("1")  // Red
	colorInfo    = lipgloss.Color("4")  // Blue
	colorMuted   = lipgloss.Color("8")  // Gray
	colorAccent  = lipgloss.Color("6")  // Cyan
)

// LabelStyle is used for field labels.
var LabelStyle = lipgloss.NewStyle().
	Foreground(colorMuted).
	Bold(true)

// ValueStyle is used for field values.
var ValueStyle = lipgloss.NewStyle().
	Foreground(lipgloss.NoColor{})

// SuccessStyle is used for success messages and status.
var SuccessStyle = lipgloss.NewStyle().
	Foreground(colorSuccess).
	Bold(true)

// WarningStyle is used for warning messages and in-progress status.
var WarningStyle = lipgloss.NewStyle().
	Foreground(colorWarning)

// ErrorStyle is used for error messages and failed status.
var ErrorStyle = lipgloss.NewStyle().
	Foreground(colorError).
	Bold(true)

// HeaderStyle is used for section headers.
var HeaderStyle = lipgloss.NewStyle().
	Foreground(colorAccent).
	Bold(true)

// DividerStyle is used for visual separators.
var DividerStyle = lipgloss.NewStyle().
	Foreground(colorMuted)

// TableHeaderStyle is used for table column headers.
var TableHeaderStyle = lipgloss.NewStyle().
	Foreground(colorInfo).
	Bold(true)

// StatusStyle returns the appropriate style based on job status.
func StatusStyle(status string) lipgloss.Style {
	switch status {
	case "COMPLETED":
		return SuccessStyle
	case "PROCESSING", "QUEUED", "PENDING":
		return WarningStyle
	case "FAILED":
		return ErrorStyle
	default:
		return ValueStyle
	}
}
