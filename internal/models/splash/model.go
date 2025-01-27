// internal/splash/model.go
package splash

import (
	"github.com/FabricSoul/auto-resume/internal/types"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type menuItem struct {
	title string
	desc  string
}

func (i menuItem) Title() string       { return i.title }
func (i menuItem) Description() string { return i.desc }
func (i menuItem) FilterValue() string { return i.title }

type SplashModel struct {
	list    list.Model
	project *types.ProjectManager
}

func NewSplasheModel(pm *types.ProjectManager) *SplashModel {
	items := []list.Item{
		menuItem{title: "[N]ew project", desc: "Start a new project"},
		menuItem{title: "[R]esume project", desc: "Resume the last project"},
		menuItem{title: "[S]ettings", desc: "Settings"},
		menuItem{title: "[Q]uit", desc: "Quit"},
	}

	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Auto Resume"
	l.SetShowTitle(true)
	l.ShowHelp()

	return &SplashModel{
		list:    l,
		project: pm,
	}
}

func (m *SplashModel) Init() tea.Cmd {
	return nil
}

func (m *SplashModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h := msg.Width
		v := msg.Height
		m.list.SetSize(h, v)
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "n", "N":
			return m, func() tea.Msg {
				return types.TransitionMsg{
					To: types.StateNewProject,
				}
			}
		case "r", "R":
			return m, nil
		case "s", "S":
			return m, nil
		}
	}
	return m, nil
}

func (m *SplashModel) View() string {
	return m.list.View()
}
