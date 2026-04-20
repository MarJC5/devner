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
	PHP       Type = "php"
	Node      Type = "node"
	NextJS    Type = "nextjs"
	Nuxt      Type = "nuxt"
	Astro     Type = "astro"
	SvelteKit Type = "sveltekit"
	Vite      Type = "vite"
)

var KnownTypes = []Type{WordPress, Laravel, PHP, Node, NextJS, Nuxt, Astro, SvelteKit, Vite}

// DevMode describes what kind of dev-time process (if any) a project wants
// when `devner dev start` is invoked. It's stored alongside the project so
// detection runs once at import/create, not on every dev command.
type DevMode string

const (
	// DevModeServer: HTTP dev server with HMR. Requires a dev_port so
	// Caddy can reverse-proxy the project's URL to it (Next/Vite/Astro
	// apps).
	DevModeServer DevMode = "server"
	// DevModeWatch: asset-compile in watch mode alongside a PHP backend.
	// No port, no Caddy routing change — PHP still serves HTTP.
	DevModeWatch DevMode = "watch"
	// DevModeNone: no dev-time process. Libraries, plain PHP without an
	// asset pipeline, static sites.
	DevModeNone DevMode = ""
)

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
//
// For the generic PHP type we assume /public if the project follows the
// common "front controller in public/" layout (Symfony, Slim, most
// modern PHP CMSes). Callers that need absolute accuracy should inspect
// the FS and override.
func (t Type) DocRoot(name string) string {
	base := "/var/www/html/" + name
	switch t {
	case Laravel, PHP:
		return base + "/public"
	case WordPress:
		return base
	case Node, NextJS, Nuxt, Astro, SvelteKit, Vite:
		// Node-based projects go through Caddy reverse_proxy to their
		// dev server — DocRoot is unused (Caddy doesn't serve files for
		// these). Kept for consistency; returned path would only matter
		// if someone switches the project to static-serve mode.
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

// ScaffoldOptions carries per-type knobs that some scaffolders need
// (Vite's --template choice, for example). Zero value is fine for
// types that don't use any option.
type ScaffoldOptions struct {
	// Template selects the upstream template for types that expose one:
	//   - Vite: react-ts (default), vue-ts, svelte-ts, solid-ts,
	//     qwik-ts, preact-ts, lit-ts, vanilla-ts.
	//   - Astro: minimal (default), basics, blog, portfolio, starlight.
	//   - SvelteKit: skeleton (default), minimal, demo.
	// Ignored by types that don't use it (Laravel, WordPress, Next, …).
	Template string
}

// Scaffold creates the project directory on the host (via container so
// ownership matches the container user) and runs the appropriate installer.
func (s *Service) Scaffold(ctx context.Context, t Type, name string, opts ScaffoldOptions) error {
	if err := ValidateName(name); err != nil {
		return err
	}
	projDir := filepath.Join(s.ProjectsDir, name)
	if _, err := os.Stat(projDir); err == nil {
		return fmt.Errorf("project %q already exists at %s", name, projDir)
	}

	switch t {
	case PHP:
		return errors.New("generic PHP projects are import-only — drop your code in projects/<name> and run `devner import`")
	case Laravel:
		return s.Exec.Exec(ctx, "frankenphp", []string{
			"bash", "-lc",
			fmt.Sprintf("cd /var/www/html && composer create-project laravel/laravel %s", shellEscape(name)),
		}, false)
	case WordPress:
		return s.Exec.Exec(ctx, "frankenphp", []string{
			"bash", "-lc",
			fmt.Sprintf("cd /var/www/html && mkdir -p %s && cd %s && wp core download --allow-root --locale=en_US", shellEscape(name), shellEscape(name)),
		}, false)
	case Node:
		return s.Exec.Exec(ctx, "frankenphp", []string{
			"bash", "-lc",
			fmt.Sprintf("cd /var/www/html && mkdir -p %s && cd %s && npm init -y", shellEscape(name), shellEscape(name)),
		}, false)
	case NextJS:
		return s.Exec.Exec(ctx, "frankenphp", []string{
			"bash", "-lc",
			fmt.Sprintf("cd /var/www/html && pnpm create next-app %s --ts --tailwind --app --no-src-dir --import-alias '@/*' --use-pnpm --yes", shellEscape(name)),
		}, false)
	case Nuxt:
		// nuxi init is the canonical scaffolder. --force to skip the
		// "target dir exists" prompt (we've already checked),
		// --no-gitInit for consistency with other types, --packageManager=pnpm
		// so the lockfile matches what the container has. After init,
		// install deps so `pnpm dev` works immediately.
		return s.Exec.Exec(ctx, "frankenphp", []string{
			"bash", "-lc",
			fmt.Sprintf("cd /var/www/html && pnpm dlx nuxi@latest init %s --force --no-gitInit --packageManager pnpm && cd %s && pnpm install",
				shellEscape(name), shellEscape(name)),
		}, false)
	case Astro:
		template := opts.Template
		if template == "" {
			template = "minimal"
		}
		return s.Exec.Exec(ctx, "frankenphp", []string{
			"bash", "-lc",
			fmt.Sprintf("cd /var/www/html && pnpm create astro@latest %s --template %s --install --no-git --yes",
				shellEscape(name), shellEscape(template)),
		}, false)
	case SvelteKit:
		// create-svelte drives SvelteKit. `--template skeleton` is the
		// barest starting point; other options include `minimal` and
		// `demo`. `--types typescript` picks the TS variant without
		// interactive prompts. After scaffold, install deps.
		template := opts.Template
		if template == "" {
			template = "skeleton"
		}
		return s.Exec.Exec(ctx, "frankenphp", []string{
			"bash", "-lc",
			fmt.Sprintf("cd /var/www/html && pnpm create svelte@latest %s --template %s --types typescript --no-prettier --no-eslint --no-playwright --no-vitest && cd %s && pnpm install",
				shellEscape(name), shellEscape(template), shellEscape(name)),
		}, false)
	case Vite:
		template := opts.Template
		if template == "" {
			template = "react-ts"
		}
		// pnpm create vite lands the project and installs deps in a
		// second step. Templates: react-ts, vue-ts, svelte-ts, solid-ts,
		// qwik-ts, preact-ts, lit-ts, vanilla-ts, etc.
		return s.Exec.Exec(ctx, "frankenphp", []string{
			"bash", "-lc",
			fmt.Sprintf("cd /var/www/html && pnpm create vite@latest %s --template %s && cd %s && pnpm install",
				shellEscape(name), shellEscape(template), shellEscape(name)),
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
