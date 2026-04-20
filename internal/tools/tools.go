// Package tools defines the Tool interface and a Registry of tools that
// the LLM can call. Each tool owns its JSON schema, an Execute function,
// and a Destructive flag that forces a TUI confirmation before running.
//
// Design note: tools DO NOT import internal/agent or internal/llm — they
// only depend on the domain services (runtime, database, project, network,
// store). This keeps the dependency graph a DAG.
package tools

import (
	"context"
	"encoding/json"
	"fmt"
)

type Tool interface {
	Name() string
	Description() string
	Destructive() bool
	Schema() json.RawMessage
	Execute(ctx context.Context, args json.RawMessage) (Result, error)
}

type Result struct {
	// Content is the textual result sent back to the LLM. Keep under 4KB
	// when possible; truncate long outputs.
	Content string
	// Data is optional structured data for the TUI to render.
	Data any
}

type Registry struct {
	tools []Tool
	byKey map[string]Tool
}

func NewRegistry(tools ...Tool) *Registry {
	r := &Registry{byKey: map[string]Tool{}}
	for _, t := range tools {
		r.Add(t)
	}
	return r
}

func (r *Registry) Add(t Tool) {
	r.tools = append(r.tools, t)
	r.byKey[t.Name()] = t
}

func (r *Registry) List() []Tool { return r.tools }

func (r *Registry) Get(name string) (Tool, bool) {
	t, ok := r.byKey[name]
	return t, ok
}

// Execute dispatches a tool call. Returns an error string as Result.Content
// when the tool fails — the LLM can then reason about the failure.
func (r *Registry) Execute(ctx context.Context, name string, args json.RawMessage) (Result, error) {
	t, ok := r.Get(name)
	if !ok {
		return Result{Content: fmt.Sprintf("unknown tool %q", name)}, fmt.Errorf("unknown tool %q", name)
	}
	return t.Execute(ctx, args)
}

// schema is a tiny helper to inline JSON schemas as RawMessage.
func schema(s string) json.RawMessage { return json.RawMessage(s) }
