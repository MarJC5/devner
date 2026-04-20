// Provisioning orchestrates the multi-step creation (scaffold + DB +
// post-setup + store + Caddy) so both the CLI `new` command and the
// agent `create_project` tool share one path.
package app

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/devner/devner/internal/database"
	"github.com/devner/devner/internal/network"
	"github.com/devner/devner/internal/project"
	"github.com/devner/devner/internal/store"
)

type CreateProjectRequest struct {
	Name       string
	Type       project.Type
	DBEngine   string // "" | "mysql" | "postgres"
	WPInstall  bool   // only meaningful when Type == WordPress
	WPTitle    string
	WPAdmin    string
	WPPassword string
	WPEmail    string
}

type CreateProjectResult struct {
	Project store.Project
	DBCreds *database.Credentials // nil when no DB requested
}

// CreateProject runs: scaffold → create DB (optional) → post-setup
// (Laravel .env + key:generate or WordPress wp-config + optional install)
// → store upsert → Caddy reload. Idempotent on the store + Caddy pieces
// but scaffold will fail if the project directory already exists.
func (d *Deps) CreateProject(ctx context.Context, req CreateProjectRequest) (*CreateProjectResult, error) {
	if err := project.ValidateName(req.Name); err != nil {
		return nil, err
	}

	if err := d.Project.Scaffold(ctx, req.Type, req.Name); err != nil {
		return nil, fmt.Errorf("scaffold: %w", err)
	}

	var creds *database.Credentials
	if req.DBEngine != "" {
		if err := database.ValidateName(req.Name); err != nil {
			return nil, fmt.Errorf("project name cannot be used as db name: %w", err)
		}
		c, err := d.DB.Create(ctx, database.Engine(req.DBEngine), req.Name)
		if err != nil {
			return nil, fmt.Errorf("db: %w", err)
		}
		creds = &c
	}

	// Post-setup writes framework config (.env / wp-config) so the
	// project is ready to serve without manual editing.
	var pCreds *project.DBCreds
	if creds != nil {
		pCreds = &project.DBCreds{
			Engine:   req.DBEngine,
			Host:     creds.Host,
			Port:     creds.Port,
			Database: creds.Database,
			User:     creds.User,
			Password: creds.Password,
		}
	}
	opts := project.PostSetupOptions{
		WPInstall:    req.WPInstall,
		WPSiteTitle:  req.WPTitle,
		WPAdminUser:  req.WPAdmin,
		WPAdminPass:  req.WPPassword,
		WPAdminEmail: req.WPEmail,
	}
	if _, err := d.Project.PostSetup(ctx, req.Type, req.Name, pCreds, opts); err != nil {
		return nil, fmt.Errorf("post-setup: %w", err)
	}

	domain := req.Name + ".localhost"
	dbName := ""
	if creds != nil {
		dbName = creds.Database
	}

	// Node-like projects need a stable internal port for Caddy to
	// reverse_proxy to. Allocated once at creation, persisted, never
	// reassigned. PHP projects get 0 (no dev server).
	devPort := 0
	if isNodeLike(req.Type) {
		p, err := d.Store.AllocateDevPort(ctx)
		if err != nil {
			return nil, fmt.Errorf("allocate dev port: %w", err)
		}
		devPort = p
	}

	storedProject := store.Project{
		Name:      req.Name,
		Type:      string(req.Type),
		DBEngine:  req.DBEngine,
		DBName:    dbName,
		Domain:    domain,
		Path:      filepath.Join(d.Config.Stack.ProjectsDir, req.Name),
		CreatedAt: time.Now(),
		DevPort:   devPort,
	}
	if err := d.Store.UpsertProject(ctx, storedProject); err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	if err := d.ApplyCaddy(ctx); err != nil {
		return nil, fmt.Errorf("caddy: %w", err)
	}

	return &CreateProjectResult{Project: storedProject, DBCreds: creds}, nil
}

// ApplyCaddy rebuilds the Caddy config from the current project list.
// Exposed so the CLI `reconcile` and the tools can share it too.
func (d *Deps) ApplyCaddy(ctx context.Context) error {
	projects, err := d.Store.ListProjects(ctx)
	if err != nil {
		return err
	}
	sites := make([]network.ProjectSite, 0, len(projects))
	for _, p := range projects {
		sites = append(sites, network.ProjectSite{
			Domain:  p.Domain,
			Root:    project.Type(p.Type).DocRoot(p.Name),
			DevPort: p.DevPort,
		})
	}
	return d.Caddy.Apply(ctx, sites)
}

// isNodeLike returns true for project types served by a JS runtime dev
// server rather than by PHP. Keeps the "which types need a dev port"
// rule in one place.
func isNodeLike(t project.Type) bool {
	switch t {
	case project.Node, project.NextJS, project.Astro, project.Vite:
		return true
	}
	return false
}
