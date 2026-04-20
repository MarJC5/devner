// Package network — Caddy configuration for devner sites.
//
// FrankenPHP bundles Caddy with a non-standard admin lifecycle: PUT/POST to
// the admin API causes Caddy to restart its own admin listener, which kills
// the in-flight HTTP connection (symptom: "stopping admin server: 10s
// timeout" every reload). Instead of fighting it, we render a Caddyfile to
// a bind-mounted path and ask Caddy to reload it via `caddy reload` — this
// is hot-reload (no downtime, no port flip, no TLS cert churn).
package network

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

type CaddyClient struct {
	// CaddyfilePath is the host-side path to the Caddyfile. Mounted into
	// the container at /etc/frankenphp/Caddyfile.
	CaddyfilePath string
	// Container is the docker compose container name, e.g. "frankenphp_devner".
	Container string
}

func NewCaddyClient() *CaddyClient {
	home, _ := os.UserHomeDir()
	return &CaddyClient{
		CaddyfilePath: filepath.Join(home, ".devner", "frankenphp", "Caddyfile"),
		Container:     "frankenphp_devner",
	}
}

// ProjectSite describes a site served by FrankenPHP for a devner project.
type ProjectSite struct {
	Domain string // e.g. myapp.localhost
	Root   string // container-side path, e.g. /var/www/html/myapp/public
}

// Apply writes a Caddyfile with one site block per project and restarts
// FrankenPHP to load it. Idempotent.
//
// Note on reload strategy: `caddy reload` / `frankenphp reload` both use the
// admin API, which FrankenPHP currently mishandles on reload (restarting
// its own admin listener times out the request — "stopping admin server:
// 10s timeout"). A `docker restart` of the container is ~3s of downtime,
// acceptable for a local dev tool, and completely reliable. TLS certs are
// persisted in the caddy_data volume so no re-issuance.
func (c *CaddyClient) Apply(ctx context.Context, sites []ProjectSite) error {
	cf := renderCaddyfile(sites)
	if err := os.MkdirAll(filepath.Dir(c.CaddyfilePath), 0o755); err != nil {
		return err
	}
	if err := os.WriteFile(c.CaddyfilePath, []byte(cf), 0o644); err != nil {
		return fmt.Errorf("write Caddyfile: %w", err)
	}

	cmd := exec.CommandContext(ctx, "docker", "restart", "-t", "5", c.Container)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("docker restart: %w: %s", err, strings.TrimSpace(string(out)))
	}
	return nil
}

// renderCaddyfile produces a Caddyfile covering the global admin block plus
// one site block per project. Sites are sorted for deterministic output.
func renderCaddyfile(sites []ProjectSite) string {
	sorted := make([]ProjectSite, len(sites))
	copy(sorted, sites)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].Domain < sorted[j].Domain })

	var b strings.Builder
	b.WriteString("{\n")
	b.WriteString("\t# Keep admin API enabled for tooling (curl, diagnostics).\n")
	b.WriteString("\tadmin 0.0.0.0:2019\n")
	b.WriteString("\tfrankenphp\n")
	b.WriteString("}\n\n")

	for _, s := range sorted {
		fmt.Fprintf(&b, "%s {\n", s.Domain)
		fmt.Fprintf(&b, "\troot * %s\n", s.Root)
		b.WriteString("\tphp_server\n")
		b.WriteString("\tfile_server\n")
		b.WriteString("}\n\n")
	}
	return b.String()
}
