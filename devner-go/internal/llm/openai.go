// OpenAI-compatible provider. Drives Infomaniak AI Tools, Ollama (via its
// /v1 endpoint), and Groq. Any provider implementing the OpenAI Chat
// Completions spec + function calling works here.
package llm

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	openai "github.com/sashabaranov/go-openai"
)

type OpenAICompat struct {
	name   string
	model  string
	client *openai.Client
}

func NewOpenAICompat(name, baseURL, apiKey, model string) *OpenAICompat {
	cfg := openai.DefaultConfig(apiKey)
	if baseURL != "" {
		cfg.BaseURL = baseURL
	}
	return &OpenAICompat{
		name:   name,
		model:  model,
		client: openai.NewClientWithConfig(cfg),
	}
}

func (p *OpenAICompat) Name() string { return p.name }

func (p *OpenAICompat) Chat(ctx context.Context, msgs []Message, tools []ToolSchema) (*Response, error) {
	req := p.buildRequest(msgs, tools, false)
	resp, err := p.client.CreateChatCompletion(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("%s chat: %w", p.name, err)
	}
	if len(resp.Choices) == 0 {
		return nil, errors.New("no choices returned")
	}
	choice := resp.Choices[0]

	out := &Response{Raw: resp}
	if len(choice.Message.ToolCalls) > 0 {
		for _, tc := range choice.Message.ToolCalls {
			out.ToolCalls = append(out.ToolCalls, ToolCall{
				ID:        tc.ID,
				Name:      tc.Function.Name,
				Arguments: json.RawMessage(tc.Function.Arguments),
			})
		}
	} else {
		out.Content = choice.Message.Content
	}
	return out, nil
}

func (p *OpenAICompat) Stream(ctx context.Context, msgs []Message, tools []ToolSchema) (<-chan StreamEvent, error) {
	req := p.buildRequest(msgs, tools, true)
	stream, err := p.client.CreateChatCompletionStream(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("%s stream: %w", p.name, err)
	}

	out := make(chan StreamEvent, 8)
	go func() {
		defer close(out)
		defer stream.Close()

		var textBuf string
		// tool_calls arrive as deltas indexed by their position; reconstruct.
		tcBuf := map[int]*partialToolCall{}

		for {
			chunk, err := stream.Recv()
			if err != nil {
				// io.EOF signals end of stream.
				if err.Error() == "EOF" || errors.Is(err, context.Canceled) {
					final := &Response{Content: textBuf}
					for _, pt := range tcBuf {
						final.ToolCalls = append(final.ToolCalls, pt.finalize())
					}
					out <- StreamEvent{Done: final}
					return
				}
				out <- StreamEvent{Err: err}
				return
			}
			if len(chunk.Choices) == 0 {
				continue
			}
			delta := chunk.Choices[0].Delta
			if delta.Content != "" {
				textBuf += delta.Content
				out <- StreamEvent{TextDelta: delta.Content}
			}
			for _, tc := range delta.ToolCalls {
				idx := 0
				if tc.Index != nil {
					idx = *tc.Index
				}
				cur := tcBuf[idx]
				if cur == nil {
					cur = &partialToolCall{}
					tcBuf[idx] = cur
				}
				if tc.ID != "" {
					cur.id = tc.ID
				}
				if tc.Function.Name != "" {
					cur.name = tc.Function.Name
				}
				if tc.Function.Arguments != "" {
					cur.args += tc.Function.Arguments
				}
			}
		}
	}()
	return out, nil
}

type partialToolCall struct {
	id, name, args string
}

func (p *partialToolCall) finalize() ToolCall {
	args := p.args
	if args == "" {
		args = "{}"
	}
	return ToolCall{ID: p.id, Name: p.name, Arguments: json.RawMessage(args)}
}

func (p *OpenAICompat) buildRequest(msgs []Message, tools []ToolSchema, stream bool) openai.ChatCompletionRequest {
	oaiMsgs := make([]openai.ChatCompletionMessage, 0, len(msgs))
	for _, m := range msgs {
		om := openai.ChatCompletionMessage{
			Role:       m.Role,
			Content:    m.Content,
			Name:       m.Name,
			ToolCallID: m.ToolCallID,
		}
		for _, tc := range m.ToolCalls {
			om.ToolCalls = append(om.ToolCalls, openai.ToolCall{
				ID:   tc.ID,
				Type: openai.ToolTypeFunction,
				Function: openai.FunctionCall{
					Name:      tc.Name,
					Arguments: string(tc.Arguments),
				},
			})
		}
		oaiMsgs = append(oaiMsgs, om)
	}

	oaiTools := make([]openai.Tool, 0, len(tools))
	for _, t := range tools {
		var schema any
		_ = json.Unmarshal(t.Parameters, &schema)
		oaiTools = append(oaiTools, openai.Tool{
			Type: openai.ToolTypeFunction,
			Function: &openai.FunctionDefinition{
				Name:        t.Name,
				Description: t.Description,
				Parameters:  schema,
			},
		})
	}

	return openai.ChatCompletionRequest{
		Model:    p.model,
		Messages: oaiMsgs,
		Tools:    oaiTools,
		Stream:   stream,
	}
}
