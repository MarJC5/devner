package runtime

import (
	"context"
	"fmt"
	"os"
	"os/exec"
)

type ContainerStatus struct {
	Name   string
	Image  string
	State  string
	Health string
}

// Exec runs a command inside a service container via `docker exec`.
func (r *Runtime) Exec(ctx context.Context, service string, cmdArgs []string, interactive bool) error {
	container := service + "_devner"
	args := []string{"exec"}
	if interactive {
		args = append(args, "-it")
	}
	args = append(args, container)
	args = append(args, cmdArgs...)

	cmd := exec.CommandContext(ctx, "docker", args...)
	if interactive {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	} else {
		cmd.Stdout = r.Stdout
		cmd.Stderr = r.Stderr
	}
	return cmd.Run()
}

func (r *Runtime) Logs(ctx context.Context, service string, tail int, follow bool) error {
	container := service + "_devner"
	args := []string{"logs"}
	if tail > 0 {
		args = append(args, "--tail", fmt.Sprintf("%d", tail))
	}
	if follow {
		args = append(args, "-f")
	}
	args = append(args, container)

	cmd := exec.CommandContext(ctx, "docker", args...)
	cmd.Stdout = r.Stdout
	cmd.Stderr = r.Stderr
	return cmd.Run()
}
