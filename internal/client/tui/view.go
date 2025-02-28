package tui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

func (m model) View() string {
	if m.login {
		return fmt.Sprintf("Main page!\nServer: %s\nEmail: %s\nGPG Key: %s\n",
			m.inputs[0].Value(), m.profile.Email, m.profile.KeyPath)

	} else {
		return m.authPage()
	}
}

func (m model) authPage() string {
	var s string

	s += headerStyle.Render("GophKeeper")
	s += "\n"

	for i := range m.inputs {
		if i == m.focused {
			s += activeStyle.Render(">") + m.inputs[i].View() + "\n"
		} else {
			s += inactiveStyle.Render(" ") + m.inputs[i].View() + "\n"
		}
	}

	loginLabel := "Go"

	if m.focused == len(m.inputs) {
		s += activeStyle.Render(">") + loginLabel + "\n"
	} else {
		s += inactiveStyle.Render(" ") + loginLabel + "\n"
	}

	if m.err != nil {
		s += errorStyle.Render("Error: "+m.err.Error()) + "\n"
	}

	s += "\nPress tab/shift+tab to switch fields\n"

	return lipgloss.NewStyle().Padding(1, 2).Render(s)
}
