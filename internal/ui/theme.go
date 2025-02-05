package ui

import "github.com/charmbracelet/lipgloss"

var (
	// Base colors
	Subtle    = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
	Highlight = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	Special   = lipgloss.AdaptiveColor{Light: "#43BF6D", Dark: "#73F59F"}
	Error     = lipgloss.AdaptiveColor{Light: "#FF0000", Dark: "#FF4444"}

	// Base styles
	BaseList = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Subtle).
			Padding(1).
			MarginRight(2)

	BaseDetails = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Subtle).
			Padding(1)

	Title = lipgloss.NewStyle().
		Foreground(Highlight).
		Bold(true).
		MarginLeft(1).
		MarginBottom(1)

	Help = lipgloss.NewStyle().
		Foreground(Subtle).
		MarginTop(1)

	SelectedItem = lipgloss.NewStyle().
			Foreground(Special).
			Bold(true)

	ErrorBox = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Error).
			Padding(1).
			Width(50)

	ErrorTitle = lipgloss.NewStyle().
			Foreground(Error).
			Bold(true).
			MarginBottom(1)

	Input = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(Subtle).
		Padding(1)

	// AppBackground style to be applied across the entire app.
	AppBackground = lipgloss.NewStyle().
			Background(lipgloss.Color("#1E1E1E")). // A solid dark color.
			Foreground(lipgloss.Color("#FFFFFF"))  // Use white text for high contrast.

	FloatBox = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Highlight).
			Padding(2).
			Width(60)

	JoinedContainer = lipgloss.NewStyle().
			Background(lipgloss.Color("#1E1E1E")).
			Padding(1)
)
