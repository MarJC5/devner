package cli

import (
	"fmt"

	"github.com/devner/devner/internal/project"
	"github.com/spf13/cobra"
)

// newDevCmd groups the subcommands that manage per-project dev-time
// processes: either an HMR dev server (Node/Next/Astro/Vite/Nuxt/SvelteKit)
// or an asset watcher running alongside a PHP backend (WordPress /
// Laravel / generic PHP with Vite, Mix, etc.). The concrete mode is
// recorded in the project's dev_mode column at import/create time.
func newDevCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dev",
		Short: "Manage per-project dev servers or asset watchers",
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
			mode := project.DevMode(p.DevMode)
			switch mode {
			case project.DevModeNone:
				return fmt.Errorf("project %q has no dev process configured (library or plain static project). Nothing to start", name)
			case project.DevModeServer:
				if p.DevPort == 0 {
					return fmt.Errorf("project %q is a dev-server project but has no allocated port — run `devner reconcile --apply`", name)
				}
				if command == "" {
					command = project.DefaultCommand(project.Type(p.Type), p.DevPort)
				}
			case project.DevModeWatch:
				if command == "" {
					command = project.WatchCommand(p.DevCommand)
				}
			default:
				return fmt.Errorf("project %q has unknown dev mode %q", name, p.DevMode)
			}
			if err := d.DevServer.Start(cmd.Context(), name, p.DevPort, command); err != nil {
				return fmt.Errorf("start: %w", err)
			}
			if mode == project.DevModeServer {
				fmt.Printf("✓ dev server starting for %s on internal port %d\n", name, p.DevPort)
				fmt.Printf("  URL     : https://%s\n", p.Domain)
			} else {
				fmt.Printf("✓ asset watcher starting for %s (PHP serves https://%s)\n", name, p.Domain)
			}
			fmt.Printf("  Command : %s\n", command)
			fmt.Printf("  Logs    : devner dev logs %s\n", name)
			fmt.Printf("\nNote: first build can take ~10s before changes are picked up.\n")
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
		Short: "Show dev/watch status for one project or all projects with a configured dev mode",
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
					if p.DevMode != "" {
						names = append(names, p.Name)
					}
				}
				if len(names) == 0 {
					fmt.Println("no projects with a configured dev mode")
					return nil
				}
			}
			for _, n := range names {
				p, err := d.Store.GetProject(cmd.Context(), n)
				if err != nil {
					fmt.Printf("  %s  — not in store\n", n)
					continue
				}
				if p.DevMode == "" {
					fmt.Printf("  %-20s  (no dev process — library or plain PHP)\n", n)
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
				switch p.DevMode {
				case string(project.DevModeServer):
					fmt.Printf("  %-20s  mode=server  port=%d  %s  https://%s\n", n, p.DevPort, mark, p.Domain)
				case string(project.DevModeWatch):
					fmt.Printf("  %-20s  mode=watch   cmd=%-8s  %s  https://%s\n", n, p.DevCommand, mark, p.Domain)
				default:
					fmt.Printf("  %-20s  mode=%s  %s\n", n, p.DevMode, mark)
				}
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
