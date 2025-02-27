package tui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

func (m model) View() string {
	if m.currentPage != "auth" {
		return fmt.Sprintf("Main page!\nAuthService: %s\nEmail: %s\nGPG Key: %s\n",
			m.inputs[0].Value(), m.inputs[1].Value(), m.filePath)
	}

	var s string
	for i := range m.inputs {
		if i == m.focused {
			s += activeStyle.Render(">") + m.inputs[i].View() + "\n"
		} else {
			s += inactiveStyle.Render(" ") + m.inputs[i].View() + "\n"
		}
	}

	fileLabel := "Login"
	if m.filePath != "" {
		fileLabel = m.filePath
	}
	if m.focused == len(m.inputs) {
		s += activeStyle.Render(">") + fileLabel + "\n"
	} else {
		s += inactiveStyle.Render(" ") + fileLabel + "\n"
	}

	if m.err != nil {
		s += errorStyle.Render("Error: "+m.err.Error()) + "\n"
	}

	s += "\nPress tab/shift+tab to switch fields\n"

	return lipgloss.NewStyle().Padding(1, 2).Render(s)
}
