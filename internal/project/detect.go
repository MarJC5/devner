package project

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Detect inspects a project directory on disk and returns its best-guess
// Type. Returns ok=false if nothing recognizable is found.
func Detect(dir string) (Type, bool) {
	exists := func(rel string) bool {
		_, err := os.Stat(filepath.Join(dir, rel))
		return err == nil
	}
	switch {
	case exists("wp-config.php"), exists("wp-config-sample.php"), exists("wp-settings.php"):
		return WordPress, true
	case exists("artisan"):
		return Laravel, true
	case exists("astro.config.mjs"), exists("astro.config.ts"), exists("astro.config.js"):
		return Astro, true
	case exists("next.config.js"), exists("next.config.ts"), exists("next.config.mjs"):
		return NextJS, true
	case exists("package.json"):
		return Node, true
	}
	return "", false
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
