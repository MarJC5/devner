package tui

import (
	"context"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/devner/devner/internal/app"
)

// promptModel is a single-scene variant of the full TUI: just the chat
// view, full-screen, closes on Esc/Ctrl+C. Designed to be the landing
// surface when a global hotkey (Cmd+D on macOS, Ctrl+Shift+D on
// Windows/Linux — bound by the user in their OS shortcuts) launches
// `devner prompt`.
type promptModel struct {
	chat   *chatScene
	width  int
	height int
}

func newPromptModel(d *app.Deps) promptModel {
	return promptModel{chat: newChatScene(d)}
}

func (m promptModel) Init() tea.Cmd { return m.chat.Init() }

func (m promptModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.chat.SetSize(msg.Width, msg.Height-2)
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" || msg.String() == "esc" {
			return m, tea.Quit
		}
	}
	updated, cmd := m.chat.Update(msg)
	m.chat = updated.(*chatScene)
	return m, cmd
}

func (m promptModel) View() string {
	header := promptHeaderStyle.Width(m.width).Render("Devner Prompt — press Esc to close")
	body := m.chat.View()
	return lipgloss.JoinVertical(lipgloss.Left, header, body)
}

var promptHeaderStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#FAFAFA")).
	Background(lipgloss.Color("#7D56F4")).
	Padding(0, 2)

// RunPrompt launches the focused single-scene prompt UI.
func RunPrompt(ctx context.Context, d *app.Deps) error {
	p := tea.NewProgram(newPromptModel(d), tea.WithAltScreen(), tea.WithContext(ctx))
	_, err := p.Run()
	if err != nil {
		return fmt.Errorf("prompt: %w", err)
	}
	return nil
}
