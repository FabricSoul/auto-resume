package models

import (
	"github.com/FabricSoul/auto-resume/internal/types"
	"github.com/FabricSoul/auto-resume/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type FloatInputModel struct {
	value    string
	prompt   string
	width    int
	height   int
	callback func(string)
}

func NewFloatInputModel(prompt string, initialValue string, callback func(string)) *FloatInputModel {
	return &FloatInputModel{
		value:    initialValue,
		prompt:   prompt,
		callback: callback,
	}
}

func (m *FloatInputModel) Init() tea.Cmd {
	return nil
}

func (m *FloatInputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			return m, func() tea.Msg {
				return types.FloatDismissMsg{}
			}
		case tea.KeyEnter:
			m.callback(m.value)
			return m, func() tea.Msg {
				return types.FloatDismissMsg{}
			}
		case tea.KeyBackspace:
			if len(m.value) > 0 {
				m.value = m.value[:len(m.value)-1]
			}
		default:
			if msg.Type == tea.KeyRunes {
				m.value += string(msg.Runes)
			}
		}
	}
	return m, nil
}

func (m *FloatInputModel) View() string {
	content := ui.Title.Render(m.prompt) + "\n\n"
	content += ui.Input.Render(m.value + "█")
	content += "\n\n" + ui.Help.Render("enter: confirm • esc: cancel")

	box := ui.FloatBox.Render(content)
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		box,
	)
}
