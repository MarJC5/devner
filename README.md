# Devner

Local dev environment orchestrator for WordPress, Laravel, Node, Next.js and Astro projects.

- Single cross-platform Go binary (macOS, Linux, Windows) — no bash, no Makefile for users
- **Three interfaces** sharing the same backend: CLI (Cobra) · TUI (Bubble Tea) · native GUI (Wails + React + Tailwind)
- macOS **menubar daemon** with Start / Stop / Open GUI and a live `stack N/M running` indicator
- Built-in AI agent: ask in natural language, it calls tools (create projects, databases, hosts, logs, exec)
- Multi-provider LLM: **Infomaniak AI Tools** (open-source models), Anthropic, local Ollama
- Docker stack: FrankenPHP (PHP+Caddy+Node), MySQL 8.4, PostgreSQL + PostGIS, Redis, Mailpit, Adminer

## Install

**macOS / Linux**
```bash
# From release binary (when published)
curl -L https://github.com/MarJC5/devner/releases/latest/download/devner_darwin_arm64.tar.gz | tar -xz
sudo mv devner /usr/local/bin/

# Or build from source
git clone -b go https://github.com/MarJC5/devner.git
cd devner
make install             # builds and installs to /usr/local/bin
```

**Windows** (PowerShell)
```powershell
# Download the .zip from Releases and extract devner.exe into PATH.
iwr -useb https://github.com/MarJC5/devner/releases/latest/download/devner_windows_amd64.zip -o devner.zip
Expand-Archive devner.zip
```

**Requirements**
- Docker Desktop (macOS, Windows) or Docker Engine + Docker Compose v2 (Linux)
- Go 1.22+ (to build from source)
- (optional) `mkcert` for system-wide HTTPS trust outside browsers

## Quickstart

```bash
# 1. Start the shared stack (pulls images the first time)
devner up

# 2. Create a Laravel project with MySQL
devner new laravel myapp --db=mysql
# → scaffolded, DB created, Caddy configured
# → https://myapp.localhost

# 3. Interactive TUI
devner tui

# 4. Ask the agent (Infomaniak by default)
#    First-run: edit ~/Library/Application Support/devner/config.toml,
#    set product_id and api_key under [llm.providers.infomaniak].
devner models                                             # sanity check auth
devner agent "create a Next.js project called web and show me logs"
```

## Migrate from an existing projects directory

```bash
devner import --source=~/path/to/old/projects --dry-run
# review what would be imported
devner import --source=~/path/to/old/projects
```

Detection:
- `wp-config.php` or `wp-settings.php` → WordPress (sniffs DB_NAME, DB_USER from wp-config.php)
- `artisan` → Laravel (sniffs DB_DATABASE from `.env`)
- `astro.config.*` → Astro
- `next.config.*` → Next.js
- `package.json` → Node

## Commands

```
devner up              # start stack
devner down            # stop stack
devner ps              # container status
devner logs <svc>      # tail container logs
devner restart [svc]   # restart stack / service
devner rebuild         # force rebuild images
devner delete --force  # destroy stack + volumes (DESTRUCTIVE)

devner new <type> <name> [--db=mysql|postgres] [--template=...]
# types: wordpress | laravel | node | nextjs | nuxt | astro | sveltekit | vite
# --template covers Vue / Svelte / Solid / Qwik / Preact / Lit / vanilla via Vite
devner remove <name>
devner list

devner dev start <name> [--command "..."]           # launch the project's dev server (Node-like only)
devner dev stop <name>
devner dev status [name]
devner dev logs <name> [--tail N]

devner db create <mysql|postgres> <name>
devner db drop   <mysql|postgres> <name>

devner hosts add <domain> [target]     # no-op for *.localhost
devner hosts remove <domain>

devner import --source=<path>          # migrate projects from a directory
devner reconcile [--apply]             # fix drift: store ↔ FS ↔ Caddy

devner certs status
devner certs install                   # mkcert -install

devner models                          # list models of the active LLM provider
devner agent "<prompt>"                # one-shot agent call
devner tui                             # interactive TUI
devner gui                             # launch the desktop app (Wails)
devner menubar                         # launch the macOS menubar daemon
devner menubar stop                    # kill it
devner menubar status                  # check whether it's running
```

## Configuration

The config file is created automatically on first run at the OS-standard location:

| OS | Path |
|---|---|
| macOS | `~/Library/Application Support/devner/config.toml` |
| Linux | `~/.config/devner/config.toml` |
| Windows | `%AppData%\devner\config.toml` |

Template (generated on first run, safe to edit):
```toml
[stack]
data_dir     = "~/.devner"
projects_dir = "~/devner/projects"

[llm]
active_provider = "infomaniak"

# Infomaniak AI Tools — OpenAI-compatible (v2 endpoint).
# Get your product_id from https://manager.infomaniak.com (AI Tools)
# or GET /1/ai with your token.
[llm.providers.infomaniak]
base_url    = "https://api.infomaniak.com/2/ai/{product_id}/openai/v1"
# Use either api_key (token in the file) OR api_key_env (env var name).
# api_key wins if both are set.
api_key     = ""
api_key_env = "INFOMANIAK_API_KEY"
product_id  = ""
model       = "qwen3"
kind        = "openai_compat"

[llm.providers.anthropic]
api_key_env = "ANTHROPIC_API_KEY"
model       = "claude-opus-4-7"
kind        = "anthropic"

# Local Ollama — no API key needed.
[llm.providers.ollama]
base_url    = "http://localhost:11434/v1"
model       = "qwen2.5-coder:14b"
kind        = "openai_compat"
```

### Infomaniak setup

1. Get your product ID: `curl -H "Authorization: Bearer $TOKEN" https://api.infomaniak.com/1/ai`
2. Put the product ID and the token in `config.toml` under `[llm.providers.infomaniak]`.
3. Verify connectivity and list available models:
   ```bash
   devner models
   ```
4. Run the agent:
   ```bash
   devner agent "list my projects"
   ```

Infomaniak v2 endpoint: `POST /2/ai/{product_id}/openai/v1/chat/completions` — see [docs](https://developer.infomaniak.com/docs/api/post/2/ai/%7Bproduct_id%7D/openai/v1/chat/completions).

## Node / Next / Nuxt / Astro / SvelteKit / Vite projects

**Framework coverage:**

| Type | Scaffold | Dev command | Notes |
|---|---|---|---|
| `nextjs` | `pnpm create next-app` | `pnpm dev` (PORT env) | App Router + TS + Tailwind by default |
| `nuxt` | `pnpm dlx nuxi init` | `pnpm dev --port N` | |
| `astro` | `pnpm create astro` | `pnpm dev --port N` | `--template minimal\|basics\|blog\|portfolio\|starlight` |
| `sveltekit` | `pnpm create svelte` | `pnpm dev --port N` | `--template skeleton\|minimal\|demo` |
| `vite` | `pnpm create vite` | `pnpm dev --port N` | `--template react-ts\|vue-ts\|svelte-ts\|solid-ts\|preact-ts\|qwik-ts\|lit-ts\|vanilla-ts` |
| `node` | `npm init -y` | `npm run start` | User controls the server code |

For Vue / Svelte / Solid / Qwik / Preact / Lit / vanilla-JS, use `devner new vite myapp --template <tpl>-ts` — no framework-specific devner type needed.


PHP projects (WordPress, Laravel) serve directly through FrankenPHP's PHP handler. Node-like projects use a different path:

1. At creation time, devner allocates an internal port in the range **3100–3999** and persists it on the project row (`projects.dev_port` in SQLite).
2. The rendered Caddyfile emits `reverse_proxy localhost:<port>` for that site instead of `php_server + file_server`.
3. Nothing listens on the port until you run `devner dev start`. Caddy returns **502** meanwhile — that's the "up but dev not running" signal.
4. `devner dev start <name>` spawns `pnpm dev` (or `npm run start` for plain Node) via `docker exec -d` inside the frankenphp container, under a new session so `devner dev stop` can terminate the whole process group (pnpm → vite/next/astro → …). A PID file lives at `/tmp/devner/dev-<name>.pid`; logs at `/tmp/devner/dev-<name>.log`.
5. The URL is the same as PHP projects: `https://<name>.localhost`. HTTPS, HMR via WebSocket, everything transparent.

### Port handling per framework

Vite and Astro don't read `PORT` from env by default — devner invokes them as `pnpm dev --port <N>`. Next.js reads `PORT` from env so devner sets it directly. Plain Node projects run `npm run start` and rely on the user's `package.json` honouring the `PORT` env var (conventional). Pass `--command "..."` to `devner dev start` to override.

### Hot reload on macOS + Docker Desktop

Bind-mount file events don't propagate reliably from the host into the container on macOS. Devner sets `CHOKIDAR_USEPOLLING=1`, `CHOKIDAR_INTERVAL=300`, and `WATCHPACK_POLLING=1` when spawning dev servers so chokidar/webpack-style watchers fall back to polling. Costs ~1–3 % of one core for instant HMR reliability. Enable Docker Desktop's VirtioFS file-sharing (Settings > General) if you prefer native fsevents without the polling overhead — the env vars become redundant but harmless.

### Stateless lifecycle

Dev servers die when the frankenphp container restarts (which happens whenever devner applies a new Caddy config — i.e. every `devner new` / `devner remove`). This is **intentional**: devner does not auto-restart them. Run `devner dev start <name>` again when you want it back.

### Quickstart

```bash
devner new vite myapp          # allocates port, scaffolds pnpm create vite react-ts
devner dev start myapp         # launches pnpm dev --port <allocated>
open https://myapp.localhost   # HTTP 200, Vite dev server with HMR
devner dev logs myapp          # tail the dev server output
devner dev stop myapp          # kill the process group, Caddy → 502
```

## Agent tools

The LLM sees these tools. Destructive ones (marked ⚠) require confirmation in the TUI or `--yes` on the CLI.

| Tool | Purpose |
|---|---|
| `list_projects` | list managed projects |
| `project_status` | details for one project |
| `create_project` | scaffold wordpress / laravel / node / nextjs / astro / vite |
| `delete_project` ⚠ | remove files + DB + Caddy entry |
| `create_database` | mysql or postgres DB + user |
| `start_dev_server` ⚠ | launch pnpm dev / npm run start for a Node-like project |
| `stop_dev_server` ⚠ | stop the dev server |
| `dev_server_status` | running/stopped + PID + port |
| `drop_database` ⚠ | drop DB + user |
| `start_stack` / `stop_stack` | lifecycle |
| `rebuild_stack` ⚠ | rebuild images + recreate |
| `tail_logs` | container logs |
| `composer` ⚠ | run composer install / require / update in the project |
| `npm` ⚠ | run npm / pnpm / yarn install / build / add in the project |
| `wp_cli` ⚠ | run WP-CLI: plugin install, option get, user create, etc. |
| `artisan` ⚠ | run `php artisan` on a Laravel project |
| `exec_in_project` ⚠ | arbitrary shell escape hatch (prefer the typed tools above) |

Every tool call is audited in `~/.devner/store.db` (`history` table).

## Architecture

```
cmd/devner/          # main (Cobra + Bubble Tea entrypoint)
internal/
  cli/               # Cobra subcommands
  tui/               # Bubble Tea scenes (projects, stack, chat, logs, config)
  agent/             # tool-use loop (max 10 iterations, destructive confirm)
  llm/               # Provider interface + OpenAI-compat + Anthropic native
  tools/             # Registry + JSON schemas
  project/           # WP/Laravel/Node/Next/Astro detect + scaffold
  database/          # MySQL + Postgres ops (validated identifiers, random passwords)
  runtime/           # docker compose wrapper + docker exec/logs
  network/           # Caddyfile rendering + docker restart reload
  store/             # SQLite metadata (modernc — pure Go, no CGO)
  app/               # Deps bundle (shared CLI + TUI wiring)
  config/            # Viper TOML
  platform/          # OS detection, hosts path, elevation check
assets/              # embed.FS: compose.yaml, Dockerfile, php.ini, Caddyfile
migrations/          # SQLite schema (.sql + embed.FS)
```

## Desktop app (GUI)

Cross-platform desktop app built with [Wails v2](https://wails.io) (Go backend + React + Tailwind frontend + native OS WebView). Shares `internal/app.Deps` with the CLI / TUI so every action goes through the same code path — no HTTP layer, no divergence.

- Dark + light themes (Sun / Moon / System — follows macOS preference)
- EN / FR UI (persisted in `localStorage`, selectable in Settings)
- Frosted-glass panels on macOS (vibrancy) and Windows (Mica); flat clean look on Linux
- Keyboard shortcuts: `⌘1..7` jump between Projects / Stack / Chat / Logs / Databases / Hosts / Settings

```bash
make gui-dev        # hot-reload Vite + auto Go rebuild on save
make gui-build      # produces gui/build/bin/devner-gui(.app|.exe|binary)
make gui-clean
```

Requirements: Go 1.22, Node 20, pnpm 9, and the Wails CLI:
```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

Launched from the CLI with `devner gui` — the CLI searches `/Applications/Devner.app`, `gui/build/bin/Devner.app`, `devner-gui` on `$PATH`, or the `$DEVNER_GUI` env var.

## Menubar daemon (macOS / Windows)

A tiny standalone process (`cmd/devner-menubar`) owns a system-tray icon with Start stack / Stop stack / Open GUI entries and a live `stack N/M running` indicator. It runs as a separate binary because Wails and `fyne.io/systray` both want ownership of the main run loop.

```bash
devner menubar          # launch (backgrounded)
devner menubar stop
devner menubar status
```

## Development

```bash
make build              # ./bin/devner
make build-menubar      # ./bin/devner-menubar
make build-all-bins     # CLI + menubar
make test               # go test ./...
make install            # copies devner to /usr/local/bin (may need sudo)
make clean              # removes ./bin and ./dist

# Without Make:
go build -o bin/devner ./cmd/devner
go test ./...
```

Cross-compile the CLI:
```bash
make build-all      # bin/devner-{darwin,linux,windows}-{amd64,arm64}
```

Release (local dry-run):
```bash
make release-snapshot    # goreleaser --snapshot --clean
```

Release CI (GitHub Actions) attaches GUI archives — `Devner_darwin_{arm64,amd64}.zip`, `Devner_linux_amd64.tar.gz`, `Devner_windows_amd64.zip` — alongside the GoReleaser-produced CLI archives on every `v*` tag push.
