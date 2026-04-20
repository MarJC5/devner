package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/devner/devner/internal/app"
)

type ExecInProject struct{ D *app.Deps }

type execArgs struct {
	Project string `json:"project"`
	Command string `json:"command"`
}

func (t *ExecInProject) Name() string { return "exec_in_project" }
func (t *ExecInProject) Description() string {
	return "Run a shell command inside the frankenphp container, cd'd into the project directory. Use for wp-cli, artisan, composer, npm, pnpm, yarn."
}
func (t *ExecInProject) Destructive() bool { return true } // arbitrary shell → always confirm
func (t *ExecInProject) Schema() json.RawMessage {
	return json.RawMessage(`{
  "type":"object",
  "properties":{
    "project":{"type":"string"},
    "command":{"type":"string","description":"Shell command to run. Examples: 'wp core install', 'composer require x/y', 'pnpm install', 'php artisan migrate'"}
  },
  "required":["project","command"],
  "additionalProperties":false
}`)
}
func (t *ExecInProject) Execute(ctx context.Context, raw json.RawMessage) (Result, error) {
	var a execArgs
	if err := json.Unmarshal(raw, &a); err != nil {
		return Result{Content: "bad args: " + err.Error()}, err
	}
	p, err := t.D.Store.GetProject(ctx, a.Project)
	if err != nil {
		return Result{Content: fmt.Sprintf("project %q not found", a.Project)}, err
	}

	buf := &captureWriter{}
	oldOut, oldErr := t.D.Runtime.Stdout, t.D.Runtime.Stderr
	t.D.Runtime.Stdout, t.D.Runtime.Stderr = buf, buf
	defer func() { t.D.Runtime.Stdout, t.D.Runtime.Stderr = oldOut, oldErr }()

	// Escape: single-quote each segment to prevent injection via project path.
	full := fmt.Sprintf("cd /var/www/html/%s && %s", shellSafe(p.Name), a.Command)
	if err := t.D.Runtime.Exec(ctx, "frankenphp", []string{"bash", "-lc", full}, false); err != nil {
		return Result{Content: fmt.Sprintf("exec failed: %s\n\n%s", err, buf.String())}, err
	}
	out := buf.String()
	if len(out) > 4000 {
		out = out[len(out)-4000:]
	}
	return Result{Content: out}, nil
}

// shellSafe is intentionally strict: project names are validated to match
// ^[a-z][a-z0-9-]+$, so this mostly guards against future callers.
func shellSafe(s string) string {
	return strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			return r
		}
		return -1
	}, s)
}
