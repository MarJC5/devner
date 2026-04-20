package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/devner/devner/internal/project"
	"github.com/devner/devner/internal/store"
	"github.com/spf13/cobra"
)

func newImportCmd() *cobra.Command {
	var source string
	var dryRun bool

	cmd := &cobra.Command{
		Use:   "import",
		Short: "Import projects from a directory (e.g. v1 devner projects/) into the store",
		RunE: func(cmd *cobra.Command, args []string) error {
			if source == "" {
				return fmt.Errorf("--source is required (path to existing projects/ directory)")
			}
			src, err := filepath.Abs(source)
			if err != nil {
				return err
			}
			info, err := os.Stat(src)
			if err != nil {
				return fmt.Errorf("source: %w", err)
			}
			if !info.IsDir() {
				return fmt.Errorf("source is not a directory")
			}

			d, err := buildDeps()
			if err != nil {
				return err
			}
			defer d.Close()
			ctx := cmd.Context()

			entries, err := os.ReadDir(src)
			if err != nil {
				return err
			}

			imported := 0
			skipped := 0
			for _, e := range entries {
				if !e.IsDir() || e.Name()[0] == '.' {
					continue
				}
				projDir := filepath.Join(src, e.Name())
				result, ok := project.Detect(projDir)
				if !ok {
					fmt.Printf("  ? skip %-30s (unrecognized)\n", e.Name())
					skipped++
					continue
				}

				if err := project.ValidateName(e.Name()); err != nil {
					fmt.Printf("  ! skip %-30s (invalid name: %v)\n", e.Name(), err)
					skipped++
					continue
				}

				db := project.SniffDB(projDir, result.Type)

				// Allocate a dev port only for projects that actually
				// want an HTTP dev server. Watch-mode (PHP+assets) and
				// library imports keep DevPort=0 — their routing stays
				// on the PHP `php_server + file_server` path (or is
				// irrelevant, for libraries).
				devPort := 0
				if !dryRun && result.DevMode == project.DevModeServer {
					port, err := d.Store.AllocateDevPort(ctx)
					if err != nil {
						fmt.Printf("  ! skip %-30s (port allocation: %v)\n", e.Name(), err)
						skipped++
						continue
					}
					devPort = port
				}

				p := store.Project{
					Name:       e.Name(),
					Type:       string(result.Type),
					DBEngine:   db.Engine,
					DBName:     db.Name,
					Domain:     e.Name() + ".localhost",
					Path:       projDir,
					CreatedAt:  time.Now(),
					DevPort:    devPort,
					DevMode:    string(result.DevMode),
					DevCommand: result.DevCommand,
				}

				dbSummary := "no-db"
				if db.Name != "" {
					dbSummary = db.Engine + "/" + db.Name
				}
				modeSummary := string(result.DevMode)
				if modeSummary == "" {
					modeSummary = "none"
				}
				if dryRun {
					fmt.Printf("  + would import  %-25s  type=%-10s  mode=%-7s  %s\n", p.Name, p.Type, modeSummary, dbSummary)
					imported++
					continue
				}
				if err := d.Store.UpsertProject(ctx, p); err != nil {
					fmt.Printf("  ! error %-30s (%v)\n", p.Name, err)
					skipped++
					continue
				}
				fmt.Printf("  + imported     %-25s  type=%-10s  mode=%-7s  %s\n", p.Name, p.Type, modeSummary, dbSummary)
				imported++
			}

			fmt.Printf("\n%d imported, %d skipped\n", imported, skipped)
			if !dryRun && imported > 0 {
				// Push the freshly-classified project set to Caddy so
				// newly-imported or re-classified projects start routing
				// without a separate `devner reconcile --apply` call.
				if err := d.ApplyCaddy(ctx); err != nil {
					fmt.Printf("\nWarning: caddy apply failed: %v (run `devner reconcile --apply`)\n", err)
				} else {
					fmt.Printf("\n✓ Caddy config updated\n")
				}
				fmt.Printf("\nNext: `devner up` to start the stack, then `devner list` to verify.\n")
				fmt.Printf("Note: project files stay at %s (not moved).\n", src)
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&source, "source", "", "path to projects directory to import")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "show what would be imported, don't write to store")
	return cmd
}
