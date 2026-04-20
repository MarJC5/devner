// Package runtime orchestrates the shared Docker Compose stack.
//
// We wrap `docker compose` CLI as a sub-process for V1. Reasons:
//   - Compose v2 is bundled with Docker Desktop, so every dev already has it.
//   - Re-implementing depends_on / healthcheck / network reconciliation with
//     the raw Docker SDK is hundreds of lines for no user-visible gain.
//   - We still use the Docker SDK directly for exec/logs/inspect (see docker.go).
package runtime

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/devner/devner/assets"
)

type Runtime struct {
	DataDir     string
	ProjectsDir string
	Project     string // compose project name, default "devner"
	Stdout      io.Writer
	Stderr      io.Writer
}

// GetStdoutErr / SetStdoutErr expose the Stdout/Stderr writers via an
// interface contract so callers that need to temporarily capture output
// (project.DevServer.execCapture, tui.silenceRuntime) can do so without
// knowing about the concrete Runtime type.
func (r *Runtime) GetStdoutErr() (io.Writer, io.Writer) { return r.Stdout, r.Stderr }
func (r *Runtime) SetStdoutErr(stdout, stderr io.Writer) {
	r.Stdout, r.Stderr = stdout, stderr
}

func New(dataDir, projectsDir string) *Runtime {
	return &Runtime{
		DataDir:     dataDir,
		ProjectsDir: projectsDir,
		Project:     "devner",
		Stdout:      os.Stdout,
		Stderr:     os.Stderr,
	}
}

// EnsureMaterialized extracts embedded compose.yaml, Dockerfile, php.ini and
// Caddyfile into DataDir so Compose can build from a stable path.
func (r *Runtime) EnsureMaterialized() error {
	if err := os.MkdirAll(r.DataDir, 0o755); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Join(r.DataDir, "dockerfiles", "frankenphp"), 0o755); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Join(r.DataDir, "frankenphp"), 0o755); err != nil {
		return err
	}

	files := map[string]string{
		"compose/compose.yaml":              filepath.Join(r.DataDir, "compose.yaml"),
		"dockerfiles/frankenphp.Dockerfile": filepath.Join(r.DataDir, "dockerfiles", "frankenphp", "Dockerfile"),
		"dockerfiles/php.ini":               filepath.Join(r.DataDir, "frankenphp", "php.ini"),
		"dockerfiles/Caddyfile":             filepath.Join(r.DataDir, "frankenphp", "Caddyfile"),
	}
	for src, dst := range files {
		data, err := assets.FS.ReadFile(src)
		if err != nil {
			return fmt.Errorf("read embed %s: %w", src, err)
		}
		if err := os.WriteFile(dst, data, 0o644); err != nil {
			return fmt.Errorf("write %s: %w", dst, err)
		}
	}
	if err := os.MkdirAll(r.ProjectsDir, 0o755); err != nil {
		return err
	}
	return nil
}

func (r *Runtime) composeCmd(ctx context.Context, args ...string) *exec.Cmd {
	composeFile := filepath.Join(r.DataDir, "compose.yaml")
	full := append([]string{"compose", "-f", composeFile, "-p", r.Project}, args...)
	cmd := exec.CommandContext(ctx, "docker", full...)
	cmd.Env = append(os.Environ(),
		"DEVNER_DATA_DIR="+r.DataDir,
		"DEVNER_PROJECTS_DIR="+r.ProjectsDir,
	)
	cmd.Stdout = r.Stdout
	cmd.Stderr = r.Stderr
	return cmd
}

func (r *Runtime) Up(ctx context.Context) error {
	if err := r.EnsureMaterialized(); err != nil {
		return err
	}
	return r.composeCmd(ctx, "up", "-d", "--build").Run()
}

func (r *Runtime) Down(ctx context.Context) error {
	return r.composeCmd(ctx, "down").Run()
}

func (r *Runtime) Stop(ctx context.Context) error {
	return r.composeCmd(ctx, "stop").Run()
}

func (r *Runtime) Restart(ctx context.Context, service string) error {
	if service == "" {
		return r.composeCmd(ctx, "restart").Run()
	}
	return r.composeCmd(ctx, "restart", service).Run()
}

func (r *Runtime) Rebuild(ctx context.Context) error {
	if err := r.EnsureMaterialized(); err != nil {
		return err
	}
	return r.composeCmd(ctx, "up", "-d", "--build", "--force-recreate").Run()
}

func (r *Runtime) Delete(ctx context.Context) error {
	return r.composeCmd(ctx, "down", "-v", "--remove-orphans").Run()
}

func (r *Runtime) PS(ctx context.Context) error {
	return r.composeCmd(ctx, "ps").Run()
}
