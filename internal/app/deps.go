// Package app wires the core services together into a Deps bundle used by
// both the CLI and the TUI. It lives in its own package so the TUI can
// depend on it without creating a cycle with the cli package.
package app

import (
	"fmt"

	"github.com/devner/devner/internal/config"
	"github.com/devner/devner/internal/database"
	"github.com/devner/devner/internal/network"
	"github.com/devner/devner/internal/project"
	"github.com/devner/devner/internal/runtime"
	"github.com/devner/devner/internal/store"
)

type Deps struct {
	Config    *config.Config
	Store     *store.Store
	Runtime   *runtime.Runtime
	DB        *database.Manager
	Project   *project.Service
	DevServer *project.DevServer
	Caddy     *network.CaddyClient
}

func Build() (*Deps, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("config: %w", err)
	}
	st, err := store.Open(cfg.Stack.DataDir)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}
	rt := runtime.New(cfg.Stack.DataDir, cfg.Stack.ProjectsDir)
	return &Deps{
		Config:    cfg,
		Store:     st,
		Runtime:   rt,
		DB:        database.NewManager(),
		Project:   project.New(rt, cfg.Stack.ProjectsDir),
		DevServer: project.NewDevServer(rt),
		Caddy:     network.NewCaddyClient(),
	}, nil
}

func (d *Deps) Close() error {
	if d.Store != nil {
		return d.Store.Close()
	}
	return nil
}
