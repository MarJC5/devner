package project

import (
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// DetectResult is the full classification of a project directory: its
// primary Type (WordPress, Laravel, PHP, Node, Next, Vite, …) plus the
// dev-time mode that `devner dev start` should respect (HMR server,
// asset watcher, or nothing). DevCommand is the pnpm script name to run
// (e.g. "dev", "watch") when DevMode != DevModeNone.
type DetectResult struct {
	Type       Type
	DevMode    DevMode
	DevCommand string
}

// Detect inspects a project directory on disk and classifies it. Returns
// ok=false only when nothing at all is recognizable (no PHP markers, no
// framework config, no package.json). Empty-but-usable classifications
// like a JS library (Type=Node, DevMode=None) still return ok=true.
func Detect(dir string) (DetectResult, bool) {
	exists := func(rel string) bool {
		_, err := os.Stat(filepath.Join(dir, rel))
		return err == nil
	}

	// Step 1: pick a Type.
	//
	// PHP markers win over JS markers because many PHP projects carry a
	// package.json for asset compilation (Vite, Tailwind, Laravel Mix).
	// If we let package.json match first, every PHP+assets project would
	// be misclassified as Node.
	var typ Type
	switch {
	case exists("wp-config.php"), exists("wp-config-sample.php"), exists("wp-settings.php"):
		typ = WordPress
	case exists("artisan"):
		typ = Laravel
	case isGenericPHP(dir, exists):
		typ = PHP
	case exists("astro.config.mjs"), exists("astro.config.ts"), exists("astro.config.js"):
		typ = Astro
	case exists("next.config.js"), exists("next.config.ts"), exists("next.config.mjs"):
		typ = NextJS
	case exists("nuxt.config.js"), exists("nuxt.config.ts"), exists("nuxt.config.mjs"):
		typ = Nuxt
	case exists("svelte.config.js"), exists("svelte.config.ts"):
		typ = SvelteKit
	case exists("vite.config.js"), exists("vite.config.ts"), exists("vite.config.mjs"):
		typ = Vite
	case exists("package.json"):
		typ = Node
	default:
		return DetectResult{}, false
	}

	// Step 2: dev mode + command from package.json, if one exists.
	mode, cmd := detectDevMode(dir, typ)
	return DetectResult{Type: typ, DevMode: mode, DevCommand: cmd}, true
}

// isGenericPHP checks for markers of a custom/framework-less PHP project
// (index.php front-controller, .htaccess with PHP rewrite rules, or a
// composer.json). WordPress and Laravel are matched earlier; this is the
// "some other PHP site" catch-all.
func isGenericPHP(dir string, exists func(string) bool) bool {
	if exists("index.php") || exists("public/index.php") || exists("composer.json") {
		return true
	}
	if data, err := os.ReadFile(filepath.Join(dir, ".htaccess")); err == nil {
		// A RewriteRule pointing at index.php (or /public/) is the
		// classic PHP front-controller signal. Without it, .htaccess
		// could belong to anything (static site, Apache config).
		lower := strings.ToLower(string(data))
		if strings.Contains(lower, "index.php") ||
			strings.Contains(lower, "/public/") ||
			strings.Contains(lower, "rewriterule") && strings.Contains(lower, ".php") {
			return true
		}
	}
	return false
}

type pkgJSON struct {
	Main    string            `json:"main"`
	Module  string            `json:"module"`
	Exports json.RawMessage   `json:"exports"`
	Files   []string          `json:"files"`
	Scripts map[string]string `json:"scripts"`
	Private *bool             `json:"private"`
}

// detectDevMode parses package.json and decides:
//   - DevModeServer: Node-app types with a framework dev script.
//   - DevModeWatch:  PHP-typed projects with any asset script (watch or
//     dev, since for a PHP backend even `vite` serves assets only).
//   - DevModeNone:   libraries (main/exports/files + no dev script), or
//     projects without package.json at all.
func detectDevMode(dir string, typ Type) (DevMode, string) {
	data, err := os.ReadFile(filepath.Join(dir, "package.json"))
	if err != nil {
		return DevModeNone, ""
	}
	var p pkgJSON
	if err := json.Unmarshal(data, &p); err != nil {
		return DevModeNone, ""
	}
	scripts := p.Scripts

	// PHP-family: any asset script means "watch mode". Prefer "watch"
	// (compile once then rebuild on change) over "dev" (launches a Vite
	// server on :5173 that we'd ignore anyway because PHP serves HTTP).
	if typ == WordPress || typ == Laravel || typ == PHP {
		if _, ok := scripts["watch"]; ok {
			return DevModeWatch, "watch"
		}
		if _, ok := scripts["dev"]; ok {
			return DevModeWatch, "dev"
		}
		return DevModeNone, ""
	}

	// Node-app types: want an HMR server. Detect it via the dev script's
	// content (to distinguish "vite" from a library's "build:dev" noise).
	if dev, ok := scripts["dev"]; ok && looksLikeDevServer(dev) {
		return DevModeServer, "dev"
	}
	if start, ok := scripts["start"]; ok && looksLikeDevServer(start) {
		return DevModeServer, "start"
	}

	// Library heuristic: no dev script, but clearly published-to-npm
	// shape. `"files"` is a strong signal (it's only meaningful for
	// published packages); `main`/`module`/`exports` reinforce it.
	if isLibraryShape(p) {
		return DevModeNone, ""
	}

	// Node type with no recognized dev script and no library shape —
	// could be a plain script project. Fall back to Node default
	// (npm run start) as a server if present, otherwise none.
	if _, ok := scripts["start"]; ok {
		return DevModeServer, "start"
	}
	if _, ok := scripts["dev"]; ok {
		return DevModeServer, "dev"
	}
	return DevModeNone, ""
}

// devServerRe matches the common "this is an HMR dev server" signatures
// we care about. Kept as a single regex for speed (runs on every detect).
var devServerRe = regexp.MustCompile(`\b(vite|next\s+dev|astro\s+dev|nuxi\s+dev|nuxt\s+dev|svelte-kit\s+dev|nodemon|tsx\s+watch|ts-node-dev|concurrently)\b`)

func looksLikeDevServer(script string) bool {
	return devServerRe.MatchString(strings.ToLower(script))
}

func isLibraryShape(p pkgJSON) bool {
	if len(p.Files) > 0 {
		return true
	}
	if p.Main != "" || p.Module != "" || len(p.Exports) > 0 {
		// These alone aren't enough (some apps set `main`); combined
		// with no dev script (the caller already checked) it's a strong
		// signal.
		return true
	}
	return false
}

var (
	wpDBNameRe = regexp.MustCompile(`(?m)define\(\s*['"]DB_NAME['"]\s*,\s*['"]([^'"]+)['"]`)
	wpDBUserRe = regexp.MustCompile(`(?m)define\(\s*['"]DB_USER['"]\s*,\s*['"]([^'"]+)['"]`)
	wpDBHostRe = regexp.MustCompile(`(?m)define\(\s*['"]DB_HOST['"]\s*,\s*['"]([^'"]+)['"]`)
)

type DBInfo struct {
	Engine string // mysql | postgres | ""
	Name   string
	User   string
}

// SniffDB attempts to extract database engine + name from a project's
// config files. Returns zero DBInfo if nothing found — this is fine, the
// import still records the project, just without DB metadata.
func SniffDB(dir string, t Type) DBInfo {
	switch t {
	case WordPress:
		return sniffWordPress(dir)
	case Laravel:
		return sniffLaravel(dir)
	}
	return DBInfo{}
}

func sniffWordPress(dir string) DBInfo {
	data, err := os.ReadFile(filepath.Join(dir, "wp-config.php"))
	if err != nil {
		return DBInfo{}
	}
	info := DBInfo{Engine: "mysql"}
	if m := wpDBNameRe.FindStringSubmatch(string(data)); len(m) == 2 {
		info.Name = m[1]
	}
	if m := wpDBUserRe.FindStringSubmatch(string(data)); len(m) == 2 {
		info.User = m[1]
	}
	if m := wpDBHostRe.FindStringSubmatch(string(data)); len(m) == 2 {
		host := strings.ToLower(m[1])
		if strings.HasPrefix(host, "postgres") {
			info.Engine = "postgres"
		}
	}
	return info
}

func sniffLaravel(dir string) DBInfo {
	data, err := os.ReadFile(filepath.Join(dir, ".env"))
	if err != nil {
		return DBInfo{Engine: "mysql"}
	}
	info := DBInfo{Engine: "mysql"}
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "#") || line == "" {
			continue
		}
		k, v, ok := cut(line, "=")
		if !ok {
			continue
		}
		v = strings.Trim(v, `"'`)
		switch k {
		case "DB_CONNECTION":
			if strings.HasPrefix(strings.ToLower(v), "pgsql") || strings.HasPrefix(strings.ToLower(v), "postgres") {
				info.Engine = "postgres"
			}
		case "DB_DATABASE":
			info.Name = v
		case "DB_USERNAME":
			info.User = v
		}
	}
	return info
}

func cut(s, sep string) (before, after string, found bool) {
	if i := strings.Index(s, sep); i >= 0 {
		return s[:i], s[i+len(sep):], true
	}
	return s, "", false
}
