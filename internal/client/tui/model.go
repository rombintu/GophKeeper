package tui

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	activeStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("62")).MarginRight(2)
	inactiveStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240")).MarginRight(2)
	errorStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
)

type model struct {
	currentPage string
	inputs      []textinput.Model
	focused     int
	filePath    string
	err         error
}

func InitialModel() model {
	ti := make([]textinput.Model, 3)
	ti[0] = textinput.New()
	ti[0].Placeholder = "AuthService URL"
	ti[0].Focus()

	ti[1] = textinput.New()
	ti[1].Placeholder = "Email"

	ti[2] = textinput.New()
	ti[2].Placeholder = "GPG Key"

	return model{
		currentPage: "auth",
		inputs:      ti,
		focused:     0,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}
