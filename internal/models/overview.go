package models

import (
	"github.com/FabricSoul/auto-resume/internal/types"
	"github.com/FabricSoul/auto-resume/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type SplashModel struct {
	width   int
	height  int
	project *types.ProjectManager
	// Add selected project index for navigation
	selectedIndex int
}

func NewSplashModel(pm *types.ProjectManager) *SplashModel {
	return &SplashModel{
		project: pm,
	}
}

func (m *SplashModel) Init() tea.Cmd {
	return nil
}

func (m *SplashModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "n":
			return m, func() tea.Msg {
				return types.TransitionMsg{
					To: types.StateNewProject,
				}
			}
		case "M":
			return m, func() tea.Msg {
				return types.TransitionMsg{
					To: types.StateLLMManager,
				}
			}
		case "up", "k":
			if m.selectedIndex > 0 {
				m.selectedIndex--
			}
		case "down", "j":
			if m.selectedIndex < len(m.project.Projects)-1 {
				m.selectedIndex++
			}
		case "enter":
			if len(m.project.Projects) > 0 {
				selectedProject := m.project.Projects[m.selectedIndex]
				return m, func() tea.Msg {
					return types.TransitionMsg{
						To:     types.StateProjectOverview,
						Params: selectedProject,
					}
				}
			}
		}
	}
	return m, nil
}

func (m *SplashModel) View() string {
	// Calculate section widths
	listWidth := m.width / 3
	detailsWidth := m.width - listWidth - 4 // Account for margins/padding

	// Projects list section
	projectsList := ui.Title.Render("Projects") + "\n"
	if len(m.project.Projects) == 0 {
		projectsList += "No projects yet"
	} else {
		for i, proj := range m.project.Projects {
			item := proj.Name
			if i == m.selectedIndex {
				item = ui.SelectedItem.Render("► " + item)
			} else {
				item = "  " + item
			}
			projectsList += item + "\n"
		}
	}

	// Details section
	detailsContent := ui.Title.Render("Project Details") + "\n"
	if len(m.project.Projects) > 0 && m.selectedIndex < len(m.project.Projects) {
		selectedProject := m.project.Projects[m.selectedIndex]
		detailsContent += "Name: " + selectedProject.Name + "\n"
		detailsContent += "Created: " + selectedProject.CreatedAt.Format("2006-01-02") + "\n"
		detailsContent += "Last Opened: " + selectedProject.LastOpened.Format("2006-01-02") + "\n"
		detailsContent += "Path: " + selectedProject.Path + "\n"
	} else {
		detailsContent += "Select a project to view details"
	}

	// Help section
	help := ui.Help.Render("n: New Project • M: Manage Models • q: Quit • ↑/↓: Navigate")

	// Layout sections
	leftSection := ui.BaseList.Width(listWidth).Height(m.height - 4).Render(projectsList)
	rightSection := ui.BaseDetails.Width(detailsWidth).Height(m.height - 4).Render(detailsContent)

	// Combine horizontally
	content := lipgloss.JoinHorizontal(lipgloss.Top, leftSection, rightSection)

	// Add help at bottom
	return lipgloss.JoinVertical(lipgloss.Left, content, help)
}
