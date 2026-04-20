// Package agent implements the tool-use loop: user → LLM → tool calls →
// tool results → LLM → ... until the model produces a final text message.
//
// Events are emitted on a channel so the TUI (or CLI verbose mode) can
// render progress without polling.
package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/devner/devner/internal/app"
	"github.com/devner/devner/internal/llm"
	"github.com/devner/devner/internal/store"
	"github.com/devner/devner/internal/tools"
)

const (
	// MaxIterations caps the tool-use loop. Each iteration is one LLM
	// round-trip; prevents runaways if the model keeps calling tools.
	MaxIterations = 10
)

// Event is emitted as the loop progresses. Exactly one field is set.
type Event struct {
	TextDelta  string
	ToolCall   *llm.ToolCall
	ToolResult *ToolResult
	// NeedConfirm is emitted before executing a destructive tool call.
	// The caller must reply via the returned confirm channel: true to
	// proceed, false to cancel.
	NeedConfirm *ConfirmRequest
	Done        bool
	Err         error
}

type ToolResult struct {
	CallID  string
	Name    string
	Content string
	OK      bool
}

type ConfirmRequest struct {
	Tool    tools.Tool
	Args    json.RawMessage
	Reply   chan<- bool
}

type Loop struct {
	Provider   llm.Provider
	Tools      *tools.Registry
	Deps       *app.Deps
	SystemMsg  string
	// AutoConfirm bypasses the destructive confirmation step. CLI uses
	// this with an explicit --yes flag; TUI never does.
	AutoConfirm bool
}

func New(d *app.Deps, p llm.Provider, r *tools.Registry) *Loop {
	return &Loop{
		Provider:  p,
		Tools:     r,
		Deps:      d,
		SystemMsg: DefaultSystemMessage(),
	}
}

func DefaultSystemMessage() string {
	return `You are Devner, an assistant that manages a local development environment on Docker (WordPress, Laravel, Node, Next.js, Astro).

BEHAVIOUR RULES:
1. Call tools first, ask questions later. The user's projects live in a local store you can query instantly — never claim you don't know a project without calling list_projects or project_status.
2. When the user mentions a project name, call project_status(name) before anything else. If it doesn't exist, then say so.
3. Pick the most specific tool:
   - composer / npm / wp_cli / artisan for their ecosystems
   - exec_in_project only when no typed tool fits (e.g. "php -v", "tail logs", "grep in files")
4. Call destructive tools (delete_project, drop_database, rebuild_stack, composer, npm, wp_cli, artisan, exec_in_project) without asking "are you sure?" — the UI prompts the user before anything destructive actually runs.
5. Keep answers short. Render results as concise markdown (tables, lists). Never dump raw stdout longer than ~15 lines; summarize instead.
6. Use read-only tools (list_projects, project_status, tail_logs) liberally to ground your answers in the current state rather than guessing from training data.`
}

// Run executes the agent loop for a single user turn.
//
// `history` is the conversation so far including system if desired (the
// loop injects the default system message if none is present). The loop
// appends assistant + tool messages to `history` as it progresses; callers
// can retain the final slice as the next turn's input.
//
// Events are sent on `events`. The channel is closed when the loop ends.
func (l *Loop) Run(ctx context.Context, history []llm.Message, userInput string, events chan<- Event) []llm.Message {
	defer close(events)

	// Inject default system if history is empty.
	if len(history) == 0 && l.SystemMsg != "" {
		history = append(history, llm.Message{Role: llm.RoleSystem, Content: l.SystemMsg})
	}
	history = append(history, llm.Message{Role: llm.RoleUser, Content: userInput})

	schemas := toSchemas(l.Tools)

	for iter := 0; iter < MaxIterations; iter++ {
		stream, err := l.Provider.Stream(ctx, history, schemas)
		if err != nil {
			events <- Event{Err: fmt.Errorf("llm stream: %w", err)}
			return history
		}

		var final *llm.Response
		for ev := range stream {
			if ev.Err != nil {
				events <- Event{Err: ev.Err}
				return history
			}
			if ev.TextDelta != "" {
				events <- Event{TextDelta: ev.TextDelta}
			}
			if ev.Done != nil {
				final = ev.Done
			}
		}
		if final == nil {
			events <- Event{Err: fmt.Errorf("provider returned no final response")}
			return history
		}

		// Append assistant message.
		history = append(history, llm.Message{
			Role:      llm.RoleAssistant,
			Content:   final.Content,
			ToolCalls: final.ToolCalls,
		})

		if len(final.ToolCalls) == 0 {
			events <- Event{Done: true}
			return history
		}

		// Execute each tool call.
		for _, tc := range final.ToolCalls {
			tc := tc
			events <- Event{ToolCall: &tc}

			t, ok := l.Tools.Get(tc.Name)
			if !ok {
				history = append(history, llm.Message{
					Role:       llm.RoleTool,
					ToolCallID: tc.ID,
					Name:       tc.Name,
					Content:    fmt.Sprintf("unknown tool %q", tc.Name),
				})
				events <- Event{ToolResult: &ToolResult{CallID: tc.ID, Name: tc.Name, Content: "unknown tool", OK: false}}
				continue
			}

			if t.Destructive() && !l.AutoConfirm {
				reply := make(chan bool, 1)
				events <- Event{NeedConfirm: &ConfirmRequest{Tool: t, Args: tc.Arguments, Reply: reply}}
				approved, ok := <-reply
				if !ok || !approved {
					res := "cancelled by user"
					history = append(history, llm.Message{
						Role:       llm.RoleTool,
						ToolCallID: tc.ID,
						Name:       tc.Name,
						Content:    res,
					})
					events <- Event{ToolResult: &ToolResult{CallID: tc.ID, Name: tc.Name, Content: res, OK: false}}
					continue
				}
			}

			start := time.Now()
			res, execErr := t.Execute(ctx, tc.Arguments)
			dur := time.Since(start)
			ok = execErr == nil

			// Audit
			_ = l.Deps.Store.RecordHistory(ctx, store.HistoryEntry{
				TS:         time.Now(),
				Tool:       tc.Name,
				ArgsJSON:   string(tc.Arguments),
				ResultJSON: res.Content,
				DurationMs: dur.Milliseconds(),
				OK:         ok,
			})

			history = append(history, llm.Message{
				Role:       llm.RoleTool,
				ToolCallID: tc.ID,
				Name:       tc.Name,
				Content:    res.Content,
			})
			events <- Event{ToolResult: &ToolResult{CallID: tc.ID, Name: tc.Name, Content: res.Content, OK: ok}}
		}
	}

	events <- Event{Err: fmt.Errorf("agent loop exceeded max iterations (%d)", MaxIterations)}
	return history
}

func toSchemas(r *tools.Registry) []llm.ToolSchema {
	all := r.List()
	out := make([]llm.ToolSchema, 0, len(all))
	for _, t := range all {
		out = append(out, llm.ToolSchema{
			Name:        t.Name(),
			Description: t.Description(),
			Parameters:  t.Schema(),
		})
	}
	return out
}
