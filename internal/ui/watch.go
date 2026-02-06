package ui

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type WatchUpdate struct {
	Progress float64
	Status   string
	TimeLeft string
	Done     bool
	Error    error
}

type WatchDisplay struct {
	w         io.Writer
	isTTY     bool
	useStyles bool
	width     int
}

func NewWatchDisplay(w io.Writer) *WatchDisplay {
	f, ok := w.(*os.File)
	isTTY := ok && IsTTY(f)
	useStyles := isTTY && ColorEnabled()

	return &WatchDisplay{
		w:         w,
		isTTY:     isTTY,
		useStyles: useStyles,
		width:     30,
	}
}

func (d *WatchDisplay) Update(update WatchUpdate) {
	if !d.isTTY {
		return
	}

	statusStyle := lipgloss.NewStyle()
	barFillStyle := lipgloss.NewStyle()
	barEmptyStyle := lipgloss.NewStyle()

	if d.useStyles {
		switch update.Status {
		case "COMPLETED":
			statusStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("2")).Bold(true)
			barFillStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
		case "FAILED":
			statusStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("1")).Bold(true)
			barFillStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("1"))
		default:
			statusStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("3"))
			barFillStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("6"))
		}
		barEmptyStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	}

	filled := int(update.Progress / 100.0 * float64(d.width))
	if filled > d.width {
		filled = d.width
	}
	if filled < 0 {
		filled = 0
	}
	empty := d.width - filled

	bar := strings.Repeat("█", filled) + strings.Repeat("░", empty)
	if d.useStyles {
		bar = barFillStyle.Render(strings.Repeat("█", filled)) +
			barEmptyStyle.Render(strings.Repeat("░", empty))
	}

	status := update.Status
	if d.useStyles {
		status = statusStyle.Render(status)
	}

	line := fmt.Sprintf("\r%s [%s] %.1f%%", status, bar, update.Progress)

	if update.TimeLeft != "" {
		line += fmt.Sprintf(" (%s)", update.TimeLeft)
	}

	line += "          "

	fmt.Fprint(d.w, line)
}

func (d *WatchDisplay) Done() {
	if d.isTTY {
		fmt.Fprintln(d.w)
	}
}

func (d *WatchDisplay) IsTTY() bool {
	return d.isTTY
}
