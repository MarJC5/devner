package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

// newHotkeyCmd bundles the OS-level integration helpers that turn
// `devner prompt` into a floating popup bound to a global shortcut.
// Today: iTerm2 Dynamic Profile on macOS. Raycast/Alfred/Hammerspoon
// users don't need this — they wire devner prompt in their launcher UI.
func newHotkeyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "hotkey",
		Short: "Install OS-level helpers so a global shortcut opens `devner prompt`",
	}
	cmd.AddCommand(newHotkeyInstallCmd(), newHotkeyUninstallCmd())
	return cmd
}

const iterm2ProfileName = "Devner Prompt"

func iterm2ProfilePath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, "Library", "Application Support", "iTerm2",
		"DynamicProfiles", "devner-prompt.json")
}

func newHotkeyInstallCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "install",
		Short: "Install the iTerm2 hotkey-window profile (macOS + iTerm2 required)",
		Long: `Creates an iTerm2 Dynamic Profile that auto-runs 'devner prompt'
in a hotkey window (slide-in terminal from the top of the screen).

After running, iTerm2 will pick up the profile within a few seconds.
To wire the Cmd+D shortcut:
  iTerm2 > Settings (Cmd+,) > Keys > Hotkey
    ✓ Show/hide iTerm2 with a system-wide hotkey
    Hotkey: Cmd+D
    Profile for hotkey window: Devner Prompt

Press Cmd+D anywhere — the window slides down with a running agent
session. Press Cmd+D again to hide. Esc inside the TUI closes it fully.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if runtime.GOOS != "darwin" {
				return fmt.Errorf("iTerm2 hotkey install is macOS-only; on %s use Raycast / AutoHotkey / Hammerspoon (see `devner prompt --help`)", runtime.GOOS)
			}
			home, err := os.UserHomeDir()
			if err != nil {
				return err
			}
			itermSupport := filepath.Join(home, "Library", "Application Support", "iTerm2")
			if _, err := os.Stat(itermSupport); err != nil {
				return fmt.Errorf("iTerm2 doesn't appear to be installed (expected %s). Install via `brew install --cask iterm2`", itermSupport)
			}

			binPath, err := exec.LookPath("devner")
			if err != nil {
				// Fall back to current executable so the profile works even
				// when devner isn't globally installed yet.
				binPath, _ = os.Executable()
			}

			guid := uuid.NewString()
			profile := map[string]any{
				"Profiles": []map[string]any{{
					"Name":          iterm2ProfileName,
					"Guid":          guid,
					"Command":       binPath + " prompt",
					"Custom Command": "Yes",
					"Close Sessions On End": true,
					"Working Directory Behavior": 0,
					// Hotkey-window-friendly defaults: rows×cols, fade on blur,
					// tight padding.
					"Rows":             24,
					"Columns":          100,
					"Screen":           -1,
					"Window Type":      4, // Top-of-screen hotkey window
					"Transparency":     0.1,
					"Blur":             true,
					"Blur Radius":      10,
					"Has Hotkey":       true,
					"HotKey Activation": 0,
					"HotKey Characters":    "",        // leave unset; user binds in iTerm2 Settings > Keys > Hotkey
					"HotKey Character Ignoring Modifiers": "",
					"HotKey Modifier Activation": 0,
					"Unlimited Scrollback": true,
				}},
			}

			if err := os.MkdirAll(filepath.Dir(iterm2ProfilePath()), 0o755); err != nil {
				return err
			}
			data, err := json.MarshalIndent(profile, "", "  ")
			if err != nil {
				return err
			}
			if err := os.WriteFile(iterm2ProfilePath(), data, 0o644); err != nil {
				return err
			}
			fmt.Printf("✓ wrote iTerm2 Dynamic Profile: %s\n", iterm2ProfilePath())
			fmt.Printf("  Command: %s prompt\n", binPath)
			fmt.Println()
			fmt.Println("Next — one-time iTerm2 setup (30 seconds):")
			fmt.Println("  1. Open iTerm2 (reopen if it was running — it picks up profiles on launch).")
			fmt.Println("  2. iTerm2 > Settings (Cmd+,) > Keys > Hotkey")
			fmt.Println("     ✓ Check \"Show/hide iTerm2 with a system-wide hotkey\"")
			fmt.Println("     Press Cmd+D (or any combo) in the hotkey field.")
			fmt.Println("     Set \"Profile for hotkey window\" to: Devner Prompt")
			fmt.Println("  3. Press the hotkey anywhere — the window slides down with devner prompt running.")
			return nil
		},
	}
}

func newHotkeyUninstallCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "uninstall",
		Short: "Remove the iTerm2 hotkey-window profile",
		RunE: func(cmd *cobra.Command, args []string) error {
			path := iterm2ProfilePath()
			if err := os.Remove(path); err != nil {
				if os.IsNotExist(err) {
					fmt.Printf("— no profile to remove at %s\n", path)
					return nil
				}
				return err
			}
			fmt.Printf("✓ removed %s\n", path)
			fmt.Println("  You may also want to clear the Cmd+D binding in iTerm2 > Settings > Keys > Hotkey.")
			return nil
		},
	}
}
