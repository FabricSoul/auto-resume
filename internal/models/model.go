// internal/models/model.go
package models

import (
	newproject "github.com/FabricSoul/auto-resume/internal/models/new-project"
	"github.com/FabricSoul/auto-resume/internal/models/splash"
	"github.com/FabricSoul/auto-resume/internal/types"
	tea "github.com/charmbracelet/bubbletea"
)

type MainModel struct {
	State         types.Appstate
	activeModel   types.Model
	previousState types.Appstate
	projects      *types.ProjectManager

	// Models
	splashScreenModel *splash.SplashModel
	newProjectModel   *newproject.NewProjectModel
}

func NewMainModel(pm *types.ProjectManager) *MainModel {
	return &MainModel{
		projects: pm,
	}
}

func (m *MainModel) Init() tea.Cmd {
	if m.activeModel == nil {
		m.activeModel = splash.NewSplasheModel(m.projects)
		m.State = types.StateSplash
	}
	return m.activeModel.Init()
}

func (m *MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case types.TransitionMsg:
		m.previousState = m.State
		m.State = msg.To

		switch msg.To {
		case types.StateNewProject:
			if m.newProjectModel == nil {
				m.newProjectModel = newproject.CreateNewProjectModel(m.projects)
			}
			m.activeModel = m.newProjectModel
		}
	}
	var cmd tea.Cmd
	m.activeModel, cmd = m.activeModel.Update(msg)
	return m, cmd
}

func (m *MainModel) View() string {
	return m.activeModel.View()
}
