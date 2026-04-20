package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/devner/devner/internal/app"
)

type configScene struct {
	deps *app.Deps
}

func newConfigScene(d *app.Deps) *configScene { return &configScene{deps: d} }

func (s *configScene) Init() tea.Cmd       { return nil }
func (s *configScene) SetSize(w, h int)    {}
func (s *configScene) Update(msg tea.Msg) (sceneModel, tea.Cmd) { return s, nil }

func (s *configScene) View() string {
	cfg := s.deps.Config
	lines := []string{
		lipgloss.NewStyle().Bold(true).Render("Configuration"),
		"",
		fmt.Sprintf("Config file   : %s", cfg.ConfigPath()),
		fmt.Sprintf("Data dir      : %s", cfg.Stack.DataDir),
		fmt.Sprintf("Projects dir  : %s", cfg.Stack.ProjectsDir),
		"",
		lipgloss.NewStyle().Bold(true).Render("LLM providers"),
		fmt.Sprintf("Active        : %s", cfg.LLM.ActiveProvider),
	}
	for name, p := range cfg.LLM.Providers {
		lines = append(lines, fmt.Sprintf("  %-12s  model=%s  key=$%s", name, p.Model, p.APIKeyEnv))
	}
	lines = append(lines, "", mutedStyle.Render("(edit "+cfg.ConfigPath()+" to change)"))

	out := ""
	for _, l := range lines {
		out += l + "\n"
	}
	return lipgloss.NewStyle().Padding(1, 2).Render(out)
}
