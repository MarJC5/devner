package store

import (
	"context"
	"database/sql"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/devner/devner/migrations"
	_ "modernc.org/sqlite"
)

type Store struct {
	db *sql.DB
}

type Project struct {
	Name      string
	Type      string
	DBEngine  string
	DBName    string
	Domain    string
	Path      string
	CreatedAt time.Time
	// DevPort is the internal container port Caddy reverse-proxies to for
	// Node/Next/Astro/Vite projects. 0 means the project doesn't need one
	// (PHP, static only). Allocated by AllocateDevPort, persisted, and
	// never changes for the life of the project.
	DevPort int
}

type HostEntry struct {
	Domain    string
	Target    string
	ManagedAt time.Time
}

type HistoryEntry struct {
	ID         int64
	TS         time.Time
	Tool       string
	ArgsJSON   string
	ResultJSON string
	DurationMs int64
	OK         bool
}

func Open(dataDir string) (*Store, error) {
	if err := os.MkdirAll(dataDir, 0o755); err != nil {
		return nil, err
	}
	path := filepath.Join(dataDir, "store.db")
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}
	if _, err := db.Exec("PRAGMA foreign_keys=ON; PRAGMA journal_mode=WAL;"); err != nil {
		db.Close()
		return nil, err
	}
	s := &Store{db: db}
	if err := s.migrate(); err != nil {
		db.Close()
		return nil, err
	}
	return s, nil
}

func (s *Store) Close() error { return s.db.Close() }

// migrate applies every .sql file in /migrations that hasn't been run
// against this database yet. Each applied file is recorded in
// schema_migrations so re-running the process is a no-op. Migration
// files can therefore use non-idempotent statements like
// `ALTER TABLE ... ADD COLUMN`.
func (s *Store) migrate() error {
	if _, err := s.db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			name TEXT PRIMARY KEY,
			applied_at INTEGER NOT NULL
		)
	`); err != nil {
		return fmt.Errorf("create schema_migrations: %w", err)
	}

	entries, err := fs.ReadDir(migrations.FS, ".")
	if err != nil {
		return err
	}
	var names []string
	for _, e := range entries {
		if !e.IsDir() && filepath.Ext(e.Name()) == ".sql" {
			names = append(names, e.Name())
		}
	}
	sort.Strings(names)

	// Load the set of already-applied migrations.
	applied := make(map[string]struct{})
	rows, err := s.db.Query(`SELECT name FROM schema_migrations`)
	if err != nil {
		return fmt.Errorf("read schema_migrations: %w", err)
	}
	for rows.Next() {
		var n string
		if err := rows.Scan(&n); err != nil {
			rows.Close()
			return err
		}
		applied[n] = struct{}{}
	}
	rows.Close()

	for _, name := range names {
		if _, ok := applied[name]; ok {
			continue
		}
		data, err := fs.ReadFile(migrations.FS, name)
		if err != nil {
			return err
		}
		if _, err := s.db.Exec(string(data)); err != nil {
			return fmt.Errorf("migration %s: %w", name, err)
		}
		if _, err := s.db.Exec(`INSERT INTO schema_migrations(name, applied_at) VALUES(?, ?)`,
			name, time.Now().Unix()); err != nil {
			return fmt.Errorf("record migration %s: %w", name, err)
		}
	}
	return nil
}

// ----- Projects -----

func (s *Store) UpsertProject(ctx context.Context, p Project) error {
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO projects(name,type,db_engine,db_name,domain,path,created_at,dev_port)
		VALUES(?,?,?,?,?,?,?,?)
		ON CONFLICT(name) DO UPDATE SET
			type=excluded.type,
			db_engine=excluded.db_engine,
			db_name=excluded.db_name,
			domain=excluded.domain,
			path=excluded.path,
			dev_port=excluded.dev_port
	`, p.Name, p.Type, p.DBEngine, p.DBName, p.Domain, p.Path, p.CreatedAt.Unix(), p.DevPort)
	return err
}

func (s *Store) DeleteProject(ctx context.Context, name string) error {
	_, err := s.db.ExecContext(ctx, `DELETE FROM projects WHERE name=?`, name)
	return err
}

func (s *Store) GetProject(ctx context.Context, name string) (*Project, error) {
	row := s.db.QueryRowContext(ctx, `SELECT name,type,db_engine,db_name,domain,path,created_at,dev_port FROM projects WHERE name=?`, name)
	var p Project
	var ts int64
	if err := row.Scan(&p.Name, &p.Type, &p.DBEngine, &p.DBName, &p.Domain, &p.Path, &ts, &p.DevPort); err != nil {
		return nil, err
	}
	p.CreatedAt = time.Unix(ts, 0)
	return &p, nil
}

func (s *Store) ListProjects(ctx context.Context) ([]Project, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT name,type,db_engine,db_name,domain,path,created_at,dev_port FROM projects ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Project
	for rows.Next() {
		var p Project
		var ts int64
		if err := rows.Scan(&p.Name, &p.Type, &p.DBEngine, &p.DBName, &p.Domain, &p.Path, &ts, &p.DevPort); err != nil {
			return nil, err
		}
		p.CreatedAt = time.Unix(ts, 0)
		out = append(out, p)
	}
	return out, rows.Err()
}

// DevPortRange bounds the internal-only port space allocated to Node-like
// projects. 3100 sits above Next.js (3000), MySQL (3306), and stays clear
// of Vite (5173), Astro (4321), Postgres (5432), Redis (6379), Caddy
// admin (2019). 3999 is a hard stop — realistic even for prolific users.
const (
	DevPortMin = 3100
	DevPortMax = 3999
)

// AllocateDevPort returns the smallest free integer in [DevPortMin,
// DevPortMax] not already assigned to a project. Returns 0 + error if
// the range is exhausted. Callers persist the returned value on the
// Project; this function does NOT reserve — two concurrent callers
// could race and get the same port. Devner's CLI is single-user, so
// the race is acceptable; if it ever becomes an issue, wrap the
// allocation + UpsertProject in a transaction.
func (s *Store) AllocateDevPort(ctx context.Context) (int, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT dev_port FROM projects WHERE dev_port >= ? AND dev_port <= ? ORDER BY dev_port`,
		DevPortMin, DevPortMax)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	taken := make(map[int]struct{})
	for rows.Next() {
		var p int
		if err := rows.Scan(&p); err != nil {
			return 0, err
		}
		taken[p] = struct{}{}
	}
	if err := rows.Err(); err != nil {
		return 0, err
	}
	for p := DevPortMin; p <= DevPortMax; p++ {
		if _, used := taken[p]; !used {
			return p, nil
		}
	}
	return 0, fmt.Errorf("dev port range %d-%d exhausted (%d projects)", DevPortMin, DevPortMax, len(taken))
}

// ----- Hosts -----

func (s *Store) UpsertHost(ctx context.Context, h HostEntry) error {
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO hosts(domain,target,managed_at) VALUES(?,?,?)
		ON CONFLICT(domain) DO UPDATE SET target=excluded.target, managed_at=excluded.managed_at
	`, h.Domain, h.Target, h.ManagedAt.Unix())
	return err
}

func (s *Store) DeleteHost(ctx context.Context, domain string) error {
	_, err := s.db.ExecContext(ctx, `DELETE FROM hosts WHERE domain=?`, domain)
	return err
}

func (s *Store) ListHosts(ctx context.Context) ([]HostEntry, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT domain,target,managed_at FROM hosts ORDER BY domain`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []HostEntry
	for rows.Next() {
		var h HostEntry
		var ts int64
		if err := rows.Scan(&h.Domain, &h.Target, &ts); err != nil {
			return nil, err
		}
		h.ManagedAt = time.Unix(ts, 0)
		out = append(out, h)
	}
	return out, rows.Err()
}

// ----- History -----

func (s *Store) RecordHistory(ctx context.Context, h HistoryEntry) error {
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO history(ts,tool,args_json,result_json,duration_ms,ok)
		VALUES(?,?,?,?,?,?)
	`, h.TS.Unix(), h.Tool, h.ArgsJSON, h.ResultJSON, h.DurationMs, boolToInt(h.OK))
	return err
}

func (s *Store) RecentHistory(ctx context.Context, limit int) ([]HistoryEntry, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id,ts,tool,args_json,result_json,duration_ms,ok FROM history ORDER BY id DESC LIMIT ?`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []HistoryEntry
	for rows.Next() {
		var h HistoryEntry
		var ts int64
		var ok int
		if err := rows.Scan(&h.ID, &ts, &h.Tool, &h.ArgsJSON, &h.ResultJSON, &h.DurationMs, &ok); err != nil {
			return nil, err
		}
		h.TS = time.Unix(ts, 0)
		h.OK = ok == 1
		out = append(out, h)
	}
	return out, rows.Err()
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
