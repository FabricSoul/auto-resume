package models

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"errors"

	"github.com/FabricSoul/auto-resume/internal/types"
	"github.com/FabricSoul/auto-resume/internal/ui"
	"github.com/teilomillet/gollm"
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
	isGenerating     bool
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
		// Reset generating state when going back to splash screen
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			m.isGenerating = false
			return m, func() tea.Msg {
				return types.TransitionMsg{To: types.StateSplash}
			}
		}

		// Don't process other key events while generating
		if m.isGenerating {
			return m, nil
		}

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
		case "tab":
			m.focusArea = (m.focusArea + 1) % 3
		case "shift+tab":
			m.focusArea = (m.focusArea + 2) % 3
		case "i":
			switch m.focusArea {
			case FocusOverview:
				var prompt, initialValue string
				var callback func(string)

				switch m.overviewField {
				case OverviewFieldProjectName:
					prompt = "Enter Project Name"
					initialValue = m.overviewProjectName
					callback = func(value string) {
						m.overviewProjectName = value
					}
				case OverviewFieldResumeInput:
					prompt = "Enter Resume Input"
					initialValue = m.resumeInput
					callback = func(value string) {
						m.resumeInput = value
					}
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
			case FocusJob:
				var prompt, initialValue string
				var callback func(string)

				current := &m.outputs[m.selectedOutputIndex]
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
			}
		case "a":
			if m.focusArea == FocusOutputs {
				newOutput := types.Output{
					Name:           "New Output",
					JobDescription: "",
				}
				m.outputs = append(m.outputs, newOutput)
				m.selectedOutputIndex = len(m.outputs) - 1
				return m, m.saveProjectConfig
			}
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

			switch msg.String() {
			case "i":
				var prompt, initialValue string
				var callback func(string)

				current := &m.outputs[m.selectedOutputIndex]
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
					currentOutput := m.outputs[m.selectedOutputIndex]
					return m, func() tea.Msg {
						return m.generateResume(currentOutput.JobDescription)
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
	if m.isGenerating {
		return lipgloss.Place(
			m.width,
			m.height,
			lipgloss.Center,
			lipgloss.Center,
			ui.FloatBox.Render("Generating resume...\nPlease stand by..."),
		)
	}

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

	// Truncate resume input if longer than 10 chars
	resumePreview := m.resumeInput
	if len(resumePreview) > 10 {
		resumePreview = resumePreview[:10] + "..."
	}
	resumeField := "Resume Input: " + resumePreview

	llmField := "LLM: "
	if len(m.llmOptions) > 0 {
		llmField += m.llmOptions[m.selectedLLMIndex].Name
	} else {
		llmField += "None"
	}

	// Highlight the active field if the overview section has focus
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

	// Truncate job description if longer than 10 chars
	descPreview := currentOutput.JobDescription
	if len(descPreview) > 10 {
		descPreview = descPreview[:10] + "..."
	}

	// Truncate output preview if longer than 10 chars
	outputPreview := currentOutput.GeneratedOutput
	if len(outputPreview) > 10 {
		outputPreview = outputPreview[:10] + "..."
	}

	nameField := "Output Name: " + currentOutput.Name
	descField := "Job Description: " + descPreview
	outputField := "Generated Output: " + outputPreview
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

	content := title + "\n" + nameField + "\n" + descField + "\n" + outputField + "\n\n" + generateButton + "    " + saveButton
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

// Add these new types for the API request/response
type GollmRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type GollmResponse struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
}

// Add this method to ProjectDetailModel
func (m *ProjectDetailModel) generateResume(jobDescription string) tea.Msg {
	if m.isGenerating {
		return nil // Prevent multiple generations at once
	}

	m.isGenerating = true
	return func() tea.Msg {
		// Always reset the generating state, even if there's an error
		defer func() {
			m.isGenerating = false
		}()

		if len(m.llmOptions) == 0 || m.selectedLLMIndex >= len(m.llmOptions) {
			return types.ErrorMsg{Error: errors.New("no LLM model selected")}
		}

		selectedModel := m.llmOptions[m.selectedLLMIndex]

		// Load project config to get resume input
		config, err := types.LoadProjectConfig(m.projectDir)
		if err != nil {
			return types.ErrorMsg{Error: fmt.Errorf("failed to load project config: %w", err)}
		}

		// Create a new LLM instance with debug logging
		llm, err := gollm.NewLLM(
			gollm.SetProvider(selectedModel.Provider),
			gollm.SetModel(selectedModel.Model),
			gollm.SetAPIKey(selectedModel.APIKey),
		)
		if err != nil {
			return types.ErrorMsg{Error: fmt.Errorf("failed to create LLM instance: %w", err)}
		}

		// Prepare the prompt
		promptText := `You are a professional resume writer. Your task is to modify the given resume to better target a specific job description.
Follow these rules:
1. Keep the same LaTeX format
2. Highlight relevant skills and experiences
3. Use keywords from the job description
4. Be concise and professional
5. Do not invent new experiences

Original Resume:
%s

Job Description:
%s

Please provide the modified resume in LaTeX format.`

		promptText = fmt.Sprintf(promptText, config.ResumeInput, jobDescription)
		prompt := gollm.NewPrompt(promptText)

		// Generate response with context and error handling
		ctx := context.Background()
		response, err := llm.Generate(ctx, prompt)
		if err != nil {
			// Handle the error based on the error message
			errMsg := err.Error()
			switch {
			case strings.Contains(errMsg, "invalid API key"):
				return types.ErrorMsg{Error: fmt.Errorf("invalid API key for %s", selectedModel.Provider)}
			case strings.Contains(errMsg, "rate limit"):
				return types.ErrorMsg{Error: fmt.Errorf("rate limit exceeded for %s", selectedModel.Provider)}
			case strings.Contains(errMsg, "model not found"):
				return types.ErrorMsg{Error: fmt.Errorf("model %s not found for %s", selectedModel.Model, selectedModel.Provider)}
			case strings.Contains(errMsg, "connection refused"):
				return types.ErrorMsg{Error: fmt.Errorf("connection to %s failed - is the service running?", selectedModel.Provider)}
			default:
				return types.ErrorMsg{Error: fmt.Errorf("LLM error: %w", err)}
			}
		}

		if response == "" {
			return types.ErrorMsg{Error: errors.New("received empty response from LLM")}
		}

		// Update project config with new output
		outputName := time.Now().Format("2006-01-02-15-04-05")
		newOutput := types.Output{
			Name:            outputName,
			JobDescription:  jobDescription,
			GeneratedOutput: response,
		}

		config.Outputs = append(config.Outputs, newOutput)
		if err := types.SaveProjectConfig(m.projectDir, config); err != nil {
			return types.ErrorMsg{Error: fmt.Errorf("failed to save project config: %w", err)}
		}

		// Update model's outputs list
		m.outputs = config.Outputs
		m.selectedOutputIndex = len(m.outputs) - 1

		// If we get here, generation was successful
		return types.GenerationCompleteMsg{} // Add this new message type
	}
}
