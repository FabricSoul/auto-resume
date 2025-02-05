package models

import (
	"github.com/FabricSoul/auto-resume/internal/types"
	"github.com/FabricSoul/auto-resume/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type NewProjectModel struct {
	project     *types.ProjectManager
	projectName string
	width       int
	height      int
}

func NewProjectScreen(pm *types.ProjectManager) *NewProjectModel {
	return &NewProjectModel{
		project:     pm,
		projectName: "",
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
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, func() tea.Msg {
				return types.TransitionMsg{
					To: types.StateSplash,
				}
			}

		case tea.KeyEnter:
			if m.projectName == "" {
				return m, func() tea.Msg {
					return types.ErrorMsg{
						Error: types.ErrEmptyProjectName,
					}
				}
			}

			if err := m.project.AddProject(m.projectName); err != nil {
				return m, func() tea.Msg {
					return types.ErrorMsg{
						Error: err,
					}
				}
			}

			return m, func() tea.Msg {
				return types.TransitionMsg{
					To: types.StateSplash,
				}
			}

		case tea.KeyBackspace:
			if len(m.projectName) > 0 {
				m.projectName = m.projectName[:len(m.projectName)-1]
			}

		default:
			if msg.Type == tea.KeyRunes {
				m.projectName += string(msg.Runes)
			}
		}
	}

	return m, nil
}

func (m *NewProjectModel) View() string {
	content := ui.Title.Render("Create New Project") + "\n\n"
	content += "Enter project name:\n"
	content += ui.Input.Render(m.projectName + "█")

	help := ui.Help.Render("enter: create • esc: cancel")

	return lipgloss.JoinVertical(
		lipgloss.Left,
		content,
		"\n",
		help,
	)
}
