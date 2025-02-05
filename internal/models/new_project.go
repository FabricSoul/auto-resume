package models

import (
	"github.com/FabricSoul/auto-resume/internal/types"
	"github.com/FabricSoul/auto-resume/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type NewProjectModel struct {
	projectName string
	width       int
	height      int
	projects    *types.ProjectManager
}

func NewProjectScreen(pm *types.ProjectManager) *NewProjectModel {
	return &NewProjectModel{
		projects: pm,
	}
}

func (m *NewProjectModel) Init() tea.Cmd {
	return nil
}

func (m *NewProjectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "i":
			return m, func() tea.Msg {
				return types.ShowFloatInputMsg{
					Prompt:       "Enter Project Name",
					InitialValue: m.projectName,
					Callback: func(value string) {
						m.projectName = value
					},
				}
			}
		case "esc":
			return m, func() tea.Msg {
				return types.TransitionMsg{To: types.StateSplash}
			}
		case "enter":
			if m.projectName == "" {
				return m, func() tea.Msg {
					return types.ErrorMsg{Error: types.ErrEmptyProjectName}
				}
			}
			if err := m.projects.AddProject(m.projectName); err != nil {
				return m, func() tea.Msg {
					return types.ErrorMsg{Error: err}
				}
			}
			return m, func() tea.Msg {
				return types.TransitionMsg{To: types.StateSplash}
			}
		}
	}
	return m, nil
}

func (m *NewProjectModel) View() string {
	content := ui.Title.Render("Create New Project") + "\n\n"
	content += "Enter project name:\n"
	content += ui.Input.Render(m.projectName)

	help := ui.Help.Render("i: edit • enter: create • esc: cancel")

	mainContent := ui.BaseDetails.Width(m.width / 2).Render(content)
	joined := ui.JoinedContainer.Render(lipgloss.JoinVertical(lipgloss.Left, mainContent, help))
	return joined

}
