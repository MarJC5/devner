package cli

import (
	"context"
	"fmt"

	"github.com/devner/devner/internal/tui"
	"github.com/spf13/cobra"
)

// newPromptCmd launches a focused chat-only TUI. The usual tab/scene
// navigation of `devner tui` is hidden — just an input, the agent's
// reply, and Esc to quit. Designed to be bound to a global OS hotkey
// (Cmd+D on macOS, Ctrl+Shift+D elsewhere) so it pops up over whatever
// you're doing when you want to ask the agent something.
func newPromptCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "prompt",
		Short: "Open a focused chat window for the agent (bind to Cmd+D / Ctrl+Shift+D)",
		Long: `Launches a full-screen chat UI that closes on Esc.
Suggested hotkey bindings:

  macOS   — Shortcuts.app: New Shortcut, Action "Run Shell Script":
              /usr/local/bin/devner prompt
            Then assign Services menu + System Settings > Keyboard > Keyboard
            Shortcuts > Services to Cmd+D.
            Easier: Raycast → "Script Commands" or Alfred workflow
            pointing at 'devner prompt', bind Cmd+D globally.

  Windows — AutoHotkey v2 one-liner:
              #d::Run 'wt.exe devner prompt'
            Save as devner.ahk, add to Startup folder.

  Linux   — GNOME: Settings > Keyboard > View and Customize Shortcuts >
            Custom > command: 'gnome-terminal -- devner prompt', key: <Super>d.
            KDE: System Settings > Shortcuts > Custom Shortcuts.

In all cases, the hotkey launches a new terminal window running
'devner prompt'. The TUI takes over until you press Esc.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			d, err := buildDeps()
			if err != nil {
				return fmt.Errorf("init: %w", err)
			}
			defer d.Close()
			return tui.RunPrompt(context.Background(), d)
		},
	}
}
