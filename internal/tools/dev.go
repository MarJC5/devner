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
	return "Start the dev server (pnpm dev / npm run start) for a Node / Next / Astro / Vite project. Spawns the process in background inside the frankenphp container and lets Caddy reverse_proxy the project's URL to it. No-op if already running."
}
func (t *StartDevServer) Destructive() bool { return true }
func (t *StartDevServer) Schema() json.RawMessage {
	return json.RawMessage(`{
  "type":"object",
  "properties":{
    "project":{"type":"string"},
    "command":{"type":"string","description":"Optional override. Default: 'pnpm dev' for next/astro/vite, 'npm run start' for node."}
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
	if p.DevPort == 0 {
		return Result{Content: fmt.Sprintf("project %q is not a Node-like project (no dev port)", a.Project)}, fmt.Errorf("no dev port")
	}
	cmdStr := a.Command
	if cmdStr == "" {
		cmdStr = project.DefaultCommand(project.Type(p.Type), p.DevPort)
	}
	if err := t.D.DevServer.Start(ctx, a.Project, p.DevPort, cmdStr); err != nil {
		return Result{Content: "start failed: " + err.Error()}, err
	}
	return Result{
		Content: fmt.Sprintf("✓ dev server starting for %s on internal port %d. URL: https://%s. First build may take ~10s.",
			a.Project, p.DevPort, p.Domain),
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
	return "Show dev server status: running/stopped, PID, port, URL. Pass a project name or call with {} for all Node-like projects."
}
func (t *DevServerStatus) Destructive() bool { return false }
func (t *DevServerStatus) Schema() json.RawMessage {
	return json.RawMessage(`{
  "type":"object",
  "properties":{"project":{"type":"string","description":"Optional project name. Omit to check all Node-like projects."}},
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
			if p.DevPort > 0 {
				names = append(names, p.Name)
			}
		}
		if len(names) == 0 {
			return Result{Content: "no Node / Next / Astro / Vite projects in store"}, nil
		}
	}

	out := ""
	for _, n := range names {
		p, err := t.D.Store.GetProject(ctx, n)
		if err != nil {
			out += fmt.Sprintf("- %s: not found\n", n)
			continue
		}
		if p.DevPort == 0 {
			out += fmt.Sprintf("- %s: not a Node-like project\n", n)
			continue
		}
		st, err := t.D.DevServer.Status(ctx, n)
		state := "stopped"
		if err == nil && st.Running {
			state = fmt.Sprintf("running pid=%d", st.PID)
		}
		out += fmt.Sprintf("- %s: %s port=%d url=https://%s\n", n, state, p.DevPort, p.Domain)
	}
	return Result{Content: out}, nil
}
