package tui

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/devner/devner/internal/agent"
	"github.com/devner/devner/internal/app"
	"github.com/devner/devner/internal/llm"
	"github.com/devner/devner/internal/tools"
)

// chatScene drives the agent loop interactively.
//
// The Bubble Tea model owns a goroutine that runs Loop.Run and a channel
// of agent.Event values. We forward channel events into tea.Msg via a
// recurring Cmd that reads one event per tick. Confirmations are rendered
// as a pending prompt the user answers with y/n.
type chatScene struct {
	deps       *app.Deps
	transcript viewport.Model
	input      textarea.Model

	provider llm.Provider
	registry *tools.Registry
	history  []llm.Message

	running bool
	events  chan agent.Event

	pendingConfirm *agent.ConfirmRequest

	// Rendered transcript lines. We keep both text and formatted form so we
	// can re-render on resize.
	log []string

	// Pending assistant buffer (streaming) — flushed into log on Done.
	asstBuf string

	initErr error
}

func newChatScene(d *app.Deps) *chatScene {
	ta := textarea.New()
	ta.Placeholder = "Ask the agent — e.g. 'create a laravel project devblog with postgres'"
	ta.Prompt = "▶ "
	ta.SetHeight(3)
	ta.ShowLineNumbers = false
	ta.CharLimit = 4000
	ta.Focus()

	vp := viewport.New(80, 15)

	s := &chatScene{
		deps:       d,
		transcript: vp,
		input:      ta,
		registry:   tools.BuildDefault(d),
	}

	provider, err := llm.FromConfig(d.Config, "")
	if err != nil {
		s.initErr = err
	} else {
		s.provider = provider
	}

	s.appendLog(mutedStyle.Render(fmt.Sprintf("provider: %s    tools: %d    type /help", s.providerLabel(), len(s.registry.List()))))
	s.refreshViewport()
	return s
}

func (s *chatScene) providerLabel() string {
	if s.provider != nil {
		return s.provider.Name()
	}
	return "(unavailable — set API key env var and restart)"
}

func (s *chatScene) Init() tea.Cmd { return textarea.Blink }

func (s *chatScene) SetSize(w, h int) {
	s.transcript.Width = w - 4
	s.transcript.Height = h - 8
	s.input.SetWidth(w - 4)
	s.refreshViewport()
}

func (s *chatScene) appendLog(line string) {
	s.log = append(s.log, line)
}

func (s *chatScene) refreshViewport() {
	content := strings.Join(s.log, "\n")
	if s.asstBuf != "" {
		content += "\n" + asstStyle.Render("agent: ") + s.asstBuf
	}
	if s.pendingConfirm != nil {
		argsPretty := string(s.pendingConfirm.Args)
		if compact, err := json.Marshal(json.RawMessage(s.pendingConfirm.Args)); err == nil {
			argsPretty = string(compact)
		}
		warn := warnStyle.Render(fmt.Sprintf(
			"⚠ Destructive tool call: %s(%s)   [y] approve   [n] cancel",
			s.pendingConfirm.Tool.Name(), argsPretty,
		))
		content += "\n" + warn
	}
	s.transcript.SetContent(content)
	s.transcript.GotoBottom()
}

// nextEventCmd produces a tea.Cmd that blocks on the events channel and
// returns the next event as a msg. Returning this in a loop is how we turn
// a Go channel into the Update() message flow without blocking.
func (s *chatScene) nextEventCmd() tea.Cmd {
	ch := s.events
	return func() tea.Msg {
		ev, ok := <-ch
		if !ok {
			return agentClosedMsg{}
		}
		return agentEventMsg(ev)
	}
}

type agentEventMsg agent.Event
type agentClosedMsg struct{}

func (s *chatScene) submit(text string) tea.Cmd {
	if s.provider == nil {
		s.appendLog(errStyle.Render("error: " + s.initErr.Error()))
		s.refreshViewport()
		return nil
	}
	if s.running {
		return nil
	}
	s.appendLog(userStyle.Render("you: ") + text)
	s.events = make(chan agent.Event, 16)
	s.running = true
	s.asstBuf = ""

	loop := agent.New(s.deps, s.provider, s.registry)
	history := s.history
	go func() {
		s.history = loop.Run(context.Background(), history, text, s.events)
	}()
	s.refreshViewport()
	return s.nextEventCmd()
}

func (s *chatScene) Update(msg tea.Msg) (sceneModel, tea.Cmd) {
	switch msg := msg.(type) {
	case agentClosedMsg:
		s.running = false
		if s.asstBuf != "" {
			s.appendLog(asstStyle.Render("agent: ") + s.asstBuf)
			s.asstBuf = ""
		}
		s.refreshViewport()
		return s, nil

	case agentEventMsg:
		ev := agent.Event(msg)
		switch {
		case ev.Err != nil:
			s.appendLog(errStyle.Render("error: " + ev.Err.Error()))
			s.running = false
		case ev.TextDelta != "":
			s.asstBuf += ev.TextDelta
		case ev.ToolCall != nil:
			s.appendLog(toolStyle.Render(fmt.Sprintf("⏺ %s(%s)", ev.ToolCall.Name, truncate(string(ev.ToolCall.Arguments), 120))))
		case ev.ToolResult != nil:
			status := "✓"
			if !ev.ToolResult.OK {
				status = "✗"
			}
			s.appendLog(mutedStyle.Render(fmt.Sprintf("  %s %s", status, truncate(ev.ToolResult.Content, 200))))
		case ev.NeedConfirm != nil:
			s.pendingConfirm = ev.NeedConfirm
			s.refreshViewport()
			// Don't read next event yet — wait for user y/n.
			return s, nil
		case ev.Done:
			// Flush any buffered text as a final assistant line.
			if s.asstBuf != "" {
				s.appendLog(asstStyle.Render("agent: ") + s.asstBuf)
				s.asstBuf = ""
			}
		}
		s.refreshViewport()
		return s, s.nextEventCmd()

	case tea.KeyMsg:
		if s.pendingConfirm != nil {
			switch msg.String() {
			case "y", "Y", "enter":
				s.pendingConfirm.Reply <- true
				s.pendingConfirm = nil
				s.refreshViewport()
				return s, s.nextEventCmd()
			case "n", "N", "esc":
				s.pendingConfirm.Reply <- false
				s.pendingConfirm = nil
				s.refreshViewport()
				return s, s.nextEventCmd()
			}
			return s, nil
		}
		if msg.String() == "enter" && !s.running {
			text := strings.TrimSpace(s.input.Value())
			if text == "" {
				return s, nil
			}
			s.input.Reset()
			cmd := s.submit(text)
			return s, cmd
		}
	}

	var cmds []tea.Cmd
	var cmd tea.Cmd
	s.input, cmd = s.input.Update(msg)
	cmds = append(cmds, cmd)
	s.transcript, cmd = s.transcript.Update(msg)
	cmds = append(cmds, cmd)
	return s, tea.Batch(cmds...)
}

func (s *chatScene) View() string {
	return lipgloss.NewStyle().Padding(1, 2).Render(
		s.transcript.View() + "\n\n" + s.input.View(),
	)
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "…"
}

var (
	userStyle = lipgloss.NewStyle().Foreground(primary).Bold(true)
	asstStyle = lipgloss.NewStyle().Foreground(okColor).Bold(true)
	toolStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFCC00"))
	warnStyle = lipgloss.NewStyle().Foreground(errColor).Bold(true)
)
