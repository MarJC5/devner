package store

import (
	"context"
	"testing"
	"time"
)

func TestProjectsCRUD(t *testing.T) {
	tmp := t.TempDir()
	s, err := Open(tmp)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer s.Close()

	ctx := context.Background()
	p := Project{
		Name:      "myapp",
		Type:      "laravel",
		DBEngine:  "mysql",
		DBName:    "myapp",
		Domain:    "myapp.localhost",
		Path:      "/tmp/myapp",
		CreatedAt: time.Unix(1700000000, 0),
	}
	if err := s.UpsertProject(ctx, p); err != nil {
		t.Fatalf("upsert: %v", err)
	}

	got, err := s.GetProject(ctx, "myapp")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.Name != p.Name || got.Type != p.Type || got.Domain != p.Domain {
		t.Errorf("got %+v want %+v", got, p)
	}

	list, err := s.ListProjects(ctx)
	if err != nil || len(list) != 1 {
		t.Fatalf("list: err=%v len=%d", err, len(list))
	}

	if err := s.DeleteProject(ctx, "myapp"); err != nil {
		t.Fatalf("delete: %v", err)
	}
	list, _ = s.ListProjects(ctx)
	if len(list) != 0 {
		t.Errorf("expected empty list, got %d", len(list))
	}
}

func TestHistoryAppend(t *testing.T) {
	tmp := t.TempDir()
	s, err := Open(tmp)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer s.Close()

	ctx := context.Background()
	for i := 0; i < 3; i++ {
		if err := s.RecordHistory(ctx, HistoryEntry{
			TS: time.Now(), Tool: "create_project", ArgsJSON: "{}", ResultJSON: "{}", DurationMs: 42, OK: true,
		}); err != nil {
			t.Fatalf("record: %v", err)
		}
	}
	list, err := s.RecentHistory(ctx, 10)
	if err != nil {
		t.Fatalf("recent: %v", err)
	}
	if len(list) != 3 {
		t.Errorf("expected 3 entries, got %d", len(list))
	}
}
