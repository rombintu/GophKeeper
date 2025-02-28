package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type errMsg struct{ err error }

func (e errMsg) Error() string { return e.err.Error() }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			// Обработка нажатия Enter на кнопке "Go"
			if s == "enter" && m.focused == len(m.inputs) {
				// Проверка заполнения всех полей
				if m.inputs[0].Value() == "" ||
					m.inputs[1].Value() == "" ||
					m.inputs[2].Value() == "" {
					m.err = fmt.Errorf("fill all fields")
					return m, nil
				}

				m.profile.Email = m.inputs[1].Value()
				m.profile.KeyPath = m.inputs[2].Value()

				// Переход на главную страницу
				m.login = true
				return m, nil
			}

			// Навигация между полями
			if s == "up" || s == "shift+tab" {
				m.focused--
			} else {
				m.focused++
			}

			// Циклическая навигация
			if m.focused > len(m.inputs) {
				m.focused = 0
			} else if m.focused < 0 {
				m.focused = len(m.inputs)
			}

			// Установка фокуса для полей ввода
			cmds := make([]tea.Cmd, len(m.inputs))
			for i := 0; i < len(m.inputs); i++ {
				if i == m.focused {
					cmds[i] = m.inputs[i].Focus()
				} else {
					m.inputs[i].Blur()
				}
			}

			return m, tea.Batch(cmds...)
		}
	case errMsg:
		m.err = msg
		return m, nil
	}

	// Обновление активного поля ввода
	if m.focused < len(m.inputs) {
		var cmd tea.Cmd
		m.inputs[m.focused], cmd = m.inputs[m.focused].Update(msg)
		return m, cmd
	}

	return m, nil
}
