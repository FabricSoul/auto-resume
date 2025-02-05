package models

import (
	"github.com/FabricSoul/auto-resume/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ErrorModel struct {
	err     error
	width   int
	height  int
	visible bool
}

func NewErrorModel() *ErrorModel {
	return &ErrorModel{
		visible: false,
	}
}

func (m *ErrorModel) SetError(err error) {
	m.err = err
	m.visible = true
}

func (m *ErrorModel) Init() tea.Cmd {
	return nil
}

func (m *ErrorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		if msg.Type == tea.KeyEsc || msg.Type == tea.KeyEnter {
			m.visible = false
			return m, nil
		}
	}
	return m, nil
}

func (m *ErrorModel) View() string {
	if !m.visible || m.err == nil {
		return ""
	}

	content := ui.ErrorTitle.Render("Error") + "\n"
	content += m.err.Error() + "\n"
	content += ui.Help.Render("press esc or enter to dismiss")

	box := ui.ErrorBox.Render(content)

	// Center the error box
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		box,
	)
}
