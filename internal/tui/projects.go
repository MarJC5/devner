package tui

import (
	"context"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/devner/devner/internal/app"
	"github.com/devner/devner/internal/store"
)

type projectItem struct{ p store.Project }

func (i projectItem) Title() string { return i.p.Name }
func (i projectItem) Description() string {
	db := "-"
	if i.p.DBEngine != "" {
		db = i.p.DBEngine + "/" + i.p.DBName
	}
	return fmt.Sprintf("%s • %s • https://%s", i.p.Type, db, i.p.Domain)
}
func (i projectItem) FilterValue() string { return i.p.Name }

type projectsLoadedMsg struct {
	items []list.Item
	err   error
}

type projectsScene struct {
	deps    *app.Deps
	list    list.Model
	width   int
	height  int
	loadErr error
}

func newProjectsScene(d *app.Deps) *projectsScene {
	delegate := list.NewDefaultDelegate()
	l := list.New(nil, delegate, 0, 0)
	l.Title = "Projects"
	l.SetShowStatusBar(false)
	l.SetShowHelp(false)
	l.DisableQuitKeybindings()
	return &projectsScene{deps: d, list: l}
}

func (s *projectsScene) Init() tea.Cmd {
	return s.loadCmd()
}

func (s *projectsScene) loadCmd() tea.Cmd {
	return func() tea.Msg {
		projects, err := s.deps.Store.ListProjects(context.Background())
		if err != nil {
			return projectsLoadedMsg{err: err}
		}
		items := make([]list.Item, 0, len(projects))
		for _, p := range projects {
			items = append(items, projectItem{p: p})
		}
		return projectsLoadedMsg{items: items}
	}
}

func (s *projectsScene) SetSize(w, h int) {
	s.width = w
	s.height = h
	s.list.SetSize(w-4, h-4)
}

func (s *projectsScene) Update(msg tea.Msg) (sceneModel, tea.Cmd) {
	switch msg := msg.(type) {
	case projectsLoadedMsg:
		s.loadErr = msg.err
		if msg.err == nil {
			s.list.SetItems(msg.items)
		}
		return s, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "r":
			return s, s.loadCmd()
		}
	}
	var cmd tea.Cmd
	s.list, cmd = s.list.Update(msg)
	return s, cmd
}

func (s *projectsScene) View() string {
	if s.loadErr != nil {
		return errStyle.Render("Failed to load projects: " + s.loadErr.Error())
	}
	if len(s.list.Items()) == 0 {
		empty := strings.Join([]string{
			"",
			"  No projects yet.",
			"",
			"  Use the Chat scene (press 3) to create one with AI,",
			"  or run:  devner new laravel myapp --db=mysql",
			"",
		}, "\n")
		return lipgloss.NewStyle().Padding(1, 2).Render(empty)
	}
	return lipgloss.NewStyle().Padding(1, 2).Render(s.list.View())
}
