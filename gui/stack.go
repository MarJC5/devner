package main

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	wailsruntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

// ContainerStatus is the flat, JS-friendly shape for one container in
// the shared stack. We build this ourselves rather than expose the raw
// Docker types — simpler TS surface.
type ContainerStatus struct {
	Name   string `json:"name"`   // e.g. "frankenphp_devner"
	Image  string `json:"image"`
	State  string `json:"state"`  // running | restarting | exited | created
	Status string `json:"status"` // "Up 5 minutes (healthy)"
}

// StackStatus queries docker once and returns the current containers.
// Used for the initial paint on the Stack screen. After that, the live
// updates arrive via the "stack:tick" event (see startStackTicker).
func (a *App) StackStatus() ([]ContainerStatus, error) {
	if a.deps == nil {
		return nil, fmt.Errorf("app deps not initialized")
	}
	return stackStatus(a.ctx)
}

func (a *App) StartStack() error {
	if a.deps == nil {
		return fmt.Errorf("app deps not initialized")
	}
	return a.deps.Runtime.Up(a.ctx)
}

func (a *App) StopStack() error {
	if a.deps == nil {
		return fmt.Errorf("app deps not initialized")
	}
	return a.deps.Runtime.Stop(a.ctx)
}

func (a *App) RebuildStack() error {
	if a.deps == nil {
		return fmt.Errorf("app deps not initialized")
	}
	return a.deps.Runtime.Rebuild(a.ctx)
}

// stackStatus shells out to `docker ps` and parses the TSV output. We
// don't use the Docker Go SDK here because it would force CGO-ish
// dependencies on the GUI module — the CLI already uses this pattern
// in internal/tui/stack.go for the same reason.
func stackStatus(ctx context.Context) ([]ContainerStatus, error) {
	timeout, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	cmd := exec.CommandContext(timeout, "docker", "ps",
		"--filter", "name=_devner",
		"--format", "{{.Names}}\t{{.Image}}\t{{.State}}\t{{.Status}}")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	var result []ContainerStatus
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, "\t", 4)
		if len(parts) != 4 {
			continue
		}
		result = append(result, ContainerStatus{
			Name: parts[0], Image: parts[1], State: parts[2], Status: parts[3],
		})
	}
	return result, nil
}

// startStackTicker emits "stack:tick" every 3 s with the current
// container list. Frontend listens via EventsOn("stack:tick", …).
// Stopped when ctx is cancelled (Wails calls shutdown → cancels).
func (a *App) startStackTicker(ctx context.Context) {
	go func() {
		tick := time.NewTicker(3 * time.Second)
		defer tick.Stop()
		// First tick immediately so the UI doesn't sit on stale state.
		a.emitStack()
		for {
			select {
			case <-ctx.Done():
				return
			case <-tick.C:
				a.emitStack()
			}
		}
	}()
}

func (a *App) emitStack() {
	if a.deps == nil {
		return
	}
	st, err := stackStatus(a.ctx)
	if err != nil {
		// Emit empty list on error — UI decides whether to show an
		// error banner; we don't spam events for transient failures.
		wailsruntime.EventsEmit(a.ctx, "stack:tick", []ContainerStatus{})
		return
	}
	wailsruntime.EventsEmit(a.ctx, "stack:tick", st)
}
