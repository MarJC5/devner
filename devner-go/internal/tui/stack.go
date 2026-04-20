package tui

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/devner/devner/internal/app"
)

type stackTickMsg struct{}

type stackStatusMsg struct {
	containers []stackContainer
	err        error
}

type stackContainer struct {
	Name   string
	State  string
	Status string
}

type stackScene struct {
	deps   *app.Deps
	items  []stackContainer
	err    error
	width  int
	height int
}

func newStackScene(d *app.Deps) *stackScene { return &stackScene{deps: d} }

func (s *stackScene) Init() tea.Cmd {
	return tea.Batch(s.refreshCmd(), s.tickCmd())
}

func (s *stackScene) SetSize(w, h int) { s.width = w; s.height = h }

func (s *stackScene) tickCmd() tea.Cmd {
	return tea.Tick(3*time.Second, func(time.Time) tea.Msg { return stackTickMsg{} })
}

func (s *stackScene) refreshCmd() tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		cmd := exec.CommandContext(ctx, "docker", "ps", "--filter", "name=_devner", "--format", "{{.Names}}\t{{.State}}\t{{.Status}}")
		out, err := cmd.Output()
		if err != nil {
			return stackStatusMsg{err: err}
		}
		var containers []stackContainer
		for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
			if line == "" {
				continue
			}
			parts := strings.SplitN(line, "\t", 3)
			if len(parts) != 3 {
				continue
			}
			containers = append(containers, stackContainer{Name: parts[0], State: parts[1], Status: parts[2]})
		}
		return stackStatusMsg{containers: containers}
	}
}

func (s *stackScene) Update(msg tea.Msg) (sceneModel, tea.Cmd) {
	switch msg := msg.(type) {
	case stackTickMsg:
		return s, tea.Batch(s.refreshCmd(), s.tickCmd())
	case stackStatusMsg:
		s.items = msg.containers
		s.err = msg.err
	case tea.KeyMsg:
		switch msg.String() {
		case "r":
			return s, s.refreshCmd()
		}
	}
	return s, nil
}

func (s *stackScene) View() string {
	var b strings.Builder
	b.WriteString(lipgloss.NewStyle().Bold(true).Render("Stack"))
	b.WriteString("\n\n")
	if s.err != nil {
		b.WriteString(errStyle.Render("Error: " + s.err.Error()))
		b.WriteString("\n\nIs Docker running?\n")
	} else if len(s.items) == 0 {
		b.WriteString(mutedStyle.Render("No devner containers running. Press `up` in CLI to start: `devner up`"))
	} else {
		for _, c := range s.items {
			dot := "○"
			color := mutedStyle
			if c.State == "running" {
				dot = "●"
				color = okStyle
			} else if c.State == "restarting" || c.State == "exited" {
				color = errStyle
			}
			line := fmt.Sprintf("%s  %-25s  %-10s  %s", color.Render(dot), c.Name, c.State, c.Status)
			b.WriteString(line)
			b.WriteString("\n")
		}
	}
	b.WriteString("\n")
	b.WriteString(mutedStyle.Render("[r] refresh   auto-refreshes every 3s"))
	return lipgloss.NewStyle().Padding(1, 2).Render(b.String())
}
