package project

import (
	"context"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// DevServer controls the lifecycle of a framework dev process (pnpm dev,
// npm run start, …) inside the frankenphp container. Each project has
// its own allocated port; the process is detached via `docker exec -d +
// nohup` so it survives this Go process. A PID file at
// /tmp/devner/dev-<name>.pid tracks the live PID for Stop and Status.
//
// Design choices:
//   - Stateless: a container restart (Caddy reload, rebuild, down) kills
//     all dev servers. We don't attempt to resurrect.
//   - Polling env vars are injected by default because Docker Desktop on
//     macOS doesn't propagate inotify events reliably through the bind
//     mount — without these, HMR misses changes.
type DevServer struct {
	Exec Executor
}

func NewDevServer(exec Executor) *DevServer {
	return &DevServer{Exec: exec}
}

// Status describes the current state of a dev server.
type Status struct {
	Name    string
	Running bool
	PID     int
	Port    int
}

// DefaultCommand returns the typical "run the dev server" command for a
// given project type, pre-wired to use the allocated `port`.
//
// Per-framework quirks:
//   - Next.js reads PORT from env, so we don't need to pass it here.
//   - Vite ignores PORT env (uses `server.port` in its config or
//     `--port` CLI). Force via `-- --port N`.
//   - Astro uses `astro dev --port N`; pass through pnpm with `--`.
//   - Node plain: whatever the user's npm script does. PORT env is
//     conventional so we set it and hope the user respects it.
func DefaultCommand(t Type, port int) string {
	switch t {
	case NextJS:
		return "pnpm dev" // PORT env wins
	case Astro, Vite:
		// pnpm forwards unrecognised flags to the underlying script;
		// `pnpm dev --port N` becomes `vite --port N` (or `astro dev
		// --port N`). Do NOT add `--` before the flags — pnpm keeps
		// `--` as a literal arg with Vite, breaking the parse.
		return fmt.Sprintf("pnpm dev --port %d", port)
	case Node:
		return "npm run start"
	}
	return "pnpm dev"
}

// Start spawns the dev server. `command` is the shell command to run
// inside the project dir (e.g. "pnpm dev"). Env vars PORT + HMR polling
// are prepended automatically so the user's script picks up the
// allocated port without manual config.
//
// If a dev server for this project is already running, Start is a
// no-op (returns nil). Call Stop first if you want to replace it.
func (d *DevServer) Start(ctx context.Context, name string, port int, command string) error {
	if command == "" {
		command = "pnpm dev"
	}
	safeName := shellEscape(name)

	// Quick "already running" probe.
	st, _ := d.Status(ctx, name)
	if st.Running {
		return nil
	}

	// Build: cd into project, launch with polling + PORT set, detach,
	// capture both streams to a log file we can tail later.
	//
	// We use `setsid` so the new bash process becomes the leader of its
	// own session + process group. `$!` then holds the PGID, which lets
	// Stop kill the whole tree (bash → pnpm → vite/next/astro → …)
	// with `kill -TERM -$pgid`. Without setsid, killing the bash
	// wrapper leaves pnpm + vite orphaned and still listening on the
	// port.
	script := fmt.Sprintf(`set -e
mkdir -p /tmp/devner
cd /var/www/html/%s
rm -f /tmp/devner/dev-%s.log
PORT=%d CHOKIDAR_USEPOLLING=1 CHOKIDAR_INTERVAL=300 WATCHPACK_POLLING=1 \
  setsid bash -lc %s > /tmp/devner/dev-%s.log 2>&1 < /dev/null &
echo $! > /tmp/devner/dev-%s.pid
`, safeName, safeName, port, shellQuoteSingle(command), safeName, safeName)

	// `docker exec -d` (detached) is crucial: without -d, the `nohup &`
	// still dies when the exec session terminates. With -d, docker gives
	// us a true fire-and-forget.
	return d.Exec.Exec(ctx, "frankenphp", []string{"bash", "-lc", script}, false)
}

// Stop reads the PID file and terminates the whole process group
// (bash → pnpm → vite/next/astro → …) via `kill -TERM -$pgid`.
// Idempotent — calling Stop on a not-running server is not an error.
func (d *DevServer) Stop(ctx context.Context, name string) error {
	safeName := shellEscape(name)
	script := fmt.Sprintf(`
if [ -f /tmp/devner/dev-%s.pid ]; then
  pid=$(cat /tmp/devner/dev-%s.pid)
  if [ -n "$pid" ]; then
    # Kill the whole process group (negative PID). 2>/dev/null swallows
    # "no such process" when children have already exited.
    kill -TERM -"$pid" 2>/dev/null || kill -TERM "$pid" 2>/dev/null || true
    for i in 1 2 3; do
      kill -0 "$pid" 2>/dev/null || break
      sleep 1
    done
    kill -KILL -"$pid" 2>/dev/null || kill -KILL "$pid" 2>/dev/null || true
  fi
  rm -f /tmp/devner/dev-%s.pid
fi
`, safeName, safeName, safeName)
	return d.Exec.Exec(ctx, "frankenphp", []string{"bash", "-lc", script}, false)
}

// Status returns the live state of a dev server. Running=true iff the
// PID file exists and the PID is still alive.
func (d *DevServer) Status(ctx context.Context, name string) (Status, error) {
	st := Status{Name: name}
	safeName := shellEscape(name)
	out := &captureWriter{}
	// Swap executor's underlying Runtime.Stdout to capture the small
	// output. project.Service used to do this; we mirror the trick.
	// Script emits "running <pid>" or "stopped" on one line.
	script := fmt.Sprintf(`
pid=""
if [ -f /tmp/devner/dev-%s.pid ]; then pid=$(cat /tmp/devner/dev-%s.pid); fi
if [ -n "$pid" ] && kill -0 "$pid" 2>/dev/null; then
  echo "running $pid"
else
  echo "stopped"
fi
`, safeName, safeName)

	if err := d.execCapture(ctx, script, out); err != nil {
		return st, err
	}
	line := strings.TrimSpace(out.String())
	if strings.HasPrefix(line, "running ") {
		st.Running = true
		pidStr := strings.TrimPrefix(line, "running ")
		if pid, perr := strconv.Atoi(pidStr); perr == nil {
			st.PID = pid
		}
	}
	return st, nil
}

// Logs returns the last `lines` lines of the dev server log file.
func (d *DevServer) Logs(ctx context.Context, name string, lines int) (string, error) {
	if lines <= 0 {
		lines = 100
	}
	safeName := shellEscape(name)
	out := &captureWriter{}
	script := fmt.Sprintf(`
if [ -f /tmp/devner/dev-%s.log ]; then
  tail -n %d /tmp/devner/dev-%s.log
else
  echo "(no log file yet — dev server never started, or was cleared on stop)"
fi
`, safeName, lines, safeName)
	if err := d.execCapture(ctx, script, out); err != nil {
		return out.String(), err
	}
	return out.String(), nil
}

// execCapture runs a bash script in frankenphp and captures stdout/stderr
// into `buf` by swapping the Executor's underlying Stdout. Ugly but
// necessary — Executor doesn't expose a "capture-to-buffer" mode.
//
// Assumes the Executor is *runtime.Runtime under the hood; if not, this
// falls back to running the command (which writes to the runtime's
// default Stdout) and returns the error. That path won't populate `buf`,
// which is fine — the Status/Logs functions degrade gracefully.
func (d *DevServer) execCapture(ctx context.Context, script string, buf *captureWriter) error {
	if r, ok := d.Exec.(stdoutSwapper); ok {
		oldOut, oldErr := r.GetStdoutErr()
		r.SetStdoutErr(buf, buf)
		defer r.SetStdoutErr(oldOut, oldErr)
	}
	return d.Exec.Exec(ctx, "frankenphp", []string{"bash", "-lc", script}, false)
}

// stdoutSwapper is implemented by runtime.Runtime (see internal/runtime).
// Declared here with the exact method shape we need so we don't have to
// import runtime (which would create a cycle).
type stdoutSwapper interface {
	GetStdoutErr() (io.Writer, io.Writer)
	SetStdoutErr(stdout, stderr io.Writer)
}

// shellQuoteSingle wraps s in single quotes and escapes any embedded
// single quote. Used for the nested command string.
func shellQuoteSingle(s string) string {
	return "'" + strings.ReplaceAll(s, "'", `'\''`) + "'"
}
