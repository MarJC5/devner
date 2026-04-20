package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/spf13/cobra"
)

// newGuiCmd launches the Wails GUI binary. We resolve the binary by
// looking in these places, in order, so the command works in all the
// common layouts:
//
//   1. $DEVNER_GUI               — explicit override
//   2. Next to the `devner` binary
//        - /path/bin/devner-gui
//        - /path/bin/devner-gui.app (macOS)
//   3. ./gui/build/bin/ relative to the devner binary (dev mode)
//   4. /Applications/Devner.app (macOS install)
//   5. PATH: devner-gui or Devner.app
//
// When a .app bundle is found on macOS we use `open` so the OS
// handles dock activation + proper environment; otherwise we fork
// the binary directly.
func newGuiCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "gui",
		Short: "Launch the devner GUI (Wails desktop app)",
		Long: `Opens the devner-gui desktop app. The CLI and GUI share the
same backend (app.Deps), so launching the GUI while the CLI's stack is
running is safe — both see the same store.db and Docker state.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			target, err := findGUI()
			if err != nil {
				return err
			}
			return launchGUI(target)
		},
	}
}

// findGUI returns the path to the GUI binary or .app bundle that we
// should launch. Empty string + error if nothing is usable.
func findGUI() (string, error) {
	if env := os.Getenv("DEVNER_GUI"); env != "" {
		if _, err := os.Stat(env); err == nil {
			return env, nil
		}
	}

	self, _ := os.Executable()
	selfDir := filepath.Dir(self)

	candidates := []string{}

	if runtime.GOOS == "darwin" {
		// Prefer the .app bundle (proper Aqua activation). The current
		// Wails config produces "Devner.app"; earlier builds produced
		// "devner-gui.app" — check both so stale installs still work.
		candidates = append(candidates,
			filepath.Join(selfDir, "Devner.app"),
			filepath.Join(selfDir, "devner-gui.app"),
			filepath.Join(selfDir, "..", "gui", "build", "bin", "Devner.app"),
			filepath.Join(selfDir, "..", "gui", "build", "bin", "devner-gui.app"),
			"/Applications/Devner.app",
			"/Applications/devner-gui.app",
		)
	}

	// Plain executable fallbacks (Linux / Windows + Mac binary outside
	// a bundle). The Wails outputfilename is currently "Devner"; the
	// older "devner-gui" name is kept for compat.
	exeNames := []string{"Devner", "devner-gui"}
	for _, exe := range exeNames {
		if runtime.GOOS == "windows" {
			exe += ".exe"
		}
		candidates = append(candidates,
			filepath.Join(selfDir, exe),
			filepath.Join(selfDir, "..", "gui", "build", "bin", exe),
		)
	}

	for _, c := range candidates {
		if _, err := os.Stat(c); err == nil {
			return c, nil
		}
	}

	// Last resort: PATH.
	for _, name := range exeNames {
		if p, err := exec.LookPath(name); err == nil {
			return p, nil
		}
	}

	return "", fmt.Errorf("devner GUI not found. Build it with `make gui-build` in the repo, or set $DEVNER_GUI to the binary/bundle path")
}

// launchGUI spawns the GUI detached so the CLI shell prompt returns
// immediately.
func launchGUI(target string) error {
	var cmd *exec.Cmd
	if runtime.GOOS == "darwin" && filepath.Ext(target) == ".app" {
		// `open -a` uses LaunchServices so the app activates properly
		// and a dock icon appears / menu bar wires up.
		cmd = exec.Command("open", target)
	} else {
		cmd = exec.Command(target)
		// Detach from our stdio so the GUI's output doesn't garble
		// this shell.
		cmd.Stdout = nil
		cmd.Stderr = nil
		cmd.Stdin = nil
	}
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("launch %s: %w", target, err)
	}
	_ = cmd.Process.Release()
	fmt.Printf("✓ launching %s\n", filepath.Base(target))
	return nil
}
