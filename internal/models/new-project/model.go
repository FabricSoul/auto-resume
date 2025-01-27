package newproject

import (
	"fmt"

	"github.com/FabricSoul/auto-resume/internal/types"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

type NewProjectModel struct {
	form *huh.Form
}

func CreateNewProjectModel(pm *types.ProjectManager) *NewProjectModel {
	return &NewProjectModel{
		form: huh.NewForm(
			huh.NewGroup(
				huh.NewInput().Title("Project Name"),
				huh.NewConfirm(),
			),
		),
	}
}

func dirDoesNotExist(s string) bool {
	return true
}

func (m *NewProjectModel) Init() tea.Cmd {
	return nil
}

func (m *NewProjectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f
	}

	return m, cmd
}

func (m *NewProjectModel) View() string {
	if m.form.State == huh.StateCompleted {
		return fmt.Sprintf("Form completed")
	}
	return m.form.View()
}
