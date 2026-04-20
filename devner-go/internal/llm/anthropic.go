// Anthropic native provider.
//
// We use the official anthropic-sdk-go. Anthropic's Messages API uses
// `tool_use` content blocks rather than OpenAI's `tool_calls` array, so
// we translate both ways here. System messages move from the `messages`
// array into the top-level `system` field.
package llm

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	anthropic "github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

type Anthropic struct {
	name   string
	model  string
	client *anthropic.Client
}

func NewAnthropic(name, apiKey, model string) *Anthropic {
	if apiKey == "" {
		apiKey = os.Getenv("ANTHROPIC_API_KEY")
	}
	c := anthropic.NewClient(option.WithAPIKey(apiKey))
	return &Anthropic{name: name, model: model, client: &c}
}

func (p *Anthropic) Name() string { return p.name }

func (p *Anthropic) Chat(ctx context.Context, msgs []Message, tools []ToolSchema) (*Response, error) {
	sys, userMsgs := splitSystem(msgs)
	apiMsgs, err := toAnthropicMessages(userMsgs)
	if err != nil {
		return nil, err
	}
	apiTools, err := toAnthropicTools(tools)
	if err != nil {
		return nil, err
	}

	params := anthropic.MessageNewParams{
		Model:     anthropic.Model(p.model),
		MaxTokens: 4096,
		Messages:  apiMsgs,
		Tools:     apiTools,
	}
	if sys != "" {
		params.System = []anthropic.TextBlockParam{{Text: sys}}
	}

	resp, err := p.client.Messages.New(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("%s chat: %w", p.name, err)
	}

	out := &Response{Raw: resp}
	for _, block := range resp.Content {
		switch b := block.AsAny().(type) {
		case anthropic.TextBlock:
			out.Content += b.Text
		case anthropic.ToolUseBlock:
			out.ToolCalls = append(out.ToolCalls, ToolCall{
				ID:        b.ID,
				Name:      b.Name,
				Arguments: b.Input,
			})
		}
	}
	return out, nil
}

// Stream falls back to non-streaming for V1 to keep the provider simple.
// Upgrading to real SSE streaming via p.client.Messages.NewStreaming is a
// mechanical change for V2.
func (p *Anthropic) Stream(ctx context.Context, msgs []Message, tools []ToolSchema) (<-chan StreamEvent, error) {
	out := make(chan StreamEvent, 1)
	go func() {
		defer close(out)
		resp, err := p.Chat(ctx, msgs, tools)
		if err != nil {
			out <- StreamEvent{Err: err}
			return
		}
		if resp.Content != "" {
			out <- StreamEvent{TextDelta: resp.Content}
		}
		out <- StreamEvent{Done: resp}
	}()
	return out, nil
}

func splitSystem(msgs []Message) (string, []Message) {
	var sys string
	var rest []Message
	for _, m := range msgs {
		if m.Role == RoleSystem {
			if sys != "" {
				sys += "\n\n"
			}
			sys += m.Content
			continue
		}
		rest = append(rest, m)
	}
	return sys, rest
}

func toAnthropicMessages(msgs []Message) ([]anthropic.MessageParam, error) {
	out := make([]anthropic.MessageParam, 0, len(msgs))
	for _, m := range msgs {
		switch m.Role {
		case RoleUser:
			out = append(out, anthropic.NewUserMessage(anthropic.NewTextBlock(m.Content)))
		case RoleAssistant:
			var blocks []anthropic.ContentBlockParamUnion
			if m.Content != "" {
				blocks = append(blocks, anthropic.NewTextBlock(m.Content))
			}
			for _, tc := range m.ToolCalls {
				var input any
				if err := json.Unmarshal(tc.Arguments, &input); err != nil {
					input = map[string]any{}
				}
				blocks = append(blocks, anthropic.NewToolUseBlock(tc.ID, input, tc.Name))
			}
			out = append(out, anthropic.NewAssistantMessage(blocks...))
		case RoleTool:
			out = append(out, anthropic.NewUserMessage(
				anthropic.NewToolResultBlock(m.ToolCallID, m.Content, false),
			))
		default:
			return nil, fmt.Errorf("anthropic: unsupported role %q", m.Role)
		}
	}
	return out, nil
}

func toAnthropicTools(tools []ToolSchema) ([]anthropic.ToolUnionParam, error) {
	if len(tools) == 0 {
		return nil, nil
	}
	out := make([]anthropic.ToolUnionParam, 0, len(tools))
	for _, t := range tools {
		var params anthropic.ToolInputSchemaParam
		if err := json.Unmarshal(t.Parameters, &params); err != nil {
			return nil, fmt.Errorf("anthropic tool %s schema: %w", t.Name, err)
		}
		out = append(out, anthropic.ToolUnionParam{
			OfTool: &anthropic.ToolParam{
				Name:        t.Name,
				Description: anthropic.String(t.Description),
				InputSchema: params,
			},
		})
	}
	return out, nil
}

var _ = errors.New
