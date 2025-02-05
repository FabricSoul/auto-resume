package models

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/FabricSoul/auto-resume/internal/types"
	"github.com/FabricSoul/auto-resume/internal/ui"
)

const (
	FocusOverview = iota
	FocusOutputs
	FocusJob
)

const (
	OverviewFieldProjectName = iota
	OverviewFieldResumeInput
	OverviewFieldLLM
)

const (
	JobFieldName = iota
	JobFieldDescription
	JobFieldGenerate
	JobFieldSavePDF
)

// ProjectDetailModel represents a detailed project configuration screen.
// The left half is split vertically into an overview form (upper part)
// and an outputs list (lower part). The right half shows job-specific fields.
type ProjectDetailModel struct {
	width, height int

	// Which screen region currently has focus:
	// 0 = overview form, 1 = outputs list, 2 = job-specific fields.
	focusArea int

	// Within the overview form, which field is being edited:
	// 0 = project name, 1 = resume input, 2 = LLM selection.
	overviewField int

	// Within the job-specific section:
	// 0 = job description input, 1 = generate button, 2 = save-to-PDF button.
	jobField int

	// Overview section fields.
	overviewProjectName string
	resumeInput         string

	// LLM options drawn from the application config.
	llmOptions []types.AIModel

	// Outputs section: list of outputs and selection index.
	outputs             []types.Output
	selectedOutputIndex int

	// The project directory where the project config (project.toml) is saved.
	projectDir string

	showLLMSelector  bool
	selectedLLMIndex int
	llmList          []types.AIModel
	selectedLLM      types.AIModel
	projects         *types.ProjectManager
}

// NewProjectDetailModel constructs and initializes the project detail model.
func NewProjectDetailModel(projectDir string, pm *types.ProjectManager) *ProjectDetailModel {
	// Load existing project config if available.
	config, err := types.LoadProjectConfig(projectDir)
	if err != nil {
		// Use default values if the config isn't present or cannot be parsed.
		config = types.ProjectConfig{
			Name:        "",
			Model:       "",
			ResumeInput: "",
			Outputs:     []types.Output{},
		}
	}
	llmOptions := pm.GetModels()
	selectedLLMIndex := 0
	for i, m := range llmOptions {
		if m.Name == config.Model {
			selectedLLMIndex = i
			break
		}
	}
	return &ProjectDetailModel{
		overviewProjectName: config.Name,
		resumeInput:         config.ResumeInput,
		outputs:             config.Outputs,
		llmOptions:          llmOptions,
		selectedLLMIndex:    selectedLLMIndex,
		projectDir:          projectDir,
		focusArea:           FocusOverview,
		overviewField:       OverviewFieldProjectName,
		jobField:            JobFieldName,
		projects:            pm,
	}
}

func (m *ProjectDetailModel) Init() tea.Cmd {
	return nil
}

func (m *ProjectDetailModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		if m.showLLMSelector {
			switch msg.String() {
			case "esc":
				m.showLLMSelector = false
			case "j", "down":
				if m.selectedLLMIndex < len(m.llmList)-1 {
					m.selectedLLMIndex++
				}
			case "k", "up":
				if m.selectedLLMIndex > 0 {
					m.selectedLLMIndex--
				}
			case "enter":
				// Set selected LLM and hide selector
				if len(m.llmList) > 0 {
					m.selectedLLM = m.llmList[m.selectedLLMIndex]
				}
				m.showLLMSelector = false
			}
			return m, nil
		}

		switch msg.String() {
		case "q", "ctrl+c":
			return m, func() tea.Msg {
				return types.TransitionMsg{To: types.StateSplash}
			}
		case "tab":
			m.focusArea = (m.focusArea + 1) % 3
		case "shift+tab":
			m.focusArea = (m.focusArea + 2) % 3
		case "l", "enter":
			if m.focusArea == FocusOverview && m.overviewField == OverviewFieldLLM {
				m.showLLMSelector = true
				m.llmList = m.projects.GetModels()
				return m, nil
			}
		case "j", "down":
			switch m.focusArea {
			case FocusOverview:
				if m.overviewField < OverviewFieldLLM {
					m.overviewField++
				}
			case FocusOutputs:
				if m.selectedOutputIndex < len(m.outputs)-1 {
					m.selectedOutputIndex++
				}
			case FocusJob:
				if m.jobField < JobFieldSavePDF {
					m.jobField++
				}
			}
		case "k", "up":
			switch m.focusArea {
			case FocusOverview:
				if m.overviewField > 0 {
					m.overviewField--
				}
			case FocusOutputs:
				if m.selectedOutputIndex > 0 {
					m.selectedOutputIndex--
				}
			case FocusJob:
				if m.jobField > JobFieldName {
					m.jobField--
				}
			}
		case "ctrl+s":
			return m, m.saveProjectConfig
		}

		if m.focusArea == FocusJob {
			if len(m.outputs) == 0 {
				break
			}
			current := &m.outputs[m.selectedOutputIndex]

			switch msg.String() {
			case "i":
				var prompt, initialValue string
				var callback func(string)

				switch m.jobField {
				case JobFieldName:
					prompt = "Enter Job Name"
					initialValue = current.Name
					callback = func(value string) {
						current.Name = value
					}
				case JobFieldDescription:
					prompt = "Enter Job Description"
					initialValue = current.JobDescription
					callback = func(value string) {
						current.JobDescription = value
					}
				case JobFieldGenerate:
					// No input needed for generate field
					break
				case JobFieldSavePDF:
					// No input needed for save PDF field
					break
				}

				if callback != nil {
					return m, func() tea.Msg {
						return types.ShowFloatInputMsg{
							Prompt:       prompt,
							InitialValue: initialValue,
							Callback:     callback,
						}
					}
				}

			case "enter":
				switch m.jobField {
				case JobFieldGenerate:
					if m.selectedLLMIndex < len(m.llmOptions) {
						current.GeneratedOutput = m.resumeInput + "\nGenerated for job: " + current.JobDescription +
							"\n[Using model: " + m.llmOptions[m.selectedLLMIndex].Name + "]"
					} else {
						return m, func() tea.Msg {
							return types.ErrorMsg{Error: fmt.Errorf("no LLM selected")}
						}
					}
				case JobFieldSavePDF:
					return m, m.saveCurrentOutputToPDF
				}
			}
		}
	}
	return m, nil
}

func (m *ProjectDetailModel) View() string {
	if m.showLLMSelector {
		return m.renderLLMSelector()
	}
	// Divide the screen into two columns.
	leftWidth := m.width / 2
	rightWidth := m.width - leftWidth - 2

	// The left column is split vertically into the overview section (upper) and outputs section (lower).
	overviewView := m.renderOverviewSection()
	outputsView := m.renderOutputsSection()
	leftSection := lipgloss.JoinVertical(lipgloss.Left, overviewView, outputsView)
	leftSection = ui.BaseList.Width(leftWidth).Render(leftSection)

	// The right column renders the job-specific fields.
	jobView := m.renderJobSection()
	rightSection := ui.BaseDetails.Width(rightWidth).Render(jobView)

	mainView := lipgloss.JoinHorizontal(lipgloss.Top, leftSection, rightSection)
	help := ui.Help.Render("tab: switch section • j/k: navigate • i: input • enter: action • ctrl+s: save")
	joined := ui.JoinedContainer.Render(lipgloss.JoinVertical(lipgloss.Left, mainView, help))
	return joined
}

func (m *ProjectDetailModel) renderOverviewSection() string {
	title := ui.Title.Render("Project Overview")
	nameField := "Project Name: " + m.overviewProjectName
	resumeField := "Resume Input: " + m.resumeInput
	llmField := "LLM: "
	if len(m.llmOptions) > 0 {
		llmField += m.llmOptions[m.selectedLLMIndex].Name
	} else {
		llmField += "None"
	}

	// Highlight the active field if the overview section has focus.
	if m.focusArea == FocusOverview {
		switch m.overviewField {
		case OverviewFieldProjectName:
			nameField = ui.SelectedItem.Render("► " + nameField)
		case OverviewFieldResumeInput:
			resumeField = ui.SelectedItem.Render("► " + resumeField)
		case OverviewFieldLLM:
			llmField = ui.SelectedItem.Render("► " + llmField)
		}
	}

	fields := []string{nameField, resumeField, llmField}
	return title + "\n" + strings.Join(fields, "\n")
}

func (m *ProjectDetailModel) renderOutputsSection() string {
	title := ui.Title.Render("Outputs")
	var outputLines []string
	if len(m.outputs) == 0 {
		outputLines = append(outputLines, "No outputs. Press 'a' to add.")
	} else {
		for i, out := range m.outputs {
			line := out.Name
			if i == m.selectedOutputIndex && m.focusArea == FocusOutputs {
				line = ui.SelectedItem.Render("► " + line)
			} else {
				line = "  " + line
			}
			outputLines = append(outputLines, line)
		}
	}
	return title + "\n" + strings.Join(outputLines, "\n")
}

func (m *ProjectDetailModel) renderJobSection() string {
	title := ui.Title.Render("Job Specific Fields")
	var currentOutput types.Output
	if len(m.outputs) > 0 {
		currentOutput = m.outputs[m.selectedOutputIndex]
	}
	nameField := "Output Name: " + currentOutput.Name
	descField := "Job Description: " + currentOutput.JobDescription
	generateButton := "[ Generate ]"
	saveButton := "[ Save to PDF ]"

	// Highlight the active element in the job-specific section.
	if m.focusArea == FocusJob {
		switch m.jobField {
		case JobFieldName:
			nameField = ui.SelectedItem.Render("► " + nameField)
		case JobFieldDescription:
			descField = ui.SelectedItem.Render("► " + descField)
		case JobFieldGenerate:
			generateButton = ui.SelectedItem.Render("► " + generateButton)
		case JobFieldSavePDF:
			saveButton = ui.SelectedItem.Render("► " + saveButton)
		}
	}

	content := title + "\n" + nameField + "\n" + descField + "\n\n" + generateButton + "    " + saveButton
	return content
}

// saveProjectConfig constructs a ProjectConfig and writes it to the project's config file.
func (m *ProjectDetailModel) saveProjectConfig() tea.Msg {
	config := types.ProjectConfig{
		Name:        m.overviewProjectName,
		Model:       "",
		ResumeInput: m.resumeInput,
		Outputs:     m.outputs,
	}
	if len(m.llmOptions) > 0 {
		config.Model = m.llmOptions[m.selectedLLMIndex].Name
	}
	err := types.SaveProjectConfig(m.projectDir, config)
	if err != nil {
		return types.ErrorMsg{Error: fmt.Errorf("failed to save project config: %w", err)}
	}
	return nil
}

// saveCurrentOutputToPDF simulates generating a PDF by writing the generated output
// to a file in the project directory.
func (m *ProjectDetailModel) saveCurrentOutputToPDF() tea.Msg {
	if len(m.outputs) == 0 {
		return types.ErrorMsg{Error: fmt.Errorf("no output to save")}
	}
	current := m.outputs[m.selectedOutputIndex]
	pdfPath := filepath.Join(m.projectDir, current.Name+".pdf")
	err := os.WriteFile(pdfPath, []byte(current.GeneratedOutput), 0644)
	if err != nil {
		return types.ErrorMsg{Error: fmt.Errorf("failed to save PDF: %w", err)}
	}
	return nil
}

func (m *ProjectDetailModel) renderLLMSelector() string {
	content := ui.Title.Render("Select LLM Model") + "\n\n"
	for i, model := range m.llmList {
		item := model.Name
		if i == m.selectedLLMIndex {
			item = ui.SelectedItem.Render("► " + item)
		} else {
			item = "  " + item
		}
		content += item + "\n"
	}
	content += "\n" + ui.Help.Render("↑/↓: Navigate • enter: Select • esc: Cancel")

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		ui.FloatBox.Render(content),
	)
}
