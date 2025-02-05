package models

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ErrorModel struct {
	err     error
	width   int
	height  int
	visible bool
}

var (
	errorBox = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#FF0000")).
			Padding(1).
			Width(50)

	errorTitle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000")).
			Bold(true).
			MarginBottom(1)

	errorHelp = lipgloss.NewStyle().
			Foreground(subtle).
			MarginTop(1)
)

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

	content := errorTitle.Render("Error") + "\n"
	content += m.err.Error() + "\n"
	content += errorHelp.Render("press esc or enter to dismiss")

	box := errorBox.Render(content)

	// Center the error box
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		box,
	)
}
