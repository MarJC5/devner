package tui

import "github.com/charmbracelet/lipgloss"

var (
	// Palette — kept aligned with the GUI's Tailwind tokens so both
	// surfaces feel like the same app.
	primary      = lipgloss.Color("#4A6CF7") // accent / active tab bg
	primaryMuted = lipgloss.Color("#F0F4FF") // very light on-dark text
	primaryFaint = lipgloss.Color("#A0ACEE") // dim accent text
	primaryDark  = lipgloss.Color("#3A58D8") // watermark / pressed state
	subtle       = lipgloss.Color("#5C5C5C")
	highlight    = lipgloss.Color("#FAFAFA")
	errColor     = lipgloss.Color("#FF5577")
	okColor      = lipgloss.Color("#00C853")

	// Unused yet but exported so scenes can pull them in without
	// re-declaring. Go complains if they're not referenced somewhere.
	_ = primaryMuted
	_ = primaryFaint
	_ = primaryDark

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
