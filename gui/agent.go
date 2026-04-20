package main

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/devner/devner/internal/agent"
	"github.com/devner/devner/internal/llm"
	"github.com/devner/devner/internal/tools"
	"github.com/google/uuid"
	wailsruntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

// Agent glue — runs agent.Loop in a goroutine per "chat turn" and
// forwards Events to the frontend as Wails runtime events. Sessions
// let the frontend distinguish concurrent turns (in practice we run
// at most one at a time, but the ID also lets confirm replies match
// the right loop instance).
//
// Event names:
//   agent:<sid>:delta        { text }
//   agent:<sid>:tool_call    { id, name, args }
//   agent:<sid>:tool_result  { id, name, ok, content }
//   agent:<sid>:need_confirm { id, tool, args }
//   agent:<sid>:done
//   agent:<sid>:err          { msg }

// agentSession tracks an in-flight loop so we can route confirm
// replies back to the right goroutine.
type agentSession struct {
	reply chan bool // destination for AgentConfirm(sid, approve)
}

var (
	agentMu       sync.Mutex
	agentSessions = map[string]*agentSession{}
	agentHistory  = map[string][]llm.Message{} // per-session transcript (retained for follow-ups)
)

// AgentChat starts a new agent turn. Returns the session id; the
// frontend should subscribe to the event names listed above using
// this id before AgentChat returns (Wails guarantees the binding
// call runs on a bg goroutine, but events only fire once the goroutine
// actually starts emitting — subscribing at send-time is the safe
// pattern).
//
// The session's message history is appended to across multiple calls
// with the same sessionID — pass "" on the first call to get a new
// session, pass the returned id to continue.
func (a *App) AgentChat(prompt, providerOverride, sessionID string) (string, error) {
	if a.deps == nil {
		return "", fmt.Errorf("app deps not initialized")
	}
	provider, err := llm.FromConfig(a.deps.Config, providerOverride)
	if err != nil {
		return "", err
	}
	registry := tools.BuildDefault(a.deps)
	loop := agent.New(a.deps, provider, registry)

	// Session id: reuse if provided (multi-turn), else generate.
	if sessionID == "" {
		sessionID = uuid.NewString()
	}
	agentMu.Lock()
	agentSessions[sessionID] = &agentSession{reply: make(chan bool, 1)}
	history := agentHistory[sessionID]
	agentMu.Unlock()

	events := make(chan agent.Event, 32)

	// Pump the loop's events into Wails runtime events.
	go func() {
		defer func() {
			agentMu.Lock()
			delete(agentSessions, sessionID)
			agentMu.Unlock()
		}()

		for ev := range events {
			switch {
			case ev.Err != nil:
				wailsruntime.EventsEmit(a.ctx, "agent:"+sessionID+":err",
					map[string]string{"msg": ev.Err.Error()})

			case ev.TextDelta != "":
				wailsruntime.EventsEmit(a.ctx, "agent:"+sessionID+":delta",
					map[string]string{"text": ev.TextDelta})

			case ev.ToolCall != nil:
				wailsruntime.EventsEmit(a.ctx, "agent:"+sessionID+":tool_call",
					map[string]any{
						"id":   ev.ToolCall.ID,
						"name": ev.ToolCall.Name,
						"args": string(ev.ToolCall.Arguments),
					})

			case ev.ToolResult != nil:
				wailsruntime.EventsEmit(a.ctx, "agent:"+sessionID+":tool_result",
					map[string]any{
						"id":      ev.ToolResult.CallID,
						"name":    ev.ToolResult.Name,
						"ok":      ev.ToolResult.OK,
						"content": ev.ToolResult.Content,
					})

			case ev.NeedConfirm != nil:
				// Emit the confirmation request. The frontend must
				// reply via AgentConfirm(sessionID, approve) which
				// delivers to the session's reply channel below.
				args := string(ev.NeedConfirm.Args)
				wailsruntime.EventsEmit(a.ctx, "agent:"+sessionID+":need_confirm",
					map[string]any{
						"id":   "", // not used for now
						"tool": ev.NeedConfirm.Tool.Name(),
						"args": args,
					})
				agentMu.Lock()
				sess := agentSessions[sessionID]
				agentMu.Unlock()
				if sess == nil {
					ev.NeedConfirm.Reply <- false
					continue
				}
				approved := <-sess.reply
				ev.NeedConfirm.Reply <- approved

			case ev.Done:
				wailsruntime.EventsEmit(a.ctx, "agent:"+sessionID+":done", nil)
			}
		}
	}()

	// Kick off the loop in another goroutine. It sends on `events`
	// (consumed above) and returns the updated history.
	go func() {
		updated := loop.Run(context.Background(), history, prompt, events)
		agentMu.Lock()
		agentHistory[sessionID] = updated
		agentMu.Unlock()
	}()

	return sessionID, nil
}

// AgentConfirm answers an earlier need_confirm event.
func (a *App) AgentConfirm(sessionID string, approve bool) error {
	agentMu.Lock()
	sess, ok := agentSessions[sessionID]
	agentMu.Unlock()
	if !ok {
		return fmt.Errorf("session %q not found (already finished?)", sessionID)
	}
	sess.reply <- approve
	return nil
}

// AgentReset drops the retained history for a session — "Clear chat".
func (a *App) AgentReset(sessionID string) {
	agentMu.Lock()
	delete(agentHistory, sessionID)
	agentMu.Unlock()
}

// ListProviders returns the provider names defined in config.
func (a *App) ListProviders() []string {
	if a.deps == nil {
		return nil
	}
	out := make([]string, 0, len(a.deps.Config.LLM.Providers))
	for name := range a.deps.Config.LLM.Providers {
		out = append(out, name)
	}
	return out
}

// ActiveProvider returns the currently configured default provider.
func (a *App) ActiveProvider() string {
	if a.deps == nil {
		return ""
	}
	return a.deps.Config.LLM.ActiveProvider
}

// unused — keeps json.RawMessage referenced so a future refactor can
// switch ToolCall.Arguments without a pointless go fmt churn.
var _ = json.RawMessage{}
