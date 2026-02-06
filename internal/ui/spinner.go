package ui

import (
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type spinnerModel struct {
	spinner spinner.Model
	message string
	done    bool
}

func (m spinnerModel) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m spinnerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case stopMsg:
		m.done = true
		return m, tea.Quit
	case updateMsg:
		m.message = string(msg)
		return m, nil
	}
	return m, nil
}

func (m spinnerModel) View() string {
	if m.done {
		return ""
	}
	return fmt.Sprintf("%s %s", m.spinner.View(), m.message)
}

type stopMsg struct{}
type updateMsg string

type Spinner struct {
	program *tea.Program
	mu      sync.Mutex
	running bool
}

func NewSpinner(w io.Writer, message string) *Spinner {
	f, ok := w.(*os.File)
	if !ok || !IsTTY(f) {
		return &Spinner{running: false}
	}

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("6"))

	model := spinnerModel{
		spinner: s,
		message: message,
	}

	opts := []tea.ProgramOption{
		tea.WithOutput(w),
	}

	p := tea.NewProgram(model, opts...)

	return &Spinner{
		program: p,
		running: false,
	}
}

func (s *Spinner) Start() {
	s.mu.Lock()
	if s.program == nil || s.running {
		s.mu.Unlock()
		return
	}
	s.running = true
	s.mu.Unlock()

	go func() {
		_, _ = s.program.Run()
	}()
}

func (s *Spinner) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.program == nil || !s.running {
		return
	}
	s.running = false
	s.program.Send(stopMsg{})
}

func (s *Spinner) UpdateMessage(msg string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.program == nil || !s.running {
		return
	}
	s.program.Send(updateMsg(msg))
}
