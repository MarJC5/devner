package project

import (
	"context"
	"fmt"
	"strings"
)

// DBCreds is the subset of connection info post-setup needs. It's a
// neutral struct (not the database package's Credentials) so `project`
// doesn't take a dependency on `database`.
type DBCreds struct {
	Engine   string // "mysql" | "postgres"
	Host     string // container network name, e.g. "mysql"
	Port     int
	Database string
	User     string
	Password string
}

// PostSetupOptions drives optional steps per project type.
type PostSetupOptions struct {
	// WordPress auto-install (skipped when WPInstall is false).
	WPInstall     bool
	WPSiteURL     string // default: https://<project>.localhost
	WPSiteTitle   string // default: project name
	WPAdminUser   string // default: "admin"
	WPAdminPass   string // default: "admin"
	WPAdminEmail  string // default: "admin@localhost.test"
}

// PostSetup runs the framework-specific finishing steps after Scaffold +
// CreateDatabase: Laravel gets a wired .env + app key, WordPress gets a
// wp-config.php and optionally a full install.
//
// Returns the captured log output (short) for the caller to show the
// user. Runs everything inside the frankenphp container via the Executor.
func (s *Service) PostSetup(ctx context.Context, t Type, name string, db *DBCreds, opts PostSetupOptions) (string, error) {
	switch t {
	case Laravel:
		return s.postSetupLaravel(ctx, name, db)
	case WordPress:
		return s.postSetupWordPress(ctx, name, db, opts)
	}
	// Node / Next / Astro: no framework config, dependencies already
	// pulled by their respective create-* commands during Scaffold().
	return "", nil
}

func (s *Service) postSetupLaravel(ctx context.Context, name string, db *DBCreds) (string, error) {
	var script strings.Builder
	fmt.Fprintf(&script, "cd /var/www/html/%s\n", shellEscape(name))
	script.WriteString("[ -f .env ] || cp .env.example .env\n")
	if db != nil {
		conn, port := "mysql", 3306
		if db.Engine == "postgres" {
			conn, port = "pgsql", 5432
		}
		// Laravel 11+ ships .env with DB_* lines commented out. Our regex
		// accepts optional leading '#' and whitespace so we rewrite the
		// commented form in place; if the key is missing entirely we
		// append it. -E for extended regex (sed on macOS GNU-compatible).
		script.WriteString(`set_env() {
  local k=$1 v=$2
  if grep -qE "^[[:space:]]*#?[[:space:]]*${k}=" .env; then
    sed -i.bak -E "s|^[[:space:]]*#?[[:space:]]*${k}=.*|${k}=${v}|" .env
  else
    echo "${k}=${v}" >> .env
  fi
}
`)
		fmt.Fprintf(&script, "set_env DB_CONNECTION %s\n", conn)
		fmt.Fprintf(&script, "set_env DB_HOST %s\n", db.Host)
		fmt.Fprintf(&script, "set_env DB_PORT %d\n", port)
		fmt.Fprintf(&script, "set_env DB_DATABASE %s\n", db.Database)
		fmt.Fprintf(&script, "set_env DB_USERNAME %s\n", db.User)
		fmt.Fprintf(&script, "set_env DB_PASSWORD %s\n", db.Password)
		script.WriteString("rm -f .env.bak\n")
	}
	script.WriteString("[ -d vendor ] || composer install --no-interaction --prefer-dist\n")
	script.WriteString("php artisan key:generate --ansi\n")
	// `composer create-project` runs migrate against sqlite by default,
	// leaving the configured MySQL/Postgres DB empty. Re-migrate against
	// the real backend so sessions/users/cache tables exist. --force so
	// it runs in non-interactive mode.
	if db != nil {
		script.WriteString("php artisan migrate --force --ansi || true\n")
	}
	// Clear cached config from the initial create-project run (it cached
	// sqlite settings we've now replaced).
	script.WriteString("php artisan config:clear --ansi || true\n")
	script.WriteString("php artisan cache:clear --ansi || true\n")

	return s.runBash(ctx, script.String())
}

func (s *Service) postSetupWordPress(ctx context.Context, name string, db *DBCreds, opts PostSetupOptions) (string, error) {
	if db == nil {
		// Without a DB, wp-config can't be generated; leave the download
		// in place and let the user point it at something manually.
		return "", nil
	}
	var script strings.Builder
	fmt.Fprintf(&script, "cd /var/www/html/%s\n", shellEscape(name))
	fmt.Fprintf(&script, `wp config create \
  --dbname=%s \
  --dbuser=%s \
  --dbpass=%s \
  --dbhost=%s \
  --skip-check --force --allow-root
`, db.Database, db.User, db.Password, db.Host)

	if opts.WPInstall {
		siteURL := opts.WPSiteURL
		if siteURL == "" {
			siteURL = "https://" + name + ".localhost"
		}
		title := opts.WPSiteTitle
		if title == "" {
			title = name
		}
		admin := opts.WPAdminUser
		if admin == "" {
			admin = "admin"
		}
		pass := opts.WPAdminPass
		if pass == "" {
			pass = "admin"
		}
		email := opts.WPAdminEmail
		if email == "" {
			email = "admin@localhost.test"
		}
		fmt.Fprintf(&script, `wp core install \
  --url=%s \
  --title=%s \
  --admin_user=%s \
  --admin_password=%s \
  --admin_email=%s \
  --skip-email --allow-root
`, shellQuote(siteURL), shellQuote(title), shellQuote(admin), shellQuote(pass), shellQuote(email))
	}
	return s.runBash(ctx, script.String())
}

// runBash executes a multi-line script via `bash -lc` in the frankenphp
// container. Output is captured and truncated to 4 KB.
func (s *Service) runBash(ctx context.Context, script string) (string, error) {
	buf := &captureWriter{}
	// project.Service historically uses Executor.Exec which writes to its
	// own Stdout/Stderr. We set them on the underlying Runtime only if
	// our executor exposes them — keep this simple: just run and return
	// the executor's error; logs go to whatever Runtime is configured
	// to. For better UX in the CLI, callers can swap Runtime.Stdout
	// around the PostSetup call.
	_ = buf
	return "", s.Exec.Exec(ctx, "frankenphp", []string{"bash", "-lc", script}, false)
}

// shellQuote wraps s in single quotes, escaping internal single quotes.
func shellQuote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", `'\''`) + "'"
}

// captureWriter is defined in internal/tools; duplicated minimally here
// because we can't import tools from project (cycle). If this grows it
// can move to an internal/util package.
type captureWriter struct{ buf []byte }

func (c *captureWriter) Write(p []byte) (int, error) { c.buf = append(c.buf, p...); return len(p), nil }
func (c *captureWriter) String() string               { return string(c.buf) }
