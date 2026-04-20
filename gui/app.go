package main

import (
	"context"
	"fmt"

	"github.com/devner/devner/internal/app"
	"github.com/devner/devner/internal/store"
	wailsruntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

// App is the root object Wails binds to the frontend. Every exported
// method becomes callable from JS via generated TypeScript wrappers in
// frontend/wailsjs/go/main/App.ts. Fields are NOT exposed; only methods.
//
// Internally we wrap an *app.Deps bundle — the same one used by the CLI
// and TUI — so the GUI, CLI and agent all share one codepath for
// projects, databases, stacks, hosts, etc.
type App struct {
	ctx  context.Context
	deps *app.Deps
}

func NewApp() *App { return &App{} }

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	d, err := app.Build()
	if err != nil {
		// We deliberately don't crash the whole GUI on deps-build
		// failure (could be missing config, docker not running, etc.).
		// Methods below return the error so the frontend can surface it.
		fmt.Println("devner-gui: app.Build error:", err)
		return
	}
	a.deps = d
	a.startStackTicker(ctx)
}

func (a *App) shutdown(ctx context.Context) {
	if a.deps != nil {
		_ = a.deps.Close()
	}
}

// beforeClose — currently a no-op (returning false means "let Wails
// quit the app normally"). Kept as a hook for a future tray-based
// close-to-tray workflow once fyne.io/systray (or Wails v3's native
// tray) cooperates reliably with Wails v2's macOS run loop.
func (a *App) beforeClose(ctx context.Context) (preventClose bool) {
	_ = wailsruntime.ScreenGetAll // keep the import alive; remove if ever unused
	return false
}

// Greet is the scaffold method — kept for Step 1 smoke testing. The
// frontend calls App.Greet("Devner") on mount and renders the result.
// We'll remove it once the real screens land in Step 3.
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s — devner GUI is wired.", name)
}

// ListProjects exposes the existing store directly. Lightweight Step-1
// binding so we can prove end-to-end data flow (frontend mount → Go →
// SQLite → typed list back).
func (a *App) ListProjects() ([]store.Project, error) {
	if a.deps == nil {
		return nil, fmt.Errorf("app deps not initialized")
	}
	return a.deps.Store.ListProjects(a.ctx)
}
