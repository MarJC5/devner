// Package llm abstracts chat completion providers behind a single interface.
// Each provider (OpenAI-compat for Infomaniak / Ollama / Groq; Anthropic
// native) implements Provider.Chat, which accepts the normalized message
// history + tool schemas and returns either a text reply or tool calls.
//
// Streaming is delivered via a tea.Msg-friendly channel so the TUI can
// render deltas incrementally without blocking Update().
package llm

import (
	"context"
	"encoding/json"
)

// Role names are the OpenAI-style set. Anthropic provider maps them.
const (
	RoleSystem    = "system"
	RoleUser      = "user"
	RoleAssistant = "assistant"
	RoleTool      = "tool"
)

type Message struct {
	Role       string          `json:"role"`
	Content    string          `json:"content,omitempty"`
	ToolCalls  []ToolCall      `json:"tool_calls,omitempty"`
	ToolCallID string          `json:"tool_call_id,omitempty"` // only for Role=tool
	Name       string          `json:"name,omitempty"`         // tool name for Role=tool
}

type ToolCall struct {
	ID        string          `json:"id"`
	Name      string          `json:"name"`
	Arguments json.RawMessage `json:"arguments"`
}

// ToolSchema is the JSON-schema description sent to the model so it knows
// what tools are available. Compatible with OpenAI function calling and
// Anthropic tool_use blocks (provider adapters translate as needed).
type ToolSchema struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Parameters  json.RawMessage `json:"parameters"` // JSON Schema object
}

// Response is the normalized result of one provider turn. Exactly one of
// Content or ToolCalls is populated (OpenAI/Anthropic both enforce this).
type Response struct {
	Content   string
	ToolCalls []ToolCall
	// Raw is the provider-native response, kept for debugging/audit.
	Raw any
}

// StreamEvent is emitted by providers that support streaming. The agent
// loop forwards these to the TUI.
type StreamEvent struct {
	// TextDelta contains an incremental assistant text token.
	TextDelta string
	// ToolCall is set when a complete tool call becomes available.
	ToolCall *ToolCall
	// Done is set on the final event with the full reconstructed Response.
	Done *Response
	// Err is set on the final event if the stream errored.
	Err error
}

type Provider interface {
	// Name returns a short identifier, e.g. "infomaniak".
	Name() string
	// Chat issues a single non-streaming turn.
	Chat(ctx context.Context, msgs []Message, tools []ToolSchema) (*Response, error)
	// Stream issues a streaming turn. Closes the channel when done.
	// Implementations may return an unbuffered channel; callers must read
	// to completion. If streaming is not supported the implementation
	// falls back to a single event with the non-streaming Response.
	Stream(ctx context.Context, msgs []Message, tools []ToolSchema) (<-chan StreamEvent, error)
}
