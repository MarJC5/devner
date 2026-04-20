// Package project scaffolds new project skeletons (WordPress, Laravel,
// Node, Next.js, Astro) by running the appropriate installer inside the
// frankenphp container where composer, wp-cli, npm and pnpm are available.
package project

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type Type string

const (
	WordPress Type = "wordpress"
	Laravel   Type = "laravel"
	Node      Type = "node"
	NextJS    Type = "nextjs"
	Astro     Type = "astro"
)

var KnownTypes = []Type{WordPress, Laravel, Node, NextJS, Astro}

func ParseType(s string) (Type, error) {
	t := Type(strings.ToLower(s))
	for _, k := range KnownTypes {
		if t == k {
			return t, nil
		}
	}
	return "", fmt.Errorf("unknown project type %q (known: %v)", s, KnownTypes)
}

// DocRoot is the container-side subpath under /var/www/html/<name> that
// Caddy should serve as document root. Determined per project type.
func (t Type) DocRoot(name string) string {
	base := "/var/www/html/" + name
	switch t {
	case Laravel:
		return base + "/public"
	case WordPress:
		return base
	case Node, NextJS, Astro:
		// Node-based projects are served via dev server (proxied) or static
		// build dir. For V1 we serve static `dist`/`out` if present; reverse
		// proxy to dev server is a V2 improvement.
		return base
	}
	return base
}

var validProjectName = regexp.MustCompile(`^[a-z][a-z0-9-]{0,62}$`)

func ValidateName(name string) error {
	if !validProjectName.MatchString(name) {
		return fmt.Errorf("invalid project name %q: lowercase letters, digits, hyphens only; must start with a letter", name)
	}
	return nil
}

// Executor runs a command inside the frankenphp container. Injected so we
// can mock in tests and keep `project` decoupled from `runtime`.
type Executor interface {
	Exec(ctx context.Context, service string, args []string, interactive bool) error
}

type Service struct {
	Exec        Executor
	ProjectsDir string
}

func New(exec Executor, projectsDir string) *Service {
	return &Service{Exec: exec, ProjectsDir: projectsDir}
}

// Scaffold creates the project directory on the host (via container so
// ownership matches the container user) and runs the appropriate installer.
func (s *Service) Scaffold(ctx context.Context, t Type, name string) error {
	if err := ValidateName(name); err != nil {
		return err
	}
	projDir := filepath.Join(s.ProjectsDir, name)
	if _, err := os.Stat(projDir); err == nil {
		return fmt.Errorf("project %q already exists at %s", name, projDir)
	}

	switch t {
	case Laravel:
		return s.Exec.Exec(ctx, "frankenphp", []string{
			"bash", "-c",
			fmt.Sprintf("cd /var/www/html && composer create-project laravel/laravel %s", shellEscape(name)),
		}, false)
	case WordPress:
		return s.Exec.Exec(ctx, "frankenphp", []string{
			"bash", "-c",
			fmt.Sprintf("cd /var/www/html && mkdir -p %s && cd %s && wp core download --allow-root --locale=en_US", shellEscape(name), shellEscape(name)),
		}, false)
	case Node:
		return s.Exec.Exec(ctx, "frankenphp", []string{
			"bash", "-c",
			fmt.Sprintf("cd /var/www/html && mkdir -p %s && cd %s && npm init -y", shellEscape(name), shellEscape(name)),
		}, false)
	case NextJS:
		return s.Exec.Exec(ctx, "frankenphp", []string{
			"bash", "-c",
			fmt.Sprintf("cd /var/www/html && pnpm create next-app %s --ts --tailwind --app --no-src-dir --import-alias '@/*' --use-pnpm --yes", shellEscape(name)),
		}, false)
	case Astro:
		return s.Exec.Exec(ctx, "frankenphp", []string{
			"bash", "-c",
			fmt.Sprintf("cd /var/www/html && pnpm create astro@latest %s --template minimal --install --no-git --yes", shellEscape(name)),
		}, false)
	}
	return errors.New("unsupported project type")
}

// Remove deletes the project directory.
func (s *Service) Remove(ctx context.Context, name string) error {
	if err := ValidateName(name); err != nil {
		return err
	}
	projDir := filepath.Join(s.ProjectsDir, name)
	return os.RemoveAll(projDir)
}

func shellEscape(s string) string {
	return "'" + strings.ReplaceAll(s, "'", `'\''`) + "'"
}
