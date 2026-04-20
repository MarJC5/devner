package tui

import "github.com/charmbracelet/lipgloss"

var (
	primary   = lipgloss.Color("#7D56F4")
	subtle    = lipgloss.Color("#5C5C5C")
	highlight = lipgloss.Color("#FAFAFA")
	errColor  = lipgloss.Color("#FF5577")
	okColor   = lipgloss.Color("#00C853")

	tabStyle = lipgloss.NewStyle().
			Padding(0, 2).
			Foreground(subtle)

	activeTabStyle = lipgloss.NewStyle().
			Padding(0, 2).
			Bold(true).
			Foreground(highlight).
			Background(primary)

	sceneBodyStyle = lipgloss.NewStyle().
			Padding(1, 2).
			Border(lipgloss.NormalBorder(), true, false, true, false).
			BorderForeground(subtle)

	footerStyle = lipgloss.NewStyle().
			Foreground(subtle).
			Padding(0, 2)

	errStyle   = lipgloss.NewStyle().Foreground(errColor)
	okStyle    = lipgloss.NewStyle().Foreground(okColor)
	mutedStyle = lipgloss.NewStyle().Foreground(subtle)
)

var _ = sceneBodyStyle // may be used by future scenes
