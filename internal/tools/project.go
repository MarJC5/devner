package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"time"

	"github.com/devner/devner/internal/app"
	"github.com/devner/devner/internal/database"
	"github.com/devner/devner/internal/network"
	"github.com/devner/devner/internal/project"
	"github.com/devner/devner/internal/store"
)

// applyCaddy rebuilds Caddy site config from the current project list.
// Shared by several tools so they all stay in sync with the store.
func applyCaddy(ctx context.Context, d *app.Deps) error {
	projects, err := d.Store.ListProjects(ctx)
	if err != nil {
		return err
	}
	sites := make([]network.ProjectSite, 0, len(projects))
	for _, p := range projects {
		sites = append(sites, network.ProjectSite{
			Domain: p.Domain,
			Root:   project.Type(p.Type).DocRoot(p.Name),
		})
	}
	return d.Caddy.Apply(ctx, sites)
}

// ---- list_projects ----

type ListProjects struct{ D *app.Deps }

func (t *ListProjects) Name() string        { return "list_projects" }
func (t *ListProjects) Description() string { return "List all devner-managed projects with type, database, and URL." }
func (t *ListProjects) Destructive() bool   { return false }
func (t *ListProjects) Schema() json.RawMessage {
	return schema(`{"type":"object","properties":{},"additionalProperties":false}`)
}
func (t *ListProjects) Execute(ctx context.Context, _ json.RawMessage) (Result, error) {
	ps, err := t.D.Store.ListProjects(ctx)
	if err != nil {
		return Result{Content: "error: " + err.Error()}, err
	}
	if len(ps) == 0 {
		return Result{Content: "No projects yet."}, nil
	}
	s := ""
	for _, p := range ps {
		db := "-"
		if p.DBEngine != "" {
			db = p.DBEngine + "/" + p.DBName
		}
		s += fmt.Sprintf("- %s  type=%s  db=%s  url=https://%s\n", p.Name, p.Type, db, p.Domain)
	}
	return Result{Content: s, Data: ps}, nil
}

// ---- create_project ----

type CreateProject struct{ D *app.Deps }

type createProjectArgs struct {
	Name   string `json:"name"`
	Type   string `json:"type"`
	DB     string `json:"db,omitempty"`
}

func (t *CreateProject) Name() string { return "create_project" }
func (t *CreateProject) Description() string {
	return "Scaffold a new project (wordpress|laravel|node|nextjs|astro) and optionally create a database. The project becomes available at https://<name>.localhost."
}
func (t *CreateProject) Destructive() bool { return false }
func (t *CreateProject) Schema() json.RawMessage {
	return schema(`{
  "type":"object",
  "properties":{
    "name":{"type":"string","description":"Project name. Lowercase letters, digits, hyphens. Must start with a letter."},
    "type":{"type":"string","enum":["wordpress","laravel","node","nextjs","astro"]},
    "db":{"type":"string","enum":["mysql","postgres",""],"description":"Optional database to create alongside the project."}
  },
  "required":["name","type"],
  "additionalProperties":false
}`)
}
func (t *CreateProject) Execute(ctx context.Context, raw json.RawMessage) (Result, error) {
	var a createProjectArgs
	if err := json.Unmarshal(raw, &a); err != nil {
		return Result{Content: "bad args: " + err.Error()}, err
	}
	typ, err := project.ParseType(a.Type)
	if err != nil {
		return Result{Content: err.Error()}, err
	}
	if err := project.ValidateName(a.Name); err != nil {
		return Result{Content: err.Error()}, err
	}

	if err := t.D.Project.Scaffold(ctx, typ, a.Name); err != nil {
		return Result{Content: "scaffold failed: " + err.Error()}, err
	}

	dbName := ""
	if a.DB != "" {
		if err := database.ValidateName(a.Name); err != nil {
			return Result{Content: "project name invalid as db name: " + err.Error()}, err
		}
		creds, err := t.D.DB.Create(ctx, database.Engine(a.DB), a.Name)
		if err != nil {
			return Result{Content: "db create failed: " + err.Error()}, err
		}
		dbName = creds.Database
	}

	domain := a.Name + ".localhost"
	p := store.Project{
		Name:      a.Name,
		Type:      string(typ),
		DBEngine:  a.DB,
		DBName:    dbName,
		Domain:    domain,
		Path:      filepath.Join(t.D.Config.Stack.ProjectsDir, a.Name),
		CreatedAt: time.Now(),
	}
	if err := t.D.Store.UpsertProject(ctx, p); err != nil {
		return Result{Content: "store failed: " + err.Error()}, err
	}
	if err := applyCaddy(ctx, t.D); err != nil {
		return Result{Content: "caddy apply failed: " + err.Error()}, err
	}
	return Result{Content: fmt.Sprintf("✓ created %s (%s). URL: https://%s", a.Name, a.Type, domain)}, nil
}

// ---- delete_project ----

type DeleteProject struct{ D *app.Deps }

type deleteProjectArgs struct {
	Name string `json:"name"`
}

func (t *DeleteProject) Name() string { return "delete_project" }
func (t *DeleteProject) Description() string {
	return "Delete a project: removes files, drops its database, and unregisters it from Caddy. DESTRUCTIVE."
}
func (t *DeleteProject) Destructive() bool { return true }
func (t *DeleteProject) Schema() json.RawMessage {
	return schema(`{
  "type":"object",
  "properties":{"name":{"type":"string"}},
  "required":["name"],
  "additionalProperties":false
}`)
}
func (t *DeleteProject) Execute(ctx context.Context, raw json.RawMessage) (Result, error) {
	var a deleteProjectArgs
	if err := json.Unmarshal(raw, &a); err != nil {
		return Result{Content: "bad args: " + err.Error()}, err
	}
	p, err := t.D.Store.GetProject(ctx, a.Name)
	if err != nil {
		return Result{Content: fmt.Sprintf("project %q not found", a.Name)}, err
	}
	if p.DBEngine != "" && p.DBName != "" {
		_ = t.D.DB.Drop(ctx, database.Engine(p.DBEngine), p.DBName)
	}
	_ = t.D.Project.Remove(ctx, a.Name)
	if err := t.D.Store.DeleteProject(ctx, a.Name); err != nil {
		return Result{Content: "store delete failed: " + err.Error()}, err
	}
	if err := applyCaddy(ctx, t.D); err != nil {
		return Result{Content: "caddy apply failed: " + err.Error()}, err
	}
	return Result{Content: "✓ deleted " + a.Name}, nil
}

// ---- project_status ----

type ProjectStatus struct{ D *app.Deps }

func (t *ProjectStatus) Name() string { return "project_status" }
func (t *ProjectStatus) Description() string {
	return "Get the URL, type, and database details for a project."
}
func (t *ProjectStatus) Destructive() bool { return false }
func (t *ProjectStatus) Schema() json.RawMessage {
	return schema(`{
  "type":"object",
  "properties":{"name":{"type":"string"}},
  "required":["name"],
  "additionalProperties":false
}`)
}
func (t *ProjectStatus) Execute(ctx context.Context, raw json.RawMessage) (Result, error) {
	var a deleteProjectArgs
	if err := json.Unmarshal(raw, &a); err != nil {
		return Result{Content: "bad args: " + err.Error()}, err
	}
	p, err := t.D.Store.GetProject(ctx, a.Name)
	if err != nil {
		return Result{Content: "not found"}, err
	}
	return Result{
		Content: fmt.Sprintf("name=%s type=%s db=%s/%s url=https://%s path=%s",
			p.Name, p.Type, p.DBEngine, p.DBName, p.Domain, p.Path),
		Data: p,
	}, nil
}
