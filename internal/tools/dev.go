package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/devner/devner/internal/app"
	"github.com/devner/devner/internal/project"
)

// ---- start_dev_server ----

type StartDevServer struct{ D *app.Deps }

type startDevArgs struct {
	Project string `json:"project"`
	Command string `json:"command,omitempty"`
}

func (t *StartDevServer) Name() string { return "start_dev_server" }
func (t *StartDevServer) Description() string {
	return "Start the project's configured dev-time process. For Node/Next/Astro/Vite apps this is the HMR dev server (reverse-proxied by Caddy). For PHP+assets projects it's the asset watcher (vite build --watch, mix watch, etc.) running alongside FrankenPHP. No-op if already running."
}
func (t *StartDevServer) Destructive() bool { return true }
func (t *StartDevServer) Schema() json.RawMessage {
	return json.RawMessage(`{
  "type":"object",
  "properties":{
    "project":{"type":"string"},
    "command":{"type":"string","description":"Optional override. Default: 'pnpm dev' for server mode, 'pnpm watch' for watch mode."}
  },
  "required":["project"],
  "additionalProperties":false
}`)
}
func (t *StartDevServer) Execute(ctx context.Context, raw json.RawMessage) (Result, error) {
	var a startDevArgs
	if err := json.Unmarshal(raw, &a); err != nil {
		return Result{Content: "bad args: " + err.Error()}, err
	}
	p, err := t.D.Store.GetProject(ctx, a.Project)
	if err != nil {
		return Result{Content: fmt.Sprintf("project %q not found", a.Project)}, err
	}
	cmdStr := a.Command
	switch project.DevMode(p.DevMode) {
	case project.DevModeNone:
		return Result{Content: fmt.Sprintf("project %q has no dev process configured (library or plain static project)", a.Project)}, fmt.Errorf("no dev mode")
	case project.DevModeServer:
		if p.DevPort == 0 {
			return Result{Content: fmt.Sprintf("project %q is a dev-server project but has no allocated port — run `devner reconcile --apply`", a.Project)}, fmt.Errorf("no dev port")
		}
		if cmdStr == "" {
			cmdStr = project.DefaultCommand(project.Type(p.Type), p.DevPort)
		}
	case project.DevModeWatch:
		if cmdStr == "" {
			cmdStr = project.WatchCommand(p.DevCommand)
		}
	default:
		return Result{Content: fmt.Sprintf("project %q has unknown dev mode %q", a.Project, p.DevMode)}, fmt.Errorf("unknown dev mode")
	}
	if err := t.D.DevServer.Start(ctx, a.Project, p.DevPort, cmdStr); err != nil {
		return Result{Content: "start failed: " + err.Error()}, err
	}
	if p.DevMode == string(project.DevModeServer) {
		return Result{
			Content: fmt.Sprintf("✓ dev server starting for %s on internal port %d. URL: https://%s. First build may take ~10s.",
				a.Project, p.DevPort, p.Domain),
		}, nil
	}
	return Result{
		Content: fmt.Sprintf("✓ asset watcher starting for %s (cmd: %s). PHP serves https://%s as usual.",
			a.Project, cmdStr, p.Domain),
	}, nil
}

// ---- stop_dev_server ----

type StopDevServer struct{ D *app.Deps }

type stopDevArgs struct {
	Project string `json:"project"`
}

func (t *StopDevServer) Name() string { return "stop_dev_server" }
func (t *StopDevServer) Description() string {
	return "Stop the dev server for a Node-like project. Idempotent — no error if already stopped."
}
func (t *StopDevServer) Destructive() bool { return true }
func (t *StopDevServer) Schema() json.RawMessage {
	return json.RawMessage(`{
  "type":"object",
  "properties":{"project":{"type":"string"}},
  "required":["project"],
  "additionalProperties":false
}`)
}
func (t *StopDevServer) Execute(ctx context.Context, raw json.RawMessage) (Result, error) {
	var a stopDevArgs
	if err := json.Unmarshal(raw, &a); err != nil {
		return Result{Content: "bad args: " + err.Error()}, err
	}
	if _, err := t.D.Store.GetProject(ctx, a.Project); err != nil {
		return Result{Content: fmt.Sprintf("project %q not found", a.Project)}, err
	}
	if err := t.D.DevServer.Stop(ctx, a.Project); err != nil {
		return Result{Content: "stop failed: " + err.Error()}, err
	}
	return Result{Content: "✓ dev server stopped for " + a.Project}, nil
}

// ---- dev_server_status ----

type DevServerStatus struct{ D *app.Deps }

type devStatusArgs struct {
	Project string `json:"project,omitempty"`
}

func (t *DevServerStatus) Name() string { return "dev_server_status" }
func (t *DevServerStatus) Description() string {
	return "Show dev process status: running/stopped, PID, port (server mode) or watch command. Pass a project name or call with {} for all projects with a configured dev mode."
}
func (t *DevServerStatus) Destructive() bool { return false }
func (t *DevServerStatus) Schema() json.RawMessage {
	return json.RawMessage(`{
  "type":"object",
  "properties":{"project":{"type":"string","description":"Optional project name. Omit to check all projects with a dev mode."}},
  "additionalProperties":false
}`)
}
func (t *DevServerStatus) Execute(ctx context.Context, raw json.RawMessage) (Result, error) {
	var a devStatusArgs
	if len(raw) > 0 {
		if err := json.Unmarshal(raw, &a); err != nil {
			return Result{Content: "bad args: " + err.Error()}, err
		}
	}

	var names []string
	if a.Project != "" {
		names = []string{a.Project}
	} else {
		all, err := t.D.Store.ListProjects(ctx)
		if err != nil {
			return Result{Content: err.Error()}, err
		}
		for _, p := range all {
			if p.DevMode != "" {
				names = append(names, p.Name)
			}
		}
		if len(names) == 0 {
			return Result{Content: "no projects with a configured dev mode in store"}, nil
		}
	}

	out := ""
	for _, n := range names {
		p, err := t.D.Store.GetProject(ctx, n)
		if err != nil {
			out += fmt.Sprintf("- %s: not found\n", n)
			continue
		}
		if p.DevMode == "" {
			out += fmt.Sprintf("- %s: no dev process (library or plain PHP)\n", n)
			continue
		}
		st, err := t.D.DevServer.Status(ctx, n)
		state := "stopped"
		if err == nil && st.Running {
			state = fmt.Sprintf("running pid=%d", st.PID)
		}
		if p.DevMode == string(project.DevModeServer) {
			out += fmt.Sprintf("- %s: %s mode=server port=%d url=https://%s\n", n, state, p.DevPort, p.Domain)
		} else {
			out += fmt.Sprintf("- %s: %s mode=watch cmd=%s url=https://%s\n", n, state, p.DevCommand, p.Domain)
		}
	}
	return Result{Content: out}, nil
}
