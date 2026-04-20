package cli

import (
	"fmt"

	"github.com/devner/devner/internal/project"
	"github.com/spf13/cobra"
)

// newDevCmd groups the subcommands that manage per-project Node / Next /
// Astro / Vite dev servers. Mirror of the agent tools start_dev_server /
// stop_dev_server / dev_server_status — consistency with other command
// groups (`devner db …`, `devner hosts …`).
func newDevCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dev",
		Short: "Manage per-project Node / Next / Astro / Vite dev servers",
	}
	cmd.AddCommand(
		newDevStartCmd(),
		newDevStopCmd(),
		newDevStatusCmd(),
		newDevLogsCmd(),
	)
	return cmd
}

func newDevStartCmd() *cobra.Command {
	var command string
	cmd := &cobra.Command{
		Use:   "start <project>",
		Short: "Start the project's dev server inside frankenphp",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			d, err := buildDeps()
			if err != nil {
				return err
			}
			defer d.Close()
			p, err := d.Store.GetProject(cmd.Context(), name)
			if err != nil {
				return fmt.Errorf("project %q not found (run `devner list`)", name)
			}
			if p.DevPort == 0 {
				return fmt.Errorf("project %q has no dev port (not a Node/Next/Astro/Vite project)", name)
			}
			if command == "" {
				command = project.DefaultCommand(project.Type(p.Type), p.DevPort)
			}
			if err := d.DevServer.Start(cmd.Context(), name, p.DevPort, command); err != nil {
				return fmt.Errorf("start: %w", err)
			}
			fmt.Printf("✓ dev server starting for %s on internal port %d\n", name, p.DevPort)
			fmt.Printf("  URL     : https://%s\n", p.Domain)
			fmt.Printf("  Command : %s\n", command)
			fmt.Printf("  Logs    : devner dev logs %s\n", name)
			fmt.Printf("\nNote: first build (Next/Astro) can take ~10s before the URL responds.\n")
			return nil
		},
	}
	cmd.Flags().StringVar(&command, "command", "", "dev command (default: pnpm dev / npm run start depending on type)")
	return cmd
}

func newDevStopCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "stop <project>",
		Short: "Stop the project's dev server",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			d, err := buildDeps()
			if err != nil {
				return err
			}
			defer d.Close()
			if _, err := d.Store.GetProject(cmd.Context(), name); err != nil {
				return fmt.Errorf("project %q not found", name)
			}
			if err := d.DevServer.Stop(cmd.Context(), name); err != nil {
				return fmt.Errorf("stop: %w", err)
			}
			fmt.Printf("✓ dev server stopped for %s\n", name)
			return nil
		},
	}
}

func newDevStatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status [project]",
		Short: "Show dev server status (for one project, or all Node-like projects)",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			d, err := buildDeps()
			if err != nil {
				return err
			}
			defer d.Close()

			var names []string
			if len(args) == 1 {
				names = []string{args[0]}
			} else {
				all, err := d.Store.ListProjects(cmd.Context())
				if err != nil {
					return err
				}
				for _, p := range all {
					if p.DevPort > 0 {
						names = append(names, p.Name)
					}
				}
				if len(names) == 0 {
					fmt.Println("no Node / Next / Astro / Vite projects")
					return nil
				}
			}
			for _, n := range names {
				p, err := d.Store.GetProject(cmd.Context(), n)
				if err != nil {
					fmt.Printf("  %s  — not in store\n", n)
					continue
				}
				if p.DevPort == 0 {
					fmt.Printf("  %-20s  (PHP project, no dev server)\n", n)
					continue
				}
				st, err := d.DevServer.Status(cmd.Context(), n)
				if err != nil {
					fmt.Printf("  %s  — status error: %v\n", n, err)
					continue
				}
				mark := "○ stopped"
				if st.Running {
					mark = fmt.Sprintf("● running (pid=%d)", st.PID)
				}
				fmt.Printf("  %-20s  port=%d  %s  https://%s\n", n, p.DevPort, mark, p.Domain)
			}
			return nil
		},
	}
}

func newDevLogsCmd() *cobra.Command {
	var tail int
	cmd := &cobra.Command{
		Use:   "logs <project>",
		Short: "Show the last N lines of the dev server log",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			d, err := buildDeps()
			if err != nil {
				return err
			}
			defer d.Close()
			out, err := d.DevServer.Logs(cmd.Context(), name, tail)
			if err != nil {
				return err
			}
			fmt.Print(out)
			return nil
		},
	}
	cmd.Flags().IntVar(&tail, "tail", 100, "number of lines")
	return cmd
}
