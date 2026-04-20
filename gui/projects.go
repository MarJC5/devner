package main

import (
	"fmt"
	"os/exec"
	"runtime"

	"github.com/devner/devner/internal/app"
	"github.com/devner/devner/internal/database"
	"github.com/devner/devner/internal/project"
)

// Per-project actions exposed to the frontend. Most wrap a single
// existing deps method — they exist so the UI has one clean call per
// button instead of composing low-level ops from JS.

// OpenInEditor launches `<editor> <project-path>` in a detached process
// using project.Open (already handles exec.Start + Process.Release).
// Empty editor string → project.DefaultEditor auto-detects by probing
// PATH (code, cursor, zed, subl, webstorm, …).
func (a *App) OpenInEditor(name, editor string) error {
	if a.deps == nil {
		return fmt.Errorf("app deps not initialized")
	}
	p, err := a.deps.Store.GetProject(a.ctx, name)
	if err != nil {
		return fmt.Errorf("project %q not found", name)
	}
	if editor == "" {
		e, err := project.DefaultEditor()
		if err != nil {
			return err
		}
		editor = e
	}
	return project.Open(editor, p.Path)
}

// ListEditors reports which of the known editors are actually on PATH
// inside the GUI's process. Used by the UI to grey out unavailable
// options instead of showing a failure when the user picks one.
func (a *App) ListEditors() []string {
	var out []string
	for _, e := range project.SupportedEditors {
		if _, err := exec.LookPath(e); err == nil {
			out = append(out, e)
		}
	}
	return out
}

// OpenProjectURL opens https://<name>.localhost in the user's default
// browser. Uses the OS-native open command (mac), xdg-open (linux),
// start (windows).
func (a *App) OpenProjectURL(name string) error {
	if a.deps == nil {
		return fmt.Errorf("app deps not initialized")
	}
	p, err := a.deps.Store.GetProject(a.ctx, name)
	if err != nil {
		return fmt.Errorf("project %q not found", name)
	}
	return openURL("https://" + p.Domain)
}

func openURL(url string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	default:
		cmd = exec.Command("xdg-open", url)
	}
	cmd.Stdout = nil
	cmd.Stderr = nil
	if err := cmd.Start(); err != nil {
		return err
	}
	_ = cmd.Process.Release()
	return nil
}

// DeleteProject removes the project entirely (files + DB + Caddy
// entry). Destructive — the frontend must show a ConfirmDialog before
// calling this.
func (a *App) DeleteProject(name string) error {
	if a.deps == nil {
		return fmt.Errorf("app deps not initialized")
	}
	p, err := a.deps.Store.GetProject(a.ctx, name)
	if err != nil {
		return fmt.Errorf("project %q not found", name)
	}
	// Mirror the CLI `devner remove` flow: drop DB, remove files,
	// delete from store, re-apply Caddy.
	if p.DBEngine != "" && p.DBName != "" {
		// Tolerate DB drop failures (engine might be down) so the
		// user can still unregister a project that's half-broken.
		_ = a.deps.DB.Drop(a.ctx, database.Engine(p.DBEngine), p.DBName)
	}
	_ = a.deps.Project.Remove(a.ctx, name)
	if err := a.deps.Store.DeleteProject(a.ctx, name); err != nil {
		return err
	}
	return a.deps.ApplyCaddy(a.ctx)
}

// ---- Dev server ----

// StartDevServer spawns `pnpm dev` (or user override) in the frankenphp
// container for a Node-like project. The project's allocated dev port
// is read from the store. Caddy reverse_proxy is already configured,
// so the moment the dev server binds the URL works.
func (a *App) StartDevServer(name, command string) error {
	if a.deps == nil {
		return fmt.Errorf("app deps not initialized")
	}
	p, err := a.deps.Store.GetProject(a.ctx, name)
	if err != nil {
		return fmt.Errorf("project %q not found", name)
	}
	if p.DevPort == 0 {
		return fmt.Errorf("project %q has no dev port (not a Node-like project)", name)
	}
	if command == "" {
		command = project.DefaultCommand(project.Type(p.Type), p.DevPort)
	}
	return a.deps.DevServer.Start(a.ctx, name, p.DevPort, command)
}

func (a *App) StopDevServer(name string) error {
	if a.deps == nil {
		return fmt.Errorf("app deps not initialized")
	}
	return a.deps.DevServer.Stop(a.ctx, name)
}

type DevServerStatusDTO struct {
	Name    string `json:"name"`
	Running bool   `json:"running"`
	PID     int    `json:"pid"`
	Port    int    `json:"port"`
}

// DevServerStatus is a flat DTO (not the internal Status struct) so
// the generated TS type is clean.
func (a *App) DevServerStatus(name string) (DevServerStatusDTO, error) {
	if a.deps == nil {
		return DevServerStatusDTO{}, fmt.Errorf("app deps not initialized")
	}
	p, err := a.deps.Store.GetProject(a.ctx, name)
	if err != nil {
		return DevServerStatusDTO{}, fmt.Errorf("project %q not found", name)
	}
	st, err := a.deps.DevServer.Status(a.ctx, name)
	if err != nil {
		return DevServerStatusDTO{}, err
	}
	return DevServerStatusDTO{
		Name: name, Running: st.Running, PID: st.PID, Port: p.DevPort,
	}, nil
}

// CreateProject proxies to the shared provision flow. Same args as
// the CLI / agent — streaming log output would be a Step-4 polish,
// for now it just blocks until the scaffold + DB + post-setup are
// done.
func (a *App) CreateProject(req app.CreateProjectRequest) (*app.CreateProjectResult, error) {
	if a.deps == nil {
		return nil, fmt.Errorf("app deps not initialized")
	}
	return a.deps.CreateProject(a.ctx, req)
}
