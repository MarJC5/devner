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

func (s *Store) migrate() error {
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
	for _, name := range names {
		data, err := fs.ReadFile(migrations.FS, name)
		if err != nil {
			return err
		}
		if _, err := s.db.Exec(string(data)); err != nil {
			return fmt.Errorf("migration %s: %w", name, err)
		}
	}
	return nil
}

// ----- Projects -----

func (s *Store) UpsertProject(ctx context.Context, p Project) error {
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO projects(name,type,db_engine,db_name,domain,path,created_at)
		VALUES(?,?,?,?,?,?,?)
		ON CONFLICT(name) DO UPDATE SET
			type=excluded.type,
			db_engine=excluded.db_engine,
			db_name=excluded.db_name,
			domain=excluded.domain,
			path=excluded.path
	`, p.Name, p.Type, p.DBEngine, p.DBName, p.Domain, p.Path, p.CreatedAt.Unix())
	return err
}

func (s *Store) DeleteProject(ctx context.Context, name string) error {
	_, err := s.db.ExecContext(ctx, `DELETE FROM projects WHERE name=?`, name)
	return err
}

func (s *Store) GetProject(ctx context.Context, name string) (*Project, error) {
	row := s.db.QueryRowContext(ctx, `SELECT name,type,db_engine,db_name,domain,path,created_at FROM projects WHERE name=?`, name)
	var p Project
	var ts int64
	if err := row.Scan(&p.Name, &p.Type, &p.DBEngine, &p.DBName, &p.Domain, &p.Path, &ts); err != nil {
		return nil, err
	}
	p.CreatedAt = time.Unix(ts, 0)
	return &p, nil
}

func (s *Store) ListProjects(ctx context.Context) ([]Project, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT name,type,db_engine,db_name,domain,path,created_at FROM projects ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Project
	for rows.Next() {
		var p Project
		var ts int64
		if err := rows.Scan(&p.Name, &p.Type, &p.DBEngine, &p.DBName, &p.Domain, &p.Path, &ts); err != nil {
			return nil, err
		}
		p.CreatedAt = time.Unix(ts, 0)
		out = append(out, p)
	}
	return out, rows.Err()
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
