package cli

import (
	"fmt"
	"strings"

	"github.com/devner/devner/internal/agent"
	"github.com/devner/devner/internal/llm"
	"github.com/devner/devner/internal/tools"
	"github.com/spf13/cobra"
)

func newAgentCmd() *cobra.Command {
	var providerOverride string
	var autoConfirm bool

	cmd := &cobra.Command{
		Use:   "agent <prompt>",
		Short: "Run one agent turn from the CLI (non-interactive)",
		Long:  "Send a single prompt to the LLM agent. Destructive tools require --yes or they will be cancelled automatically.",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			prompt := strings.Join(args, " ")
			d, err := buildDeps()
			if err != nil {
				return err
			}
			defer d.Close()

			provider, err := llm.FromConfig(d.Config, providerOverride)
			if err != nil {
				return err
			}
			registry := tools.BuildDefault(d)
			loop := agent.New(d, provider, registry)
			loop.AutoConfirm = autoConfirm

			events := make(chan agent.Event, 32)
			go loop.Run(cmd.Context(), nil, prompt, events)

			for ev := range events {
				switch {
				case ev.Err != nil:
					return ev.Err
				case ev.TextDelta != "":
					fmt.Print(ev.TextDelta)
				case ev.ToolCall != nil:
					fmt.Printf("\n⏺ %s(%s)\n", ev.ToolCall.Name, truncateCLI(string(ev.ToolCall.Arguments), 200))
				case ev.ToolResult != nil:
					status := "✓"
					if !ev.ToolResult.OK {
						status = "✗"
					}
					fmt.Printf("  %s %s\n", status, truncateCLI(ev.ToolResult.Content, 500))
				case ev.NeedConfirm != nil:
					if autoConfirm {
						ev.NeedConfirm.Reply <- true
					} else {
						fmt.Fprintf(cmd.ErrOrStderr(), "\n⚠ Destructive tool %s cancelled (pass --yes to auto-approve)\n", ev.NeedConfirm.Tool.Name())
						ev.NeedConfirm.Reply <- false
					}
				case ev.Done:
					fmt.Println()
				}
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&providerOverride, "provider", "", "override active LLM provider (infomaniak|anthropic|ollama)")
	cmd.Flags().BoolVar(&autoConfirm, "yes", false, "auto-approve destructive tools")
	return cmd
}

func truncateCLI(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "…"
}
