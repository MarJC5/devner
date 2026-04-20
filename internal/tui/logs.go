package tui

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/devner/devner/internal/app"
)

var stackServices = []string{"frankenphp", "mysql", "postgres", "redis", "mailpit", "adminer"}

type logsMsg struct {
	service string
	lines   string
	err     error
}

type logsScene struct {
	deps     *app.Deps
	service  int // index in stackServices
	viewport viewport.Model
	lines    string
	err      error
}

func newLogsScene(d *app.Deps) *logsScene {
	vp := viewport.New(80, 20)
	return &logsScene{deps: d, viewport: vp}
}

func (s *logsScene) Init() tea.Cmd { return s.loadCmd() }

func (s *logsScene) SetSize(w, h int) {
	s.viewport.Width = w - 4
	s.viewport.Height = h - 6
}

func (s *logsScene) loadCmd() tea.Cmd {
	svc := stackServices[s.service]
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		out, err := exec.CommandContext(ctx, "docker", "logs", "--tail", "200", svc+"_devner").CombinedOutput()
		return logsMsg{service: svc, lines: string(out), err: err}
	}
}

func (s *logsScene) Update(msg tea.Msg) (sceneModel, tea.Cmd) {
	switch msg := msg.(type) {
	case logsMsg:
		s.err = msg.err
		s.lines = msg.lines
		s.viewport.SetContent(s.lines)
		s.viewport.GotoBottom()
		return s, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "n", "right":
			s.service = (s.service + 1) % len(stackServices)
			return s, s.loadCmd()
		case "p", "left":
			s.service = (s.service - 1 + len(stackServices)) % len(stackServices)
			return s, s.loadCmd()
		case "r":
			return s, s.loadCmd()
		}
	}
	var cmd tea.Cmd
	s.viewport.Update(msg)
	s.viewport, cmd = s.viewport.Update(msg)
	return s, cmd
}

func (s *logsScene) View() string {
	var tabs []string
	for i, svc := range stackServices {
		style := tabStyle
		if i == s.service {
			style = activeTabStyle
		}
		tabs = append(tabs, style.Render(svc))
	}
	header := lipgloss.JoinHorizontal(lipgloss.Top, tabs...)

	body := s.viewport.View()
	if s.err != nil {
		body = errStyle.Render("Error: " + s.err.Error())
	} else if strings.TrimSpace(s.lines) == "" {
		body = mutedStyle.Render("(no logs — container may not be running)")
	}

	hint := mutedStyle.Render("[n/p] switch service   [r] refresh   [↑/↓] scroll")
	return lipgloss.NewStyle().Padding(1, 2).Render(fmt.Sprintf("%s\n\n%s\n\n%s", header, body, hint))
}
