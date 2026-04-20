package tools

import (
	"context"
	"encoding/json"

	"github.com/devner/devner/internal/app"
)

type StartStack struct{ D *app.Deps }

func (t *StartStack) Name() string             { return "start_stack" }
func (t *StartStack) Description() string      { return "Start the shared dev stack (mysql, postgres, frankenphp, redis, mailpit, adminer)." }
func (t *StartStack) Destructive() bool        { return false }
func (t *StartStack) Schema() json.RawMessage  { return emptySchema }
func (t *StartStack) Execute(ctx context.Context, _ json.RawMessage) (Result, error) {
	if err := t.D.Runtime.Up(ctx); err != nil {
		return Result{Content: "up failed: " + err.Error()}, err
	}
	return Result{Content: "✓ stack up"}, nil
}

type StopStack struct{ D *app.Deps }

func (t *StopStack) Name() string             { return "stop_stack" }
func (t *StopStack) Description() string      { return "Stop the shared dev stack (keeps volumes)." }
func (t *StopStack) Destructive() bool        { return false }
func (t *StopStack) Schema() json.RawMessage  { return emptySchema }
func (t *StopStack) Execute(ctx context.Context, _ json.RawMessage) (Result, error) {
	if err := t.D.Runtime.Stop(ctx); err != nil {
		return Result{Content: "stop failed: " + err.Error()}, err
	}
	return Result{Content: "✓ stack stopped"}, nil
}

type RebuildStack struct{ D *app.Deps }

func (t *RebuildStack) Name() string             { return "rebuild_stack" }
func (t *RebuildStack) Description() string      { return "Rebuild images and recreate containers. Slightly destructive: short downtime, but no data loss." }
func (t *RebuildStack) Destructive() bool        { return true }
func (t *RebuildStack) Schema() json.RawMessage  { return emptySchema }
func (t *RebuildStack) Execute(ctx context.Context, _ json.RawMessage) (Result, error) {
	if err := t.D.Runtime.Rebuild(ctx); err != nil {
		return Result{Content: "rebuild failed: " + err.Error()}, err
	}
	return Result{Content: "✓ stack rebuilt"}, nil
}

type TailLogs struct{ D *app.Deps }

type tailLogsArgs struct {
	Service string `json:"service"`
	Lines   int    `json:"lines"`
}

func (t *TailLogs) Name() string        { return "tail_logs" }
func (t *TailLogs) Description() string { return "Return the last N log lines of a stack service (mysql, postgres, frankenphp, redis, mailpit, adminer)." }
func (t *TailLogs) Destructive() bool   { return false }
func (t *TailLogs) Schema() json.RawMessage {
	return json.RawMessage(`{
  "type":"object",
  "properties":{
    "service":{"type":"string","enum":["frankenphp","mysql","postgres","redis","mailpit","adminer"]},
    "lines":{"type":"integer","minimum":1,"maximum":1000,"default":100}
  },
  "required":["service"],
  "additionalProperties":false
}`)
}
func (t *TailLogs) Execute(ctx context.Context, raw json.RawMessage) (Result, error) {
	var a tailLogsArgs
	if err := json.Unmarshal(raw, &a); err != nil {
		return Result{Content: "bad args: " + err.Error()}, err
	}
	if a.Lines <= 0 {
		a.Lines = 100
	}
	// Capture logs output to a buffer via Runtime.Stdout swap.
	buf := &captureWriter{}
	oldStdout := t.D.Runtime.Stdout
	oldStderr := t.D.Runtime.Stderr
	t.D.Runtime.Stdout = buf
	t.D.Runtime.Stderr = buf
	defer func() {
		t.D.Runtime.Stdout = oldStdout
		t.D.Runtime.Stderr = oldStderr
	}()
	_ = t.D.Runtime.Logs(ctx, a.Service, a.Lines, false)
	out := buf.String()
	if len(out) > 4000 {
		out = out[len(out)-4000:]
	}
	return Result{Content: out}, nil
}

type captureWriter struct{ buf []byte }

func (c *captureWriter) Write(p []byte) (int, error) { c.buf = append(c.buf, p...); return len(p), nil }
func (c *captureWriter) String() string               { return string(c.buf) }

var emptySchema = json.RawMessage(`{"type":"object","properties":{},"additionalProperties":false}`)
