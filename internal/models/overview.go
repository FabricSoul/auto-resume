package models

import (
	"github.com/FabricSoul/auto-resume/internal/types"
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

var (
	// Style definitions
	subtle    = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
	highlight = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	special   = lipgloss.AdaptiveColor{Light: "#43BF6D", Dark: "#73F59F"}

	list = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(subtle).
		Padding(1).
		MarginRight(2)

	details = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(subtle).
		Padding(1)

	title = lipgloss.NewStyle().
		Foreground(highlight).
		Bold(true).
		MarginLeft(1).
		MarginBottom(1)

	helpStyle = lipgloss.NewStyle().
			Foreground(subtle).
			MarginTop(1)

	selectedItem = lipgloss.NewStyle().
			Foreground(special).
			Bold(true)
)

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
		case "up", "k":
			if m.selectedIndex > 0 {
				m.selectedIndex--
			}
		case "down", "j":
			if m.selectedIndex < len(m.project.Projects)-1 {
				m.selectedIndex++
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
	projectsList := title.Render("Projects") + "\n"
	if len(m.project.Projects) == 0 {
		projectsList += "No projects yet"
	} else {
		for i, proj := range m.project.Projects {
			item := proj.Name
			if i == m.selectedIndex {
				item = selectedItem.Render("► " + item)
			} else {
				item = "  " + item
			}
			projectsList += item + "\n"
		}
	}

	// Details section
	detailsContent := title.Render("Project Details") + "\n"
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
	help := helpStyle.Render("n: New Project • q: Quit • ↑/↓: Navigate")

	// Layout sections
	leftSection := list.Width(listWidth).Height(m.height - 4).Render(projectsList)
	rightSection := details.Width(detailsWidth).Height(m.height - 4).Render(detailsContent)

	// Combine horizontally
	content := lipgloss.JoinHorizontal(lipgloss.Top, leftSection, rightSection)

	// Add help at bottom
	return lipgloss.JoinVertical(lipgloss.Left, content, help)
}
