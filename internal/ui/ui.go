package ui

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

const (
	padding  = 2
	maxWidth = 80
)

var (
	titleStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true)
)

func ShowFiles(files []string) {
	for _, f := range files {
		println(" -", f)
	}
}

func GetExplanation() (string, error) {
	var explanation string
	err := huh.NewInput().
		Title("What did you do in these files?").
		Value(&explanation).
		Run()
	return explanation, err
}

// Simple spinner model for fast mode
type loadingModel struct {
	spinner  spinner.Model
	done     bool
	err      error
	result   string
	callback func() (string, error)
}

func newLoadingModel(callback func() (string, error)) loadingModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return loadingModel{
		spinner:  s,
		callback: callback,
	}
}

func (m loadingModel) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		func() tea.Msg {
			result, err := m.callback()
			if err != nil {
				return errMsg{err}
			}
			return successMsg{result}
		},
	)
}

func (m loadingModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}
	case errMsg:
		m.err = msg.error
		return m, tea.Quit
	case successMsg:
		m.result = msg.data
		m.done = true
		return m, tea.Quit
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m loadingModel) View() string {
	if m.err != nil {
		return fmt.Sprintf("\nError: %v\n", m.err)
	}
	if m.done {
		return ""
	}
	return fmt.Sprintf("\n %s Generating commit message...\n", m.spinner.View())
}

// Fancy progress model for normal mode
type progressModel struct {
	progress progress.Model
	done     bool
	err      error
	result   string
	callback func() (string, error)
	started  bool
}

type tickMsg time.Time

func newProgressModel(callback func() (string, error)) progressModel {
	p := progress.New(
		progress.WithDefaultGradient(),
		progress.WithWidth(maxWidth),
		progress.WithoutPercentage(),
	)

	return progressModel{
		progress: p,
		callback: callback,
	}
}

func (m progressModel) Init() tea.Cmd {
	return tea.Batch(
		tickCmd(),
		func() tea.Msg {
			result, err := m.callback()
			if err != nil {
				return errMsg{err}
			}
			return successMsg{result}
		},
	)
}

func (m progressModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.progress.Width = msg.Width - padding*2 - 4
		if m.progress.Width > maxWidth {
			m.progress.Width = maxWidth
		}
		return m, nil

	case errMsg:
		m.err = msg.error
		return m, tea.Quit

	case successMsg:
		m.result = msg.data
		m.done = true
		cmd := m.progress.SetPercent(1.0)
		return m, tea.Batch(cmd, tea.Quit)

	case tickMsg:
		if !m.started {
			m.started = true
			cmd := m.progress.SetPercent(0.0)
			return m, tea.Batch(tickCmd(), cmd)
		}
		if m.progress.Percent() >= 0.9 {
			return m, nil
		}
		cmd := m.progress.IncrPercent(0.1)
		return m, tea.Batch(tickCmd(), cmd)

	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd
	}

	return m, nil
}

func (m progressModel) View() string {
	if m.err != nil {
		return fmt.Sprintf("\nError: %v\n", m.err)
	}
	if m.done {
		return ""
	}

	pad := strings.Repeat(" ", padding)
	title := titleStyle.Render("Generating commit message...")
	return "\n" +
		pad + title + "\n" +
		pad + m.progress.View() + "\n"
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

type errMsg struct {
	error
}

type successMsg struct {
	data string
}

func WithLoading(callback func() (string, error)) (string, error) {
	return withLoadingInternal(callback, false)
}

func WithFastLoading(callback func() (string, error)) (string, error) {
	return withLoadingInternal(callback, true)
}

func withLoadingInternal(callback func() (string, error), fastMode bool) (string, error) {
	var p *tea.Program
	if fastMode {
		p = tea.NewProgram(newLoadingModel(callback))
	} else {
		p = tea.NewProgram(newProgressModel(callback))
	}

	m, err := p.Run()
	if err != nil {
		return "", fmt.Errorf("failed to run loading UI: %w", err)
	}

	if fastMode {
		model := m.(loadingModel)
		if model.err != nil {
			return "", model.err
		}
		return model.result, nil
	}

	model := m.(progressModel)
	if model.err != nil {
		return "", model.err
	}
	return model.result, nil
}

func ConfirmCommit(message string) (bool, error) {
	fmt.Printf("\nGenerated commit message:\n%s\n\n", message)

	var confirm bool
	form := huh.NewConfirm().
		Title("Use this commit message and commit?").
		Value(&confirm)

	err := form.Run()
	return confirm, err
}

func ShowError(format string, args ...interface{}) {
	fmt.Printf("\nError: "+format+"\n", args...)
	os.Exit(1)
}

func ShowSuccess(message string) {
	fmt.Printf("\n✓ %s\n", message)
}
