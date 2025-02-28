package tui

import "github.com/charmbracelet/lipgloss"

var (
	headerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("12")).Align(lipgloss.Center).
			Bold(true).
			Foreground(lipgloss.Color("12")).
			MarginBottom(1).
			Width(20)
	activeStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("62")).MarginRight(2)
	inactiveStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240")).MarginRight(2)
	errorStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
)
