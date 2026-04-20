package cli

import (
	"fmt"

	"github.com/devner/devner/internal/project"
	"github.com/spf13/cobra"
)

// newOpenCmd + newCodeCmd + newZedCmd + newCursorCmd share the same
// underlying handler. The named shortcuts exist because v1 users are
// used to `devner code <project>` and we keep the muscle memory.

func newOpenCmd() *cobra.Command {
	var editor string
	cmd := &cobra.Command{
		Use:   "open <project>",
		Short: "Open a project in an editor (auto-detects code / cursor / zed / subl / jetbrains)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return openProject(cmd, args[0], editor)
		},
	}
	cmd.Flags().StringVar(&editor, "editor", "", "editor binary name (code|cursor|zed|subl|webstorm|phpstorm|idea). Default: first found.")
	return cmd
}

func newCodeCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "code <project>",
		Short: "Open a project in Visual Studio Code (alias of `open --editor=code`)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return openProject(cmd, args[0], "code")
		},
	}
}

func newZedCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "zed <project>",
		Short: "Open a project in Zed (alias of `open --editor=zed`)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return openProject(cmd, args[0], "zed")
		},
	}
}

func newCursorCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "cursor <project>",
		Short: "Open a project in Cursor (alias of `open --editor=cursor`)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return openProject(cmd, args[0], "cursor")
		},
	}
}

func openProject(cmd *cobra.Command, name, editor string) error {
	if err := project.ValidateName(name); err != nil {
		return err
	}
	d, err := buildDeps()
	if err != nil {
		return err
	}
	defer d.Close()

	p, err := d.Store.GetProject(cmd.Context(), name)
	if err != nil {
		return fmt.Errorf("project %q not found in store (run `devner list` to see known projects)", name)
	}
	if editor == "" {
		resolved, err := project.DefaultEditor()
		if err != nil {
			return err
		}
		editor = resolved
	}
	if err := project.Open(editor, p.Path); err != nil {
		return err
	}
	fmt.Printf("✓ opened %s in %s\n  → %s\n", p.Name, editor, p.Path)
	return nil
}
