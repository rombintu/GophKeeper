package tui

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Profile struct {
	Server  string
	Email   string
	KeyPath string
	Token   string
}

type model struct {
	login       bool
	currentPage string
	inputs      []textinput.Model
	focused     int
	profile     Profile
	err         error
}

func InitialModel() model {
	ti := make([]textinput.Model, 3)
	ti[0] = textinput.New()
	ti[0].Placeholder = "Server"
	ti[0].Focus()

	ti[1] = textinput.New()
	ti[1].Placeholder = "Email"

	ti[2] = textinput.New()
	ti[2].Placeholder = "Path to GPG Key"

	return model{
		currentPage: "auth",
		inputs:      ti,
		focused:     0,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}
