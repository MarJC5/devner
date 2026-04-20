package cli

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/devner/devner/internal/project"
	"github.com/spf13/cobra"
)

func newReconcileCmd() *cobra.Command {
	var apply bool
	cmd := &cobra.Command{
		Use:   "reconcile",
		Short: "Reconcile the store with the filesystem (and optionally re-push Caddy config)",
		RunE: func(cmd *cobra.Command, args []string) error {
			d, err := buildDeps()
			if err != nil {
				return err
			}
			defer d.Close()
			ctx := cmd.Context()

			stored, err := d.Store.ListProjects(ctx)
			if err != nil {
				return err
			}

			missingFS := 0
			staleFS := 0
			missingPort := 0
			for _, p := range stored {
				if _, err := os.Stat(p.Path); err != nil {
					fmt.Printf("  ! %s: path missing %s\n", p.Name, p.Path)
					missingFS++
					if apply {
						if err := d.Store.DeleteProject(ctx, p.Name); err != nil {
							fmt.Printf("    delete from store failed: %v\n", err)
						} else {
							fmt.Printf("    → removed from store\n")
						}
					}
					continue
				}
				// Heal dev-server projects that never got a port (imports
				// from before migration 003, or failed allocations).
				if p.DevMode == string(project.DevModeServer) && p.DevPort == 0 {
					fmt.Printf("  ! %s: dev_mode=server but dev_port=0\n", p.Name)
					missingPort++
					if apply {
						port, err := d.Store.AllocateDevPort(ctx)
						if err != nil {
							fmt.Printf("    port allocation failed: %v\n", err)
							continue
						}
						p.DevPort = port
						if err := d.Store.UpsertProject(ctx, p); err != nil {
							fmt.Printf("    update failed: %v\n", err)
							continue
						}
						fmt.Printf("    → allocated port %d\n", port)
					}
				}
			}

			storedSet := map[string]bool{}
			for _, p := range stored {
				storedSet[p.Name] = true
			}

			entries, err := os.ReadDir(d.Config.Stack.ProjectsDir)
			if err != nil && !os.IsNotExist(err) {
				return err
			}
			for _, e := range entries {
				if !e.IsDir() || e.Name()[0] == '.' || storedSet[e.Name()] {
					continue
				}
				projDir := filepath.Join(d.Config.Stack.ProjectsDir, e.Name())
				if result, ok := project.Detect(projDir); ok {
					modeSummary := string(result.DevMode)
					if modeSummary == "" {
						modeSummary = "none"
					}
					fmt.Printf("  ? %s: on disk but not in store (type=%s, mode=%s)\n", e.Name(), result.Type, modeSummary)
					staleFS++
					if apply {
						fmt.Printf("    (use `devner import --source=%s` to register)\n", d.Config.Stack.ProjectsDir)
					}
				}
			}

			if apply {
				fmt.Printf("\n→ pushing Caddy config\n")
				if err := reapplyCaddy(ctx, d); err != nil {
					fmt.Printf("  caddy apply failed: %v\n", err)
				} else {
					fmt.Printf("  ✓ caddy updated (%d sites)\n", len(stored)-missingFS)
				}
			}

			fmt.Printf("\n%d store entries missing files, %d on-disk projects not in store, %d dev-server projects missing a port\n", missingFS, staleFS, missingPort)
			if !apply {
				fmt.Printf("\nRun again with --apply to fix drift.\n")
			}
			return nil
		},
	}
	cmd.Flags().BoolVar(&apply, "apply", false, "actually apply fixes (default: dry-run)")
	return cmd
}

// reapplyCaddy recomputes and pushes the current project list to Caddy admin API.
func reapplyCaddy(ctx context.Context, d *Deps) error {
	return applyCaddySites(ctx, d)
}
