package cli

import (
	"context"
	"fmt"
	"time"

	"github.com/devner/devner/internal/database"
	"github.com/devner/devner/internal/network"
	"github.com/devner/devner/internal/project"
	"github.com/devner/devner/internal/store"
	"github.com/spf13/cobra"
)

// registerCommands attaches the full Phase 1 command tree to the root.
func registerCommands(root *cobra.Command) {
	root.AddCommand(
		newNewCmd(),
		newRemoveCmd(),
		newPSCmd(),
		newLogsCmd(),
		newRestartCmd(),
		newDBCmd(),
		newHostsCmd(),
		newRebuildCmd(),
		newDeleteCmd(),
		newImportCmd(),
		newReconcileCmd(),
		newCertsCmd(),
		newAgentCmd(),
		newModelsCmd(),
	)
}

// ---- new ----

func newNewCmd() *cobra.Command {
	var dbEngine string
	cmd := &cobra.Command{
		Use:   "new <type> <name>",
		Short: "Create a new project (wordpress|laravel|node|nextjs|astro)",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			typ, err := project.ParseType(args[0])
			if err != nil {
				return err
			}
			name := args[1]
			if err := project.ValidateName(name); err != nil {
				return err
			}

			d, err := buildDeps()
			if err != nil {
				return err
			}
			defer d.Store.Close()

			ctx := cmd.Context()

			fmt.Printf("→ scaffolding %s project %q\n", typ, name)
			if err := d.Project.Scaffold(ctx, typ, name); err != nil {
				return fmt.Errorf("scaffold: %w", err)
			}

			dbName := ""
			if dbEngine != "" {
				if err := database.ValidateName(name); err != nil {
					return fmt.Errorf("project name cannot be used as db name: %w", err)
				}
				fmt.Printf("→ creating %s database %q\n", dbEngine, name)
				creds, err := d.DB.Create(ctx, database.Engine(dbEngine), name)
				if err != nil {
					return fmt.Errorf("db: %w", err)
				}
				dbName = creds.Database
				fmt.Printf("  user=%s host=%s port=%d\n", creds.User, creds.Host, creds.Port)
			}

			domain := name + ".localhost"
			p := store.Project{
				Name:      name,
				Type:      string(typ),
				DBEngine:  dbEngine,
				DBName:    dbName,
				Domain:    domain,
				Path:      d.Config.Stack.ProjectsDir + "/" + name,
				CreatedAt: time.Now(),
			}
			if err := d.Store.UpsertProject(ctx, p); err != nil {
				return fmt.Errorf("store: %w", err)
			}

			if err := applyCaddySites(ctx, d); err != nil {
				return fmt.Errorf("caddy apply: %w", err)
			}
			fmt.Printf("✓ %s ready at https://%s\n", name, domain)
			return nil
		},
	}
	cmd.Flags().StringVar(&dbEngine, "db", "", "database engine: mysql | postgres")
	return cmd
}

// ---- remove ----

func newRemoveCmd() *cobra.Command {
	var keepFiles bool
	cmd := &cobra.Command{
		Use:   "remove <name>",
		Short: "Delete a project (files, database, hosts entry)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			if err := project.ValidateName(name); err != nil {
				return err
			}
			d, err := buildDeps()
			if err != nil {
				return err
			}
			defer d.Store.Close()
			ctx := cmd.Context()

			p, err := d.Store.GetProject(ctx, name)
			if err != nil {
				return fmt.Errorf("project %q not found in store", name)
			}

			if p.DBEngine != "" && p.DBName != "" {
				fmt.Printf("→ dropping %s database %q\n", p.DBEngine, p.DBName)
				if err := d.DB.Drop(ctx, database.Engine(p.DBEngine), p.DBName); err != nil {
					fmt.Printf("  (warning) drop db: %v\n", err)
				}
			}
			if !keepFiles {
				fmt.Printf("→ removing files %s\n", p.Path)
				if err := d.Project.Remove(ctx, name); err != nil {
					fmt.Printf("  (warning) remove files: %v\n", err)
				}
			}
			if err := d.Store.DeleteProject(ctx, name); err != nil {
				return fmt.Errorf("store delete: %w", err)
			}
			if err := applyCaddySites(ctx, d); err != nil {
				fmt.Printf("  (warning) caddy apply: %v\n", err)
			}
			fmt.Printf("✓ removed %s\n", name)
			return nil
		},
	}
	cmd.Flags().BoolVar(&keepFiles, "keep-files", false, "do not delete project files")
	return cmd
}

// ---- ps ----

func newPSCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "ps",
		Short: "Show stack container status",
		RunE: func(cmd *cobra.Command, args []string) error {
			d, err := buildDeps()
			if err != nil {
				return err
			}
			defer d.Store.Close()
			return d.Runtime.PS(cmd.Context())
		},
	}
}

// ---- logs ----

func newLogsCmd() *cobra.Command {
	var tail int
	var follow bool
	cmd := &cobra.Command{
		Use:   "logs <service>",
		Short: "Show logs for a stack service",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			d, err := buildDeps()
			if err != nil {
				return err
			}
			defer d.Store.Close()
			return d.Runtime.Logs(cmd.Context(), args[0], tail, follow)
		},
	}
	cmd.Flags().IntVar(&tail, "tail", 100, "number of lines")
	cmd.Flags().BoolVarP(&follow, "follow", "f", false, "follow log output")
	return cmd
}

// ---- restart ----

func newRestartCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "restart [service]",
		Short: "Restart stack (or one service)",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			d, err := buildDeps()
			if err != nil {
				return err
			}
			defer d.Store.Close()
			svc := ""
			if len(args) == 1 {
				svc = args[0]
			}
			return d.Runtime.Restart(cmd.Context(), svc)
		},
	}
}

// ---- rebuild ----

func newRebuildCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "rebuild",
		Short: "Rebuild images and recreate containers",
		RunE: func(cmd *cobra.Command, args []string) error {
			d, err := buildDeps()
			if err != nil {
				return err
			}
			defer d.Store.Close()
			return d.Runtime.Rebuild(cmd.Context())
		},
	}
}

// ---- delete ----

func newDeleteCmd() *cobra.Command {
	var force bool
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete stack containers and volumes (destructive)",
		RunE: func(cmd *cobra.Command, args []string) error {
			if !force {
				return fmt.Errorf("refusing to delete without --force (all data will be lost)")
			}
			d, err := buildDeps()
			if err != nil {
				return err
			}
			defer d.Store.Close()
			return d.Runtime.Delete(cmd.Context())
		},
	}
	cmd.Flags().BoolVar(&force, "force", false, "confirm destructive delete")
	return cmd
}

// ---- db ----

func newDBCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "db", Short: "Database operations"}
	cmd.AddCommand(
		&cobra.Command{
			Use:   "create <engine> <name>",
			Short: "Create database + user (engine: mysql|postgres)",
			Args:  cobra.ExactArgs(2),
			RunE: func(cmd *cobra.Command, args []string) error {
				d, err := buildDeps()
				if err != nil {
					return err
				}
				defer d.Store.Close()
				creds, err := d.DB.Create(cmd.Context(), database.Engine(args[0]), args[1])
				if err != nil {
					return err
				}
				fmt.Printf("✓ %s/%s\n  user=%s\n  password=%s\n", args[0], creds.Database, creds.User, creds.Password)
				return nil
			},
		},
		&cobra.Command{
			Use:   "drop <engine> <name>",
			Short: "Drop database + user",
			Args:  cobra.ExactArgs(2),
			RunE: func(cmd *cobra.Command, args []string) error {
				d, err := buildDeps()
				if err != nil {
					return err
				}
				defer d.Store.Close()
				return d.DB.Drop(cmd.Context(), database.Engine(args[0]), args[1])
			},
		},
	)
	return cmd
}

// ---- hosts ----

func newHostsCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "hosts", Short: "Manage /etc/hosts entries"}
	cmd.AddCommand(
		&cobra.Command{
			Use:   "add <domain> [target]",
			Short: "Add hosts entry (no-op for *.localhost)",
			Args:  cobra.RangeArgs(1, 2),
			RunE: func(cmd *cobra.Command, args []string) error {
				target := "127.0.0.1"
				if len(args) == 2 {
					target = args[1]
				}
				hm, err := network.NewHostsManager()
				if err != nil {
					return err
				}
				added, err := hm.Add(args[0], target)
				if err != nil {
					return err
				}
				if !added {
					fmt.Printf("no-op: %s resolves natively or already present\n", args[0])
					return nil
				}
				if err := hm.Save(); err != nil {
					return err
				}
				fmt.Printf("✓ %s → %s\n", args[0], target)
				return nil
			},
		},
		&cobra.Command{
			Use:   "remove <domain>",
			Short: "Remove hosts entry",
			Args:  cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				hm, err := network.NewHostsManager()
				if err != nil {
					return err
				}
				removed, err := hm.Remove(args[0])
				if err != nil {
					return err
				}
				if !removed {
					fmt.Printf("no-op: %s not present\n", args[0])
					return nil
				}
				return hm.Save()
			},
		},
	)
	return cmd
}

// applyCaddySites pushes the current list of projects to Caddy admin API.
func applyCaddySites(ctx context.Context, d *Deps) error {
	projects, err := d.Store.ListProjects(ctx)
	if err != nil {
		return err
	}
	var sites []network.ProjectSite
	for _, p := range projects {
		sites = append(sites, network.ProjectSite{
			Domain: p.Domain,
			Root:   project.Type(p.Type).DocRoot(p.Name),
		})
	}
	return d.Caddy.Apply(ctx, sites)
}
