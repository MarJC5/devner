package cli

import (
	"fmt"

	"github.com/devner/devner/internal/platform"
	"github.com/spf13/cobra"
)

func NewRootCmd(version string) *cobra.Command {
	root := &cobra.Command{
		Use:           "devner",
		Short:         "Devner - local dev environment orchestrator",
		Long:          "Devner manages local WordPress, Laravel, Node, Next, and Astro projects with Docker.",
		Version:       version,
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	root.AddCommand(
		newVersionCmd(version),
		newTUICmd(),
		newUpCmd(),
		newDownCmd(),
		newListCmd(),
	)
	registerCommands(root)

	return root
}

func newVersionCmd(version string) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("devner %s (%s)\n", version, platform.Detect())
		},
	}
}

func newListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List managed projects",
		RunE: func(cmd *cobra.Command, args []string) error {
			d, err := buildDeps()
			if err != nil {
				return err
			}
			defer d.Store.Close()
			projects, err := d.Store.ListProjects(cmd.Context())
			if err != nil {
				return err
			}
			if len(projects) == 0 {
				fmt.Printf("No projects yet. Use `devner new <type> <name>` to create one.\n")
				fmt.Printf("Projects dir: %s\n", d.Config.Stack.ProjectsDir)
				return nil
			}
			for _, p := range projects {
				db := "-"
				if p.DBEngine != "" {
					db = p.DBEngine + "/" + p.DBName
				}
				fmt.Printf("%-20s  %-10s  %-15s  https://%s\n", p.Name, p.Type, db, p.Domain)
			}
			return nil
		},
	}
}

func newUpCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "up",
		Short: "Start the shared dev stack (mysql, postgres, redis, frankenphp, ...)",
		RunE: func(cmd *cobra.Command, args []string) error {
			d, err := buildDeps()
			if err != nil {
				return err
			}
			defer d.Store.Close()
			if err := d.Runtime.Up(cmd.Context()); err != nil {
				return err
			}
			// Rewrite the Caddyfile so infra routes (adminer.localhost,
			// mailpit.localhost) and any existing project sites are
			// applied on boot — not only after `devner new`.
			return d.ApplyCaddy(cmd.Context())
		},
	}
}

func newDownCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "down",
		Short: "Stop the shared dev stack",
		RunE: func(cmd *cobra.Command, args []string) error {
			d, err := buildDeps()
			if err != nil {
				return err
			}
			defer d.Store.Close()
			return d.Runtime.Down(cmd.Context())
		},
	}
}

func newTUICmd() *cobra.Command {
	return &cobra.Command{
		Use:   "tui",
		Short: "Launch the interactive TUI",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runTUI()
		},
	}
}
