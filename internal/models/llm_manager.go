package models

import (
	"errors"
	"fmt"

	"github.com/FabricSoul/auto-resume/internal/types"
	"github.com/FabricSoul/auto-resume/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// LLMManagerModel manages AI models with a two‐column layout and vim-like keys.
type LLMManagerModel struct {
	width, height int

	pm            *types.ProjectManager
	models        []types.AIModel
	selectedIndex int

	editing        bool
	editFieldIndex int // 0: Name, 1: Provider, 2: Model, 3: APIKey, 4: Submit button
	tempModel      types.AIModel
	isNew          bool
}

// NewLLMManagerModel creates a new instance of the model manager.
func NewLLMManagerModel(pm *types.ProjectManager) *LLMManagerModel {
	return &LLMManagerModel{
		pm:            pm,
		models:        pm.GetModels(),
		selectedIndex: 0,
	}
}

func (m *LLMManagerModel) Init() tea.Cmd {
	return nil
}

func (m *LLMManagerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		if m.editing {
			switch msg.String() {
			case "esc":
				m.editing = false
				return m, nil
			case "i":
				if m.editFieldIndex < 4 { // Don't show float input for submit button
					var prompt, initialValue string
					var callback func(string)

					switch m.editFieldIndex {
					case 0:
						prompt = "Enter Model Name"
						initialValue = m.tempModel.Name
						callback = func(value string) {
							m.tempModel.Name = value
						}
					case 1:
						prompt = "Enter Provider"
						initialValue = m.tempModel.Provider
						callback = func(value string) {
							m.tempModel.Provider = value
						}
					case 2:
						prompt = "Enter Model"
						initialValue = m.tempModel.Model
						callback = func(value string) {
							m.tempModel.Model = value
						}
					case 3:
						prompt = "Enter API Key"
						initialValue = m.tempModel.APIKey
						callback = func(value string) {
							m.tempModel.APIKey = value
						}
					}

					return m, func() tea.Msg {
						return types.ShowFloatInputMsg{
							Prompt:       prompt,
							InitialValue: initialValue,
							Callback:     callback,
						}
					}
				}
			case "j", "down":
				if m.editFieldIndex < 4 {
					m.editFieldIndex++
				}
			case "k", "up":
				if m.editFieldIndex > 0 {
					m.editFieldIndex--
				}
			case "enter":
				if m.editFieldIndex == 4 { // Submit button
					if m.tempModel.Name == "" {
						return m, func() tea.Msg {
							return types.ErrorMsg{Error: types.ErrEmptyModelName}
						}
					}
					// Check name uniqueness
					for i, model := range m.models {
						if model.Name == m.tempModel.Name && (m.isNew || i != m.selectedIndex) {
							return m, func() tea.Msg {
								return types.ErrorMsg{Error: errors.New("model name must be unique")}
							}
						}
					}

					// Save model
					if m.isNew {
						m.models = append(m.models, m.tempModel)
						m.selectedIndex = len(m.models) - 1
					} else {
						m.models[m.selectedIndex] = m.tempModel
					}
					if err := m.pm.SaveModels(m.models); err != nil {
						return m, func() tea.Msg {
							return types.ErrorMsg{Error: err}
						}
					}
					m.editing = false
				}
			}
		} else {
			// Non-editing mode key handling.
			switch msg.String() {
			case "q", "ctrl+c":
				// Return to parent model (splash screen) instead of quitting.
				return m, func() tea.Msg {
					return types.TransitionMsg{To: types.StateSplash}
				}
			case "j", "down":
				if m.selectedIndex < len(m.models)-1 {
					m.selectedIndex++
				}
			case "k", "up":
				if m.selectedIndex > 0 {
					m.selectedIndex--
				}
			case "a":
				// Enter creation mode.
				m.editing = true
				m.isNew = true
				m.editFieldIndex = 0
				m.tempModel = types.AIModel{}
			case "e":
				// Enter edit mode for the selected model.
				if len(m.models) > 0 {
					m.editing = true
					m.isNew = false
					m.editFieldIndex = 0
					m.tempModel = m.models[m.selectedIndex]
				}
			}
		}
	}
	return m, nil
}

func (m *LLMManagerModel) renderModelsList() string {
	if len(m.models) == 0 {
		return "No models available"
	}

	var content string
	for i, model := range m.models {
		line := model.Name
		if i == m.selectedIndex {
			line = ui.SelectedItem.Render("► " + line)
		} else {
			line = "  " + line
		}
		content += line + "\n"
	}
	return content
}

func (m *LLMManagerModel) renderDetailsSection() string {
	if !m.editing {
		content := ui.Title.Render("Model Details") + "\n\n"
		if len(m.models) > 0 {
			current := m.models[m.selectedIndex]
			return fmt.Sprintf("%sName: %s\nProvider: %s\nModel: %s\nAPI Key: %s\n",
				content, current.Name, current.Provider, current.Model, current.APIKey)
		}
		return content + "Select a model or press 'a' to add a new one"
	}

	header := "Editing " + map[bool]string{true: "New Model", false: "Model"}[m.isNew]
	content := ui.Title.Render(header) + "\n\n"

	fields := []struct {
		label, value string
	}{
		{"Name", m.tempModel.Name},
		{"Provider", m.tempModel.Provider},
		{"Model", m.tempModel.Model},
		{"API Key", m.tempModel.APIKey},
		{"Submit", ""},
	}

	for i, field := range fields {
		prefix := map[bool]string{true: ui.SelectedItem.Render("► "), false: "  "}[i == m.editFieldIndex]
		if field.label == "Submit" {
			content += fmt.Sprintf("%s[%s]\n", prefix, field.label)
			continue
		}
		content += fmt.Sprintf("%s%s: %s\n", prefix, field.label, field.value)
	}

	return content + "\nPress enter to move between fields or submit, esc to cancel"
}

func (m *LLMManagerModel) View() string {
	listWidth := m.width / 3
	detailsWidth := m.width - listWidth - 4

	listContent := ui.Title.Render("AI Models") + "\n" + m.renderModelsList()
	leftSection := ui.BaseList.Width(listWidth).Height(m.height - 4).Render(listContent)

	detailsContent := m.renderDetailsSection()
	rightSection := ui.BaseDetails.Width(detailsWidth).Height(m.height - 4).Render(detailsContent)

	content := lipgloss.JoinHorizontal(lipgloss.Top, leftSection, rightSection)
	help := ui.Help.Render("a: add • e: edit • j/k: navigate • i: input • enter: submit • esc: cancel • q: back")
	return ui.JoinedContainer.Render(lipgloss.JoinVertical(lipgloss.Left, content, help))
}
