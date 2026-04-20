package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/devner/devner/internal/app"
)

// runInProject is shared by every tool that shells into the frankenphp
// container for a given project. Captures combined output, truncates the
// tail to 4 KB to keep tool-result payloads below LLM context limits.
func runInProject(ctx context.Context, d *app.Deps, projectName, cmdLine string) (string, error) {
	p, err := d.Store.GetProject(ctx, projectName)
	if err != nil {
		return fmt.Sprintf("project %q not found", projectName), err
	}

	buf := &captureWriter{}
	oldOut, oldErr := d.Runtime.Stdout, d.Runtime.Stderr
	d.Runtime.Stdout, d.Runtime.Stderr = buf, buf
	defer func() { d.Runtime.Stdout, d.Runtime.Stderr = oldOut, oldErr }()

	full := fmt.Sprintf("cd /var/www/html/%s && %s", shellSafe(p.Name), cmdLine)
	if err := d.Runtime.Exec(ctx, "frankenphp", []string{"bash", "-lc", full}, false); err != nil {
		return fmt.Sprintf("exec failed: %s\n\n%s", err, truncateTail(buf.String(), 4000)), err
	}
	return truncateTail(buf.String(), 4000), nil
}

func truncateTail(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[len(s)-n:]
}

// ExecInProject is the generic shell escape hatch. Prefer the typed tools
// (Composer, NPM, WPCli, Artisan) when the operation fits one of them —
// they give the LLM clearer affordances and more precise schemas.
type ExecInProject struct{ D *app.Deps }

type execArgs struct {
	Project string `json:"project"`
	Command string `json:"command"`
}

func (t *ExecInProject) Name() string { return "exec_in_project" }
func (t *ExecInProject) Description() string {
	return "Run an arbitrary shell command inside the frankenphp container, cd'd into the project directory. Escape hatch — prefer the typed tools (composer, npm, wp_cli, artisan) when they fit."
}
func (t *ExecInProject) Destructive() bool { return true }
func (t *ExecInProject) Schema() json.RawMessage {
	return json.RawMessage(`{
  "type":"object",
  "properties":{
    "project":{"type":"string"},
    "command":{"type":"string","description":"Shell command. Example: 'mkdir public/uploads && chmod 0755 public/uploads'"}
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
	out, err := runInProject(ctx, t.D, a.Project, a.Command)
	return Result{Content: out}, err
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

// ---- Composer ----

type Composer struct{ D *app.Deps }

type composerArgs struct {
	Project string   `json:"project"`
	Args    []string `json:"args"`
}

func (t *Composer) Name() string { return "composer" }
func (t *Composer) Description() string {
	return "Run composer in the project. Pass the composer subcommand + args. Examples: ['install'], ['require','laravel/sanctum'], ['update'], ['dump-autoload'], ['show']."
}
func (t *Composer) Destructive() bool { return true }
func (t *Composer) Schema() json.RawMessage {
	return json.RawMessage(`{
  "type":"object",
  "properties":{
    "project":{"type":"string"},
    "args":{"type":"array","items":{"type":"string"},"description":"composer args, e.g. ['require','laravel/sanctum'] or ['install']"}
  },
  "required":["project","args"],
  "additionalProperties":false
}`)
}
func (t *Composer) Execute(ctx context.Context, raw json.RawMessage) (Result, error) {
	var a composerArgs
	if err := json.Unmarshal(raw, &a); err != nil {
		return Result{Content: "bad args: " + err.Error()}, err
	}
	out, err := runInProject(ctx, t.D, a.Project, "composer "+joinArgs(a.Args))
	return Result{Content: out}, err
}

// ---- NPM / pnpm / yarn ----

type NPM struct{ D *app.Deps }

type npmArgs struct {
	Project string   `json:"project"`
	Manager string   `json:"manager"`
	Args    []string `json:"args"`
}

func (t *NPM) Name() string { return "npm" }
func (t *NPM) Description() string {
	return "Run npm / pnpm / yarn in the project. Specify the manager. Examples: manager='pnpm' args=['install']; manager='npm' args=['run','build']; manager='yarn' args=['add','react']."
}
func (t *NPM) Destructive() bool { return true }
func (t *NPM) Schema() json.RawMessage {
	return json.RawMessage(`{
  "type":"object",
  "properties":{
    "project":{"type":"string"},
    "manager":{"type":"string","enum":["npm","pnpm","yarn"]},
    "args":{"type":"array","items":{"type":"string"}}
  },
  "required":["project","manager","args"],
  "additionalProperties":false
}`)
}
func (t *NPM) Execute(ctx context.Context, raw json.RawMessage) (Result, error) {
	var a npmArgs
	if err := json.Unmarshal(raw, &a); err != nil {
		return Result{Content: "bad args: " + err.Error()}, err
	}
	if a.Manager != "npm" && a.Manager != "pnpm" && a.Manager != "yarn" {
		return Result{Content: "manager must be npm|pnpm|yarn"}, fmt.Errorf("invalid manager")
	}
	out, err := runInProject(ctx, t.D, a.Project, a.Manager+" "+joinArgs(a.Args))
	return Result{Content: out}, err
}

// ---- WP-CLI ----

type WPCli struct{ D *app.Deps }

type wpArgs struct {
	Project string   `json:"project"`
	Args    []string `json:"args"`
}

func (t *WPCli) Name() string { return "wp_cli" }
func (t *WPCli) Description() string {
	return "Run WP-CLI on a WordPress project. Examples: args=['core','install','--url=example.localhost','--title=Test','--admin_user=admin','--admin_email=a@b.c']; args=['plugin','install','woocommerce','--activate']; args=['option','get','siteurl']."
}
func (t *WPCli) Destructive() bool { return true }
func (t *WPCli) Schema() json.RawMessage {
	return json.RawMessage(`{
  "type":"object",
  "properties":{
    "project":{"type":"string"},
    "args":{"type":"array","items":{"type":"string"}}
  },
  "required":["project","args"],
  "additionalProperties":false
}`)
}
func (t *WPCli) Execute(ctx context.Context, raw json.RawMessage) (Result, error) {
	var a wpArgs
	if err := json.Unmarshal(raw, &a); err != nil {
		return Result{Content: "bad args: " + err.Error()}, err
	}
	// --allow-root since the container runs as root by default.
	cmd := "wp --allow-root " + joinArgs(a.Args)
	out, err := runInProject(ctx, t.D, a.Project, cmd)
	return Result{Content: out}, err
}

// ---- Artisan ----

type Artisan struct{ D *app.Deps }

type artisanArgs struct {
	Project string   `json:"project"`
	Args    []string `json:"args"`
}

func (t *Artisan) Name() string { return "artisan" }
func (t *Artisan) Description() string {
	return "Run `php artisan` on a Laravel project. Examples: args=['migrate']; args=['make:model','Post','-mc']; args=['route:list']; args=['tinker']."
}
func (t *Artisan) Destructive() bool { return true }
func (t *Artisan) Schema() json.RawMessage {
	return json.RawMessage(`{
  "type":"object",
  "properties":{
    "project":{"type":"string"},
    "args":{"type":"array","items":{"type":"string"}}
  },
  "required":["project","args"],
  "additionalProperties":false
}`)
}
func (t *Artisan) Execute(ctx context.Context, raw json.RawMessage) (Result, error) {
	var a artisanArgs
	if err := json.Unmarshal(raw, &a); err != nil {
		return Result{Content: "bad args: " + err.Error()}, err
	}
	out, err := runInProject(ctx, t.D, a.Project, "php artisan "+joinArgs(a.Args))
	return Result{Content: out}, err
}

// joinArgs shell-escapes each element with single quotes and joins with spaces.
// The container runs bash -lc so we need to guard arg boundaries.
func joinArgs(args []string) string {
	parts := make([]string, 0, len(args))
	for _, a := range args {
		parts = append(parts, "'"+strings.ReplaceAll(a, "'", `'\''`)+"'")
	}
	return strings.Join(parts, " ")
}
