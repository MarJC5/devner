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
				if typ, ok := project.Detect(projDir); ok {
					fmt.Printf("  ? %s: on disk but not in store (type=%s)\n", e.Name(), typ)
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

			fmt.Printf("\n%d store entries missing files, %d on-disk projects not in store\n", missingFS, staleFS)
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
