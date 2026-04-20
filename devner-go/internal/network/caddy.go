// Caddy admin API client.
//
// The shared stack runs Caddy (via FrankenPHP) with the admin API exposed on
// :2019. We push per-project site blocks as JSON config and delete them the
// same way — no Caddyfile concat, no reload-via-SIGHUP race condition.
package network

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type CaddyClient struct {
	BaseURL string
	HTTP    *http.Client
}

func NewCaddyClient() *CaddyClient {
	return &CaddyClient{
		BaseURL: "http://127.0.0.1:2019",
		HTTP:    &http.Client{Timeout: 10 * time.Second},
	}
}

// ProjectSite describes a site served by FrankenPHP for a devner project.
type ProjectSite struct {
	Domain string // e.g. myapp.localhost
	Root   string // container-side path, e.g. /var/www/html/myapp/public
}

// Apply replaces the full Caddy config with one derived from `sites`.
// Using PUT /config is atomic — no partial state if a single site is invalid.
func (c *CaddyClient) Apply(ctx context.Context, sites []ProjectSite) error {
	cfg := buildConfig(sites)
	body, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL+"/load", bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return fmt.Errorf("caddy admin: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		msg, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("caddy admin %d: %s", resp.StatusCode, msg)
	}
	return nil
}

func buildConfig(sites []ProjectSite) map[string]any {
	var servers map[string]any
	if len(sites) == 0 {
		servers = map[string]any{}
	} else {
		var routes []any
		for _, s := range sites {
			routes = append(routes, map[string]any{
				"match": []any{map[string]any{"host": []string{s.Domain}}},
				"handle": []any{
					map[string]any{
						"handler": "subroute",
						"routes": []any{
							map[string]any{
								"handle": []any{
									map[string]any{"handler": "vars", "root": s.Root},
									map[string]any{
										"handler":     "php",
										"root":        s.Root,
										"try_files":   []string{"{http.request.uri.path}", "{http.request.uri.path}/index.php", "index.php"},
										"split_path":  []string{".php"},
									},
									map[string]any{"handler": "file_server", "root": s.Root},
								},
							},
						},
					},
				},
				"terminal": true,
			})
		}
		servers = map[string]any{
			"srv0": map[string]any{
				"listen": []string{":443"},
				"routes": routes,
			},
			"srv_http": map[string]any{
				"listen": []string{":80"},
				"routes": routes,
			},
		}
	}

	return map[string]any{
		"admin": map[string]any{"listen": "0.0.0.0:2019"},
		"apps": map[string]any{
			"http": map[string]any{
				"servers": servers,
			},
		},
	}
}
