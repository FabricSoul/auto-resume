// internal/models/main_model.go
package models

import (
	"github.com/FabricSoul/auto-resume/internal/types"
	"github.com/FabricSoul/auto-resume/internal/ui"
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
	projectModel      *ProjectDetailModel
	floatModel        tea.Model
	showFloat         bool
	isEditing         bool // Global editing state
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
	switch msg := msg.(type) {
	case types.ShowFloatInputMsg:
		m.floatModel = NewFloatInputModel(msg.Prompt, msg.InitialValue, msg.Callback)
		m.showFloat = true
		m.isEditing = true
		return m, nil
	case types.FloatDismissMsg:
		m.showFloat = false
		m.isEditing = false
		return m, nil
	}

	if m.showFloat {
		var cmd tea.Cmd
		m.floatModel, cmd = m.floatModel.Update(msg)
		return m, cmd
	}

	if m.errorModel.visible {
		var cmd tea.Cmd
		newModel, cmd := m.errorModel.Update(msg)
		m.errorModel = newModel.(*ErrorModel)
		return m, cmd
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
		case types.StateProjectOverview:
			if m.projectModel == nil {
				project := msg.Params.(types.Project)
				m.projectModel = NewProjectDetailModel(project.Path, m.projects)
			}
			m.activeModel = m.projectModel
		}
	}

	// Update active model
	var cmd tea.Cmd
	m.activeModel, cmd = m.activeModel.Update(msg)
	return m, cmd
}

func (m *MainModel) View() string {
	var content string
	if m.showFloat {
		content = m.floatModel.View()
	} else if m.errorModel.visible {
		content = m.errorModel.View()
	} else {
		content = m.activeModel.View()
	}
	// Wrap the view with the centralized AppBackground style.
	return ui.AppBackground.Render(content)
}
