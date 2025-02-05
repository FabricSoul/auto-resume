// internal/models/main_model.go
package models

import (
	"github.com/FabricSoul/auto-resume/internal/types"
	tea "github.com/charmbracelet/bubbletea"
)

type MainModel struct {
	State         types.Appstate
	activeModel   types.Model
	previousState types.Appstate
	projects      *types.ProjectManager

	// Models
	splashScreenModel *SplashModel
	newProjectModel   *NewProjectModel
	errorModel        *ErrorModel
	llmManagerModel   *LLMManagerModel
}

func NewMainModel(pm *types.ProjectManager) *MainModel {
	return &MainModel{
		projects:   pm,
		errorModel: NewErrorModel(),
	}
}

func (m *MainModel) Init() tea.Cmd {
	if m.activeModel == nil {
		m.activeModel = NewSplashModel(m.projects)
		m.State = types.StateSplash
	}
	return m.activeModel.Init()
}

func (m *MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	// Handle error messages first
	switch msg := msg.(type) {
	case types.ErrorMsg:
		m.errorModel.SetError(msg.Error)
		return m, nil
	case tea.KeyMsg:
		// Allow error dismissal from any screen if error is visible
		if m.errorModel.visible {
			if msg.Type == tea.KeyEsc || msg.Type == tea.KeyEnter {
				m.errorModel.visible = false
				return m, nil
			}
			return m, nil
		}
	}

	// Handle state transitions
	switch msg := msg.(type) {
	case types.TransitionMsg:
		m.previousState = m.State
		m.State = msg.To

		switch msg.To {
		case types.StateNewProject:
			if m.newProjectModel == nil {
				m.newProjectModel = NewProjectScreen(m.projects)
			}
			m.activeModel = m.newProjectModel
		case types.StateSplash:
			if m.splashScreenModel == nil {
				m.splashScreenModel = NewSplashModel(m.projects)
			}
			m.activeModel = m.splashScreenModel
		case types.StateLLMManager:
			if m.llmManagerModel == nil {
				m.llmManagerModel = NewLLMManagerModel(m.projects)
			}
			m.activeModel = m.llmManagerModel
		}
	}

	// Update active model
	m.activeModel, cmd = m.activeModel.Update(msg)
	return m, cmd
}

func (m *MainModel) View() string {
	if m.errorModel.visible {
		return m.errorModel.View()
	}
	return m.activeModel.View()
}
