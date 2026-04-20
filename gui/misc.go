package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/devner/devner/internal/database"
	"github.com/devner/devner/internal/network"
	"github.com/devner/devner/internal/store"
)

// Bindings for Logs / DBs / Hosts / Settings. Thin wrappers around
// app.Deps so Wails can expose them to the frontend.

// ---- Logs ----

// TailLogs returns the last N lines of a service's docker log as a
// single string. Used by the Logs screen's tabs.
func (a *App) TailLogs(service string, lines int) (string, error) {
	if lines <= 0 {
		lines = 200
	}
	ctx, cancel := context.WithTimeout(a.ctx, 5*time.Second)
	defer cancel()
	container := service + "_devner"
	out, err := exec.CommandContext(ctx, "docker", "logs",
		"--tail", fmt.Sprintf("%d", lines), container).CombinedOutput()
	if err != nil {
		return string(out), fmt.Errorf("docker logs %s: %w", container, err)
	}
	return string(out), nil
}

// ---- Databases ----

type DatabaseInfo struct {
	Engine string `json:"engine"`
	Name   string `json:"name"`
}

func (a *App) ListDatabases() ([]DatabaseInfo, error) {
	if a.deps == nil {
		return nil, fmt.Errorf("app deps not initialized")
	}
	var out []DatabaseInfo
	ctx, cancel := context.WithTimeout(a.ctx, 5*time.Second)
	defer cancel()

	// Query via docker exec rather than opening a second connection
	// from the GUI. Keeps the gui module free of direct DB deps that
	// aren't already transitively in internal/*.
	if rows, err := queryDockerSQL(ctx,
		"mysql_devner",
		[]string{"mysql", "-uroot", "-pdevner", "-Nse",
			"SELECT schema_name FROM information_schema.schemata ORDER BY schema_name;"},
	); err == nil {
		for _, name := range rows {
			if isMySQLSystem(name) {
				continue
			}
			out = append(out, DatabaseInfo{Engine: "mysql", Name: name})
		}
	}

	if rows, err := queryDockerSQL(ctx,
		"postgres_devner",
		[]string{"psql", "-U", "devner", "-d", "devner", "-Atc",
			"SELECT datname FROM pg_database ORDER BY datname;"},
	); err == nil {
		for _, name := range rows {
			if isPostgresSystem(name) {
				continue
			}
			out = append(out, DatabaseInfo{Engine: "postgres", Name: name})
		}
	}
	return out, nil
}

func queryDockerSQL(ctx context.Context, container string, cmd []string) ([]string, error) {
	args := append([]string{"exec", container}, cmd...)
	output, err := exec.CommandContext(ctx, "docker", args...).Output()
	if err != nil {
		return nil, err
	}
	var rows []string
	var start int
	s := string(output)
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			if line := s[start:i]; line != "" {
				rows = append(rows, line)
			}
			start = i + 1
		}
	}
	if start < len(s) {
		if line := s[start:]; line != "" {
			rows = append(rows, line)
		}
	}
	return rows, nil
}

func isMySQLSystem(name string) bool {
	switch name {
	case "information_schema", "mysql", "performance_schema", "sys":
		return true
	}
	return false
}

func isPostgresSystem(name string) bool {
	switch name {
	case "template0", "template1", "postgres":
		return true
	}
	return false
}

type CreateDatabaseResult struct {
	Engine   string `json:"engine"`
	Database string `json:"database"`
	User     string `json:"user"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
}

func (a *App) CreateDatabase(engine, name string) (*CreateDatabaseResult, error) {
	if a.deps == nil {
		return nil, fmt.Errorf("app deps not initialized")
	}
	creds, err := a.deps.DB.Create(a.ctx, database.Engine(engine), name)
	if err != nil {
		return nil, err
	}
	return &CreateDatabaseResult{
		Engine: engine, Database: creds.Database, User: creds.User,
		Password: creds.Password, Host: creds.Host, Port: creds.Port,
	}, nil
}

func (a *App) DropDatabase(engine, name string) error {
	if a.deps == nil {
		return fmt.Errorf("app deps not initialized")
	}
	return a.deps.DB.Drop(a.ctx, database.Engine(engine), name)
}

// ---- Hosts ----

type HostEntryDTO struct {
	Domain    string `json:"domain"`
	Target    string `json:"target"`
	ManagedAt int64  `json:"managed_at"` // unix seconds
}

func (a *App) ListHosts() ([]HostEntryDTO, error) {
	if a.deps == nil {
		return nil, fmt.Errorf("app deps not initialized")
	}
	entries, err := a.deps.Store.ListHosts(a.ctx)
	if err != nil {
		return nil, err
	}
	out := make([]HostEntryDTO, 0, len(entries))
	for _, e := range entries {
		out = append(out, HostEntryDTO{
			Domain: e.Domain, Target: e.Target,
			ManagedAt: e.ManagedAt.Unix(),
		})
	}
	return out, nil
}

func (a *App) AddHost(domain, target string) error {
	if a.deps == nil {
		return fmt.Errorf("app deps not initialized")
	}
	if target == "" {
		target = "127.0.0.1"
	}
	hm, err := network.NewHostsManager()
	if err != nil {
		return err
	}
	if _, err := hm.Add(domain, target); err != nil {
		return err
	}
	if err := hm.Save(); err != nil {
		return err
	}
	return a.deps.Store.UpsertHost(a.ctx, store.HostEntry{
		Domain: domain, Target: target, ManagedAt: time.Now(),
	})
}

func (a *App) RemoveHost(domain string) error {
	if a.deps == nil {
		return fmt.Errorf("app deps not initialized")
	}
	hm, err := network.NewHostsManager()
	if err != nil {
		return err
	}
	if _, err := hm.Remove(domain); err != nil {
		return err
	}
	if err := hm.Save(); err != nil {
		return err
	}
	return a.deps.Store.DeleteHost(a.ctx, domain)
}

// ---- Settings ----

type ConfigInfo struct {
	ConfigPath     string                 `json:"config_path"`
	DataDir        string                 `json:"data_dir"`
	ProjectsDir    string                 `json:"projects_dir"`
	ActiveProvider string                 `json:"active_provider"`
	Providers      map[string]ProviderDTO `json:"providers"`
}

type ProviderDTO struct {
	BaseURL   string `json:"base_url"`
	Model     string `json:"model"`
	Kind      string `json:"kind"`
	APIKeyEnv string `json:"api_key_env"`
	HasKey    bool   `json:"has_key"`
}

func (a *App) GetConfigInfo() (*ConfigInfo, error) {
	if a.deps == nil {
		return nil, fmt.Errorf("app deps not initialized")
	}
	c := a.deps.Config
	out := &ConfigInfo{
		ConfigPath:     c.ConfigPath(),
		DataDir:        c.Stack.DataDir,
		ProjectsDir:    c.Stack.ProjectsDir,
		ActiveProvider: c.LLM.ActiveProvider,
		Providers:      map[string]ProviderDTO{},
	}
	for name, p := range c.LLM.Providers {
		hasKey := p.APIKey != ""
		if !hasKey && p.APIKeyEnv != "" {
			// Quick "can we authenticate" probe without a network call.
			hasKey = os.Getenv(p.APIKeyEnv) != ""
		}
		out.Providers[name] = ProviderDTO{
			BaseURL:   p.BaseURL,
			Model:     p.Model,
			Kind:      p.Kind,
			APIKeyEnv: p.APIKeyEnv,
			HasKey:    hasKey,
		}
	}
	return out, nil
}

// OpenConfigFile opens the devner config.toml in the user's default
// editor via `open` / xdg-open / start.
func (a *App) OpenConfigFile() error {
	if a.deps == nil {
		return fmt.Errorf("app deps not initialized")
	}
	return openURL(a.deps.Config.ConfigPath())
}

func (a *App) CertsStatus() (network.CertsStatus, error) {
	return network.DetectCerts(a.ctx), nil
}

func (a *App) InstallCertsCA() error {
	return network.InstallRootCA(a.ctx)
}
