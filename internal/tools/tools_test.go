package tools

import (
	"encoding/json"
	"testing"

	"github.com/devner/devner/internal/app"
	"github.com/devner/devner/internal/config"
)

// TestBuildDefault verifies every tool in the default registry has a valid
// JSON schema and a non-empty description. The Deps used here are a shell
// — we don't actually call the tools, only inspect metadata.
func TestBuildDefault(t *testing.T) {
	deps := &app.Deps{Config: &config.Config{}}
	r := BuildDefault(deps)
	if len(r.List()) < 10 {
		t.Fatalf("expected >= 10 tools, got %d", len(r.List()))
	}
	seen := map[string]bool{}
	for _, tool := range r.List() {
		if tool.Name() == "" {
			t.Errorf("tool with empty name: %T", tool)
		}
		if seen[tool.Name()] {
			t.Errorf("duplicate tool name %q", tool.Name())
		}
		seen[tool.Name()] = true

		if tool.Description() == "" {
			t.Errorf("tool %q has empty description", tool.Name())
		}

		var schema map[string]any
		if err := json.Unmarshal(tool.Schema(), &schema); err != nil {
			t.Errorf("tool %q has invalid JSON schema: %v", tool.Name(), err)
			continue
		}
		if schema["type"] != "object" {
			t.Errorf("tool %q schema root must be type=object, got %v", tool.Name(), schema["type"])
		}
	}
}

func TestRegistryDispatch(t *testing.T) {
	r := NewRegistry()
	if _, ok := r.Get("nope"); ok {
		t.Error("empty registry should not return anything")
	}
}
