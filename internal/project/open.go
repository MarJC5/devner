package project

import (
	"fmt"
	"os"
	"os/exec"
)

// SupportedEditors lists the editors we know how to launch. The string is
// the binary name expected on PATH. Order matters for DefaultEditor's
// auto-detection: earliest = highest priority.
var SupportedEditors = []string{"code", "cursor", "zed", "subl", "webstorm", "phpstorm", "idea"}

// DefaultEditor picks the first editor from SupportedEditors that is on
// PATH. Falls back to $VISUAL, then $EDITOR, then returns an error.
func DefaultEditor() (string, error) {
	for _, e := range SupportedEditors {
		if _, err := exec.LookPath(e); err == nil {
			return e, nil
		}
	}
	for _, env := range []string{"VISUAL", "EDITOR"} {
		if v := os.Getenv(env); v != "" {
			if _, err := exec.LookPath(v); err == nil {
				return v, nil
			}
		}
	}
	return "", fmt.Errorf("no editor found on PATH (tried %v and $VISUAL / $EDITOR)", SupportedEditors)
}

// Open launches `editor <path>` in the background and returns once the
// editor has successfully exec'd (we don't wait for it to close — GUI
// editors are typically non-blocking by design but stay attached to the
// parent otherwise).
func Open(editor, path string) error {
	if editor == "" {
		e, err := DefaultEditor()
		if err != nil {
			return err
		}
		editor = e
	}
	if _, err := exec.LookPath(editor); err != nil {
		return fmt.Errorf("%q not on PATH", editor)
	}

	cmd := exec.Command(editor, path)
	// Detach stdio — we don't want the editor's startup noise polluting
	// our output, and we must not hold a handle to its process.
	cmd.Stdout = nil
	cmd.Stderr = nil
	cmd.Stdin = nil
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("start %s: %w", editor, err)
	}
	// Release so the editor survives this process exiting.
	_ = cmd.Process.Release()
	return nil
}
