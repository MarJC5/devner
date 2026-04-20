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
		Short: "Open a focused chat window for the agent",
		Long:  `Launches a full-screen chat UI that closes on Esc. Run it from a terminal or alias it however you like.`,
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
