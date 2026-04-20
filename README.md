# Devner

Local dev environment orchestrator for WordPress, Laravel, Node, Next.js and Astro projects.

- Single cross-platform Go binary (macOS, Linux, Windows) — no bash, no Makefile for users
- TUI (Bubble Tea) + scripting CLI (Cobra)
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

devner new <type> <name> [--db=mysql|postgres]
devner remove <name>
devner list

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

## Agent tools

The LLM sees these tools. Destructive ones (marked ⚠) require confirmation in the TUI or `--yes` on the CLI.

| Tool | Purpose |
|---|---|
| `list_projects` | list managed projects |
| `project_status` | details for one project |
| `create_project` | scaffold wordpress / laravel / node / nextjs / astro |
| `delete_project` ⚠ | remove files + DB + Caddy entry |
| `create_database` | mysql or postgres DB + user |
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

## Development

```bash
make build      # builds ./bin/devner
make test       # go test ./...
make install    # copies binary to /usr/local/bin/devner (may need sudo)
make clean      # removes ./bin and ./dist

# Without Make:
go build -o bin/devner ./cmd/devner
go test ./...
```

Cross-compile:
```bash
make build-all      # bin/devner-{darwin,linux,windows}-{amd64,arm64}
```

Release (local dry-run):
```bash
make release-snapshot    # goreleaser --snapshot --clean
```
