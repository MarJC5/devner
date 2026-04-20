package tui

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/devner/devner/internal/app"
)

type Scene int

const (
	SceneProjects Scene = iota
	SceneStack
	SceneChat
	SceneLogs
	SceneConfig
)

var sceneNames = []string{"Projects", "Stack", "Chat", "Logs", "Config"}

type sceneModel interface {
	Init() tea.Cmd
	Update(tea.Msg) (sceneModel, tea.Cmd)
	View() string
	SetSize(width, height int)
}

type RootModel struct {
	scene  Scene
	width  int
	height int
	deps   *app.Deps
	scenes [5]sceneModel
	status string
}

func NewRootModel(d *app.Deps) RootModel {
	m := RootModel{scene: SceneProjects, deps: d}
	m.scenes[SceneProjects] = newProjectsScene(d)
	m.scenes[SceneStack] = newStackScene(d)
	m.scenes[SceneChat] = newChatScene(d)
	m.scenes[SceneLogs] = newLogsScene(d)
	m.scenes[SceneConfig] = newConfigScene(d)
	return m
}

func (m RootModel) Init() tea.Cmd {
	var cmds []tea.Cmd
	for _, s := range m.scenes {
		if c := s.Init(); c != nil {
			cmds = append(cmds, c)
		}
	}
	return tea.Batch(cmds...)
}

func (m RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		for i, s := range m.scenes {
			s.SetSize(msg.Width, m.bodyHeight())
			m.scenes[i] = s
		}
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "q":
			if m.scene != SceneChat {
				return m, tea.Quit
			}
		case "tab":
			m.scene = (m.scene + 1) % Scene(len(sceneNames))
			return m, m.scenes[m.scene].Init()
		case "shift+tab":
			m.scene = (m.scene - 1 + Scene(len(sceneNames))) % Scene(len(sceneNames))
			return m, m.scenes[m.scene].Init()
		case "1", "2", "3", "4", "5":
			idx := int(msg.String()[0] - '1')
			if idx >= 0 && idx < len(sceneNames) {
				m.scene = Scene(idx)
				return m, m.scenes[m.scene].Init()
			}
		}
	case statusMsg:
		m.status = string(msg)
		return m, nil
	}

	updated, cmd := m.scenes[m.scene].Update(msg)
	m.scenes[m.scene] = updated
	return m, cmd
}

func (m RootModel) View() string {
	var tabs []string
	for i, name := range sceneNames {
		style := tabStyle
		if Scene(i) == m.scene {
			style = activeTabStyle
		}
		tabs = append(tabs, style.Render(name))
	}
	header := lipgloss.JoinHorizontal(lipgloss.Top, tabs...)

	body := m.scenes[m.scene].View()

	hint := "[tab] next scene   [1-5] jump   [q] quit"
	if m.scene == SceneChat {
		hint = "[tab] leave chat   [enter] send   [ctrl+c] quit"
	}
	footerText := hint
	if m.status != "" {
		footerText = m.status + "  •  " + hint
	}
	footer := footerStyle.Width(m.width).Render(footerText)

	return lipgloss.JoinVertical(lipgloss.Left, header, body, footer)
}

func (m RootModel) bodyHeight() int {
	h := m.height - 3
	if h < 10 {
		h = 10
	}
	return h
}

// statusMsg lets scenes push a status line into the footer.
type statusMsg string

func SetStatus(s string) tea.Cmd {
	return func() tea.Msg { return statusMsg(s) }
}

// Run starts the TUI. Accepts a ctx for graceful shutdown of background
// tickers; current impl relies on tea.Quit instead.
func Run(ctx context.Context, d *app.Deps) error {
	p := tea.NewProgram(NewRootModel(d), tea.WithAltScreen(), tea.WithContext(ctx))
	_, err := p.Run()
	return err
}

// _ keeps ctx imported when not used directly.
var _ = context.Background

