package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/devner/devner/internal/app"
	"github.com/devner/devner/internal/database"
	"github.com/devner/devner/internal/network"
	"github.com/devner/devner/internal/project"
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
	Name         string `json:"name"`
	Type         string `json:"type"`
	DB           string `json:"db,omitempty"`
	Template     string `json:"template,omitempty"`
	WPInstall    bool   `json:"wp_install,omitempty"`
	WPSiteTitle  string `json:"wp_title,omitempty"`
	WPAdminUser  string `json:"wp_admin,omitempty"`
	WPAdminPass  string `json:"wp_password,omitempty"`
	WPAdminEmail string `json:"wp_email,omitempty"`
}

func (t *CreateProject) Name() string { return "create_project" }
func (t *CreateProject) Description() string {
	return "Scaffold a new project and wire it up end-to-end: optional database + user, framework config (Laravel .env + key:generate, WordPress wp-config.php, optional wp core install), Caddy site + HTTPS. The project becomes reachable at https://<name>.localhost."
}
func (t *CreateProject) Destructive() bool { return false }
func (t *CreateProject) Schema() json.RawMessage {
	return schema(`{
  "type":"object",
  "properties":{
    "name":{"type":"string","description":"Project name. Lowercase letters, digits, hyphens. Must start with a letter."},
    "type":{"type":"string","enum":["wordpress","laravel","node","nextjs","nuxt","astro","sveltekit","vite"]},
    "db":{"type":"string","enum":["mysql","postgres",""],"description":"Optional database engine."},
    "template":{"type":"string","description":"Scaffold template. Vite: react-ts (default) | vue-ts | svelte-ts | solid-ts | preact-ts | qwik-ts | lit-ts | vanilla-ts. Astro: minimal (default) | basics | blog | portfolio | starlight. SvelteKit: skeleton (default) | minimal | demo. Ignored by other types."},
    "wp_install":{"type":"boolean","description":"WordPress only: run 'wp core install' after setup so the site is immediately usable."},
    "wp_title":{"type":"string"},
    "wp_admin":{"type":"string"},
    "wp_password":{"type":"string"},
    "wp_email":{"type":"string"}
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

	res, err := t.D.CreateProject(ctx, app.CreateProjectRequest{
		Name:       a.Name,
		Type:       typ,
		DBEngine:   a.DB,
		Template:   a.Template,
		WPInstall:  a.WPInstall,
		WPTitle:    a.WPSiteTitle,
		WPAdmin:    a.WPAdminUser,
		WPPassword: a.WPAdminPass,
		WPEmail:    a.WPAdminEmail,
	})
	if err != nil {
		return Result{Content: err.Error()}, err
	}
	msg := fmt.Sprintf("✓ created %s (%s). URL: https://%s", res.Project.Name, a.Type, res.Project.Domain)
	if res.DBCreds != nil {
		msg += fmt.Sprintf("\n  DB: %s/%s user=%s host=%s:%d", a.DB, res.DBCreds.Database, res.DBCreds.User, res.DBCreds.Host, res.DBCreds.Port)
	}
	if a.WPInstall {
		admin := a.WPAdminUser
		if admin == "" {
			admin = "admin"
		}
		pass := a.WPAdminPass
		if pass == "" {
			pass = "admin"
		}
		msg += fmt.Sprintf("\n  WP admin: %s / %s", admin, pass)
	}
	_ = res
	return Result{Content: msg}, nil
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
