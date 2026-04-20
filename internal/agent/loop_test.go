package agent

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/devner/devner/internal/app"
	"github.com/devner/devner/internal/config"
	"github.com/devner/devner/internal/llm"
	"github.com/devner/devner/internal/store"
	"github.com/devner/devner/internal/tools"
)

// stubProvider scripts a deterministic sequence of Responses so we can
// exercise the loop without a real LLM.
type stubProvider struct {
	responses []llm.Response
	idx       int
}

func (s *stubProvider) Name() string { return "stub" }
func (s *stubProvider) Chat(ctx context.Context, msgs []llm.Message, tools []llm.ToolSchema) (*llm.Response, error) {
	r := s.responses[s.idx]
	s.idx++
	return &r, nil
}
func (s *stubProvider) Stream(ctx context.Context, msgs []llm.Message, tools []llm.ToolSchema) (<-chan llm.StreamEvent, error) {
	r := s.responses[s.idx]
	s.idx++
	ch := make(chan llm.StreamEvent, 2)
	if r.Content != "" {
		ch <- llm.StreamEvent{TextDelta: r.Content}
	}
	rc := r
	ch <- llm.StreamEvent{Done: &rc}
	close(ch)
	return ch, nil
}

// echoTool is a trivial non-destructive tool that returns its input.
type echoTool struct{}

func (e *echoTool) Name() string             { return "echo" }
func (e *echoTool) Description() string      { return "Returns the input." }
func (e *echoTool) Destructive() bool        { return false }
func (e *echoTool) Schema() json.RawMessage  {
	return json.RawMessage(`{"type":"object","properties":{"msg":{"type":"string"}},"required":["msg"]}`)
}
func (e *echoTool) Execute(ctx context.Context, args json.RawMessage) (tools.Result, error) {
	return tools.Result{Content: string(args)}, nil
}

func newTestDeps(t *testing.T) *app.Deps {
	t.Helper()
	s, err := store.Open(t.TempDir())
	if err != nil {
		t.Fatalf("store: %v", err)
	}
	return &app.Deps{Config: &config.Config{}, Store: s}
}

func TestLoopSimpleText(t *testing.T) {
	d := newTestDeps(t)
	defer d.Store.Close()

	provider := &stubProvider{responses: []llm.Response{
		{Content: "hello"},
	}}
	registry := tools.NewRegistry(&echoTool{})
	loop := New(d, provider, registry)

	events := make(chan Event, 16)
	go loop.Run(context.Background(), nil, "hi", events)

	var texts []string
	var done bool
	for ev := range events {
		if ev.Err != nil {
			t.Fatalf("unexpected err: %v", ev.Err)
		}
		if ev.TextDelta != "" {
			texts = append(texts, ev.TextDelta)
		}
		if ev.Done {
			done = true
		}
	}
	if !done {
		t.Fatal("expected Done event")
	}
	if len(texts) == 0 || texts[0] != "hello" {
		t.Errorf("expected text 'hello', got %v", texts)
	}
}

func TestLoopToolCall(t *testing.T) {
	d := newTestDeps(t)
	defer d.Store.Close()

	provider := &stubProvider{responses: []llm.Response{
		{ToolCalls: []llm.ToolCall{{ID: "t1", Name: "echo", Arguments: json.RawMessage(`{"msg":"hi"}`)}}},
		{Content: "done"},
	}}
	registry := tools.NewRegistry(&echoTool{})
	loop := New(d, provider, registry)
	loop.AutoConfirm = true

	events := make(chan Event, 16)
	go loop.Run(context.Background(), nil, "run echo", events)

	var sawToolCall bool
	var sawToolResult bool
	var sawDone bool
	for ev := range events {
		if ev.Err != nil {
			t.Fatalf("unexpected err: %v", ev.Err)
		}
		if ev.ToolCall != nil {
			sawToolCall = true
		}
		if ev.ToolResult != nil {
			sawToolResult = true
			if !ev.ToolResult.OK {
				t.Errorf("tool result not OK: %s", ev.ToolResult.Content)
			}
		}
		if ev.Done {
			sawDone = true
		}
	}
	if !sawToolCall || !sawToolResult || !sawDone {
		t.Errorf("missing events: toolCall=%v toolResult=%v done=%v", sawToolCall, sawToolResult, sawDone)
	}
}

func TestLoopDestructiveConfirmCancel(t *testing.T) {
	d := newTestDeps(t)
	defer d.Store.Close()

	registry := tools.NewRegistry(&destructiveEcho{})
	provider := &stubProvider{responses: []llm.Response{
		{ToolCalls: []llm.ToolCall{{ID: "t1", Name: "echo_danger", Arguments: json.RawMessage(`{}`)}}},
		{Content: "acknowledged cancel"},
	}}
	loop := New(d, provider, registry)

	events := make(chan Event, 16)
	go loop.Run(context.Background(), nil, "do it", events)

	cancelled := false
	for ev := range events {
		if ev.NeedConfirm != nil {
			ev.NeedConfirm.Reply <- false
		}
		if ev.ToolResult != nil && ev.ToolResult.Content == "cancelled by user" {
			cancelled = true
		}
	}
	if !cancelled {
		t.Error("expected cancellation")
	}
}

type destructiveEcho struct{ echoTool }

func (d *destructiveEcho) Name() string      { return "echo_danger" }
func (d *destructiveEcho) Destructive() bool { return true }
