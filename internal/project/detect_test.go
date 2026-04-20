package project

import (
	"os"
	"path/filepath"
	"testing"
)

// fixture builds a temp dir from a map of relative-path → contents. Empty
// contents means "create an empty file"; a path ending in `/` creates a
// directory. Returns the absolute dir path.
func fixture(t *testing.T, files map[string]string) string {
	t.Helper()
	dir := t.TempDir()
	for rel, body := range files {
		full := filepath.Join(dir, rel)
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatalf("mkdir %s: %v", full, err)
		}
		if err := os.WriteFile(full, []byte(body), 0o644); err != nil {
			t.Fatalf("write %s: %v", full, err)
		}
	}
	return dir
}

func TestDetect(t *testing.T) {
	cases := []struct {
		name     string
		files    map[string]string
		wantType Type
		wantMode DevMode
		wantCmd  string
	}{
		{
			name:     "wordpress",
			files:    map[string]string{"wp-config.php": ""},
			wantType: WordPress,
			wantMode: DevModeNone,
		},
		{
			name:     "wordpress-with-vite-assets",
			files:    map[string]string{"wp-config.php": "", "package.json": `{"scripts":{"watch":"vite build --watch"}}`},
			wantType: WordPress,
			wantMode: DevModeWatch,
			wantCmd:  "watch",
		},
		{
			name:     "laravel",
			files:    map[string]string{"artisan": ""},
			wantType: Laravel,
			wantMode: DevModeNone,
		},
		{
			name:     "laravel-with-mix-watch",
			files:    map[string]string{"artisan": "", "package.json": `{"scripts":{"watch":"mix watch","dev":"mix"}}`},
			wantType: Laravel,
			wantMode: DevModeWatch,
			wantCmd:  "watch",
		},
		{
			name: "php-generic-with-htaccess",
			files: map[string]string{
				".htaccess":         "RewriteEngine On\nRewriteRule ^(.*)$ /public/$1 [L]",
				"public/index.php":  "<?php",
				"package.json":      `{"scripts":{"dev":"vite","watch":"vite build --watch"}}`,
			},
			wantType: PHP,
			wantMode: DevModeWatch,
			wantCmd:  "watch",
		},
		{
			name: "php-generic-index-php-only",
			files: map[string]string{
				"index.php": "<?php echo 'hello';",
			},
			wantType: PHP,
			wantMode: DevModeNone,
		},
		{
			name:     "nextjs",
			files:    map[string]string{"next.config.js": "", "package.json": `{"scripts":{"dev":"next dev"}}`},
			wantType: NextJS,
			wantMode: DevModeServer,
			wantCmd:  "dev",
		},
		{
			name:     "nuxt",
			files:    map[string]string{"nuxt.config.ts": "", "package.json": `{"scripts":{"dev":"nuxi dev"}}`},
			wantType: Nuxt,
			wantMode: DevModeServer,
			wantCmd:  "dev",
		},
		{
			name:     "astro",
			files:    map[string]string{"astro.config.mjs": "", "package.json": `{"scripts":{"dev":"astro dev"}}`},
			wantType: Astro,
			wantMode: DevModeServer,
			wantCmd:  "dev",
		},
		{
			name:     "sveltekit",
			files:    map[string]string{"svelte.config.js": "", "package.json": `{"scripts":{"dev":"vite dev"}}`},
			wantType: SvelteKit,
			wantMode: DevModeServer,
			wantCmd:  "dev",
		},
		{
			name:     "vite-app",
			files:    map[string]string{"vite.config.js": "", "package.json": `{"scripts":{"dev":"vite"}}`},
			wantType: Vite,
			wantMode: DevModeServer,
			wantCmd:  "dev",
		},
		{
			name: "node-library",
			files: map[string]string{
				"package.json": `{
					"name":"my-lib","main":"dist/lib.js","module":"dist/lib.mjs",
					"exports":{".":"./dist/lib.mjs"},
					"files":["dist"],
					"scripts":{"build":"rollup -c","test":"jest"}
				}`,
			},
			wantType: Node,
			wantMode: DevModeNone,
		},
		{
			name: "node-app-with-start",
			files: map[string]string{
				"package.json": `{"scripts":{"start":"nodemon server.js"}}`,
			},
			wantType: Node,
			wantMode: DevModeServer,
			wantCmd:  "start",
		},
		{
			name:     "unrecognized",
			files:    map[string]string{"README.md": "# hello"},
			wantType: "",
			wantMode: DevModeNone,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			dir := fixture(t, tc.files)
			got, ok := Detect(dir)
			if tc.wantType == "" {
				if ok {
					t.Fatalf("expected unrecognized, got %+v", got)
				}
				return
			}
			if !ok {
				t.Fatalf("expected type=%s, got ok=false", tc.wantType)
			}
			if got.Type != tc.wantType {
				t.Errorf("Type = %q, want %q", got.Type, tc.wantType)
			}
			if got.DevMode != tc.wantMode {
				t.Errorf("DevMode = %q, want %q", got.DevMode, tc.wantMode)
			}
			if got.DevCommand != tc.wantCmd {
				t.Errorf("DevCommand = %q, want %q", got.DevCommand, tc.wantCmd)
			}
		})
	}
}
