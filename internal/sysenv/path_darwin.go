//go:build darwin

// Package sysenv contains process-environment helpers that are too
// small to deserve their own home but are shared between the GUI, the
// menubar daemon, and (potentially) the CLI.
package sysenv

import (
	"os"
	"strings"
)

// EnsureDevPath prepends the standard Homebrew and Docker Desktop bin
// directories to $PATH if they are not already present.
//
// macOS apps launched via LaunchServices (Finder, Dock, Spotlight,
// Login Items, `open Foo.app`) inherit a minimal
// `/usr/bin:/bin:/usr/sbin:/sbin` PATH that excludes both
// `/usr/local/bin` and `/opt/homebrew/bin`. Docker Desktop installs
// `docker` as a symlink under `/usr/local/bin`, so any
// `exec.Command("docker", …)` from a .app bundle fails with
// "executable file not found".
//
// Call this once, as early as possible in main(), before any code
// shells out. Idempotent: when launched from a terminal the inherited
// PATH already contains these directories and the call is a no-op.
func EnsureDevPath() {
	want := []string{"/opt/homebrew/bin", "/usr/local/bin"}
	cur := os.Getenv("PATH")
	parts := strings.Split(cur, string(os.PathListSeparator))

	has := make(map[string]bool, len(parts))
	for _, p := range parts {
		has[p] = true
	}

	var prepend []string
	for _, p := range want {
		if !has[p] {
			prepend = append(prepend, p)
		}
	}
	if len(prepend) == 0 {
		return
	}
	_ = os.Setenv("PATH", strings.Join(append(prepend, parts...), string(os.PathListSeparator)))
}
