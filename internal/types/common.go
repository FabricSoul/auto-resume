// internal/types/common.go
package types

import tea "github.com/charmbracelet/bubbletea"

type Appstate int

const (
	StateSplash Appstate = iota
	StateNewProject
	StateProjectList
	StateProjectOverview
	StateSettings
)

type Model interface {
	Init() tea.Cmd
	Update(tea.Msg) (tea.Model, tea.Cmd)
	View() string
}
