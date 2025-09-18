package utils

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// tea model for spinner
type spinnerModel struct {
	sp      spinner.Model
	message string
	stopped bool
}

func (m spinnerModel) Init() tea.Cmd { return spinner.Tick }

func (m spinnerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case tea.KeyMsg:
		return m, nil
	case tea.WindowSizeMsg:
		return m, nil
	}
	var cmd tea.Cmd
	m.sp, cmd = m.sp.Update(msg)
	if m.stopped {
		return m, tea.Quit
	}
	return m, cmd
}

func (m spinnerModel) View() string {
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("#4ECDC4")).Margin(1, 0)
	return style.Render(fmt.Sprintf("%s %s", m.sp.View(), m.message))
}

// ShowLoading displays a Bubble Tea spinner with the given message and returns a stop func
func ShowLoading(message string) func() {
	sp := spinner.New()
	sp.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#4ECDC4"))
	sp.Spinner = spinner.Line

	model := spinnerModel{sp: sp, message: message}
	p := tea.NewProgram(model, tea.WithoutCatchPanics())

	// run program in background and signal when it exits
	done := make(chan struct{})
	go func() { _ = p.Start(); close(done) }()

	// Return function to stop loading
	return func() {
		// Ask program to quit gracefully
		p.Quit()
		// Wait for exit with a timeout; if it hangs, kill
		select {
		case <-done:
			// ok
		case <-time.After(500 * time.Millisecond):
			p.Kill()
		}
	}
}
