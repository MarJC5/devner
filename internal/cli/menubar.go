package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)

// The menubar tray is a separate daemon (cmd/devner-menubar) that
// owns the macOS NSStatusItem. Keeping it out of the Wails process
// sidesteps fyne.io/systray's main-thread requirement, which the Wails
// v2 event loop refuses to yield.
//
// `devner menubar` launches the daemon. `devner menubar stop` kills
// it by name. Only one instance per user is expected.

func newMenubarCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "menubar",
		Short: "Launch the devner menubar icon (macOS top-right status item)",
	}
	cmd.AddCommand(newMenubarStartCmd(), newMenubarStopCmd(), newMenubarStatusCmd())
	// Running `devner menubar` directly is the common case — treat it
	// as the start alias so users don't have to remember the subword.
	cmd.RunE = newMenubarStartCmd().RunE
	return cmd
}

func newMenubarStartCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "start",
		Short: "Launch the menubar daemon in background",
		RunE: func(cmd *cobra.Command, args []string) error {
			if isMenubarRunning() {
				fmt.Println("menubar already running.")
				return nil
			}
			bin, err := findMenubarBin()
			if err != nil {
				return err
			}
			c := exec.Command(bin)
			c.Stdout = nil
			c.Stderr = nil
			c.Stdin = nil
			if err := c.Start(); err != nil {
				return fmt.Errorf("launch %s: %w", bin, err)
			}
			_ = c.Process.Release()
			fmt.Printf("✓ menubar launched (pid=%d). Look top-right of your screen.\n", c.Process.Pid)
			return nil
		},
	}
}

func newMenubarStopCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "stop",
		Short: "Stop the running menubar daemon",
		RunE: func(cmd *cobra.Command, args []string) error {
			n, err := killMenubar()
			if err != nil {
				return err
			}
			if n == 0 {
				fmt.Println("no menubar daemon was running.")
				return nil
			}
			fmt.Printf("✓ stopped %d menubar process(es).\n", n)
			return nil
		},
	}
}

func newMenubarStatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show whether the menubar daemon is running",
		RunE: func(cmd *cobra.Command, args []string) error {
			if isMenubarRunning() {
				fmt.Println("● menubar is running")
			} else {
				fmt.Println("○ menubar is not running")
			}
			return nil
		},
	}
}

// findMenubarBin searches the usual places for the devner-menubar
// binary. Similar lookup logic to findGUI (gui.go) to keep UX
// predictable across both commands.
func findMenubarBin() (string, error) {
	if env := os.Getenv("DEVNER_MENUBAR"); env != "" {
		if _, err := os.Stat(env); err == nil {
			return env, nil
		}
	}

	exe := "devner-menubar"
	if runtime.GOOS == "windows" {
		exe += ".exe"
	}

	self, _ := os.Executable()
	selfDir := filepath.Dir(self)
	candidates := []string{
		filepath.Join(selfDir, exe),
		filepath.Join(selfDir, "..", "bin", exe),
	}
	for _, c := range candidates {
		if _, err := os.Stat(c); err == nil {
			return c, nil
		}
	}
	if p, err := exec.LookPath(exe); err == nil {
		return p, nil
	}
	return "", fmt.Errorf("devner-menubar binary not found. Build with `make build` in the repo or set $DEVNER_MENUBAR")
}

// isMenubarRunning uses `pgrep -x` to check for the binary name. Works
// on macOS and Linux; on Windows we return false (the menubar daemon
// isn't packaged there anyway).
func isMenubarRunning() bool {
	if runtime.GOOS == "windows" {
		return false
	}
	out, err := exec.Command("pgrep", "-x", "devner-menubar").Output()
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(out)) != ""
}

// killMenubar sends SIGTERM to every devner-menubar process. Returns
// the number killed. Using pkill avoids walking /proc ourselves.
func killMenubar() (int, error) {
	if runtime.GOOS == "windows" {
		return 0, fmt.Errorf("menubar stop is macOS/Linux only")
	}
	if !isMenubarRunning() {
		return 0, nil
	}
	// pkill returns 0 when at least one process was killed, 1 when
	// none matched. We already checked, so expect success.
	if err := exec.Command("pkill", "-TERM", "-x", "devner-menubar").Run(); err != nil {
		return 0, fmt.Errorf("pkill: %w", err)
	}
	// pgrep count — simpler than parsing pkill output.
	out, _ := exec.Command("pgrep", "-c", "devner-menubar").Output()
	_ = out
	return 1, nil
}
