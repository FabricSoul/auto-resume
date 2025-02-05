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
			// Editing mode key handling.
			switch msg.Type {
			case tea.KeyEsc:
				// Cancel editing.
				m.editing = false
				return m, nil
			case tea.KeyBackspace:
				switch m.editFieldIndex {
				case 0:
					if len(m.tempModel.Name) > 0 {
						m.tempModel.Name = m.tempModel.Name[:len(m.tempModel.Name)-1]
					}
				case 1:
					if len(m.tempModel.Provider) > 0 {
						m.tempModel.Provider = m.tempModel.Provider[:len(m.tempModel.Provider)-1]
					}
				case 2:
					if len(m.tempModel.Model) > 0 {
						m.tempModel.Model = m.tempModel.Model[:len(m.tempModel.Model)-1]
					}
				case 3:
					if len(m.tempModel.APIKey) > 0 {
						m.tempModel.APIKey = m.tempModel.APIKey[:len(m.tempModel.APIKey)-1]
					}
				}
				return m, nil
			case tea.KeyEnter:
				if m.editFieldIndex < 4 {
					// Move to the next field.
					m.editFieldIndex++
				} else {
					// Submit pressed.
					// Ensure model name is not empty.
					if m.tempModel.Name == "" {
						return m, func() tea.Msg {
							return types.ErrorMsg{Error: types.ErrEmptyModelName}
						}
					}
					// Check model name uniqueness.
					for i, model := range m.models {
						// When editing, ignore the current model.
						if model.Name == m.tempModel.Name && (m.isNew || i != m.selectedIndex) {
							return m, func() tea.Msg {
								return types.ErrorMsg{Error: errors.New("model name must be unique")}
							}
						}
					}

					// Save or update the model.
					if m.isNew {
						m.models = append(m.models, m.tempModel)
						m.selectedIndex = len(m.models) - 1
					} else {
						m.models[m.selectedIndex] = m.tempModel
					}
					// Save the updated models to config.
					if err := m.pm.SaveModels(m.models); err != nil {
						return m, func() tea.Msg {
							return types.ErrorMsg{Error: err}
						}
					}
					m.editing = false
				}
				return m, nil
			}
			// In editing mode, if not on the Submit button, add runes to the active field.
			if msg.Type == tea.KeyRunes && m.editFieldIndex < 4 {
				char := string(msg.Runes)
				switch m.editFieldIndex {
				case 0:
					m.tempModel.Name += char
				case 1:
					m.tempModel.Provider += char
				case 2:
					m.tempModel.Model += char
				case 3:
					m.tempModel.APIKey += char
				}
			}
			return m, nil
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
				// case "s":
				// 	// Save the current models config outside editing.
				// 	if err := m.pm.SaveModels(m.models); err != nil {
				// 		return m, func() tea.Msg {
				// 			return types.ErrorMsg{Error: err}
				// 		}
				// 	}
			}
		}
	}
	return m, nil
}

func (m *LLMManagerModel) View() string {
	listWidth := m.width / 3
	detailsWidth := m.width - listWidth - 4

	// Left column: list of AI models.
	listContent := ui.Title.Render("AI Models") + "\n"
	if len(m.models) == 0 {
		listContent += "No models available"
	} else {
		for i, model := range m.models {
			line := model.Name
			if i == m.selectedIndex {
				line = ui.SelectedItem.Render("► " + line)
			} else {
				line = "  " + line
			}
			listContent += line + "\n"
		}
	}
	leftSection := ui.BaseList.Width(listWidth).Height(m.height - 4).Render(listContent)

	// Right column: model details or editing form.
	var detailsContent string
	if m.editing {
		header := "Editing "
		if m.isNew {
			header += "New Model"
		} else {
			header += "Model"
		}
		detailsContent += ui.Title.Render(header) + "\n\n"
		// List all fields plus a Submit button.
		fields := []struct {
			label string
			value string
		}{
			{"Name", m.tempModel.Name},
			{"Provider", m.tempModel.Provider},
			{"Model", m.tempModel.Model},
			{"API Key", m.tempModel.APIKey},
			{"Submit", ""},
		}
		for i, field := range fields {
			prefix := "  "
			if i == m.editFieldIndex {
				prefix = ui.SelectedItem.Render("► ")
			}
			// For Submit button, display it as a button-like text.
			if field.label == "Submit" {
				detailsContent += fmt.Sprintf("%s[%s]\n", prefix, field.label)
			} else {
				detailsContent += fmt.Sprintf("%s%s: %s\n", prefix, field.label, field.value)
			}
		}
		detailsContent += "\nPress enter to move between fields or submit, esc to cancel"
	} else {
		detailsContent += ui.Title.Render("Model Details") + "\n\n"
		if len(m.models) > 0 {
			current := m.models[m.selectedIndex]
			detailsContent += fmt.Sprintf("Name: %s\n", current.Name)
			detailsContent += fmt.Sprintf("Provider: %s\n", current.Provider)
			detailsContent += fmt.Sprintf("Model: %s\n", current.Model)
			detailsContent += fmt.Sprintf("API Key: %s\n", current.APIKey)
		} else {
			detailsContent += "Select a model or press 'a' to add a new one"
		}
	}
	rightSection := ui.BaseDetails.Width(detailsWidth).Height(m.height - 4).Render(detailsContent)

	help := ui.Help.Render("a: add • e: edit • j/k: navigate • enter: next/submit • esc: cancel • q: back")
	content := lipgloss.JoinHorizontal(lipgloss.Top, leftSection, rightSection)
	return lipgloss.JoinVertical(lipgloss.Left, content, help)
}
