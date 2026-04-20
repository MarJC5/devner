package cli

import (
	"github.com/devner/devner/internal/app"
)

// Deps is an alias for app.Deps retained here only so the cli package's
// internal commands keep a short name. New code should reference app.Deps
// directly.
type Deps = app.Deps

func buildDeps() (*Deps, error) { return app.Build() }
