package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)

// Marker enclosing the devner block in rc files. Used to detect an
// already-installed setup and to let the user find/remove it by hand.
const (
	shellMarkerBegin = "# >>> devner shell integration >>>"
	shellMarkerEnd   = "# <<< devner shell integration <<<"
)

func newShellCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "shell",
		Short: "Shell integration helpers (alias, tab completion, PATH)",
	}
	cmd.AddCommand(newShellInitCmd(), newShellUninstallCmd())
	return cmd
}

func newShellInitCmd() *cobra.Command {
	var shell string
	var install bool

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Print or install shell integration (PATH + completion) for zsh, bash, or fish",
		Long: `Emit a shell snippet that:
  - puts the devner binary on PATH (using its current absolute location),
  - wires tab completion for subcommands, flags and argument values.

Without --install, the snippet goes to stdout so you can source it or
eval it. With --install, the snippet is appended to the detected rc file
between marker comments, so re-running is idempotent and the block can
be removed via 'devner shell uninstall'.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			detected, rc, err := resolveShell(shell)
			if err != nil {
				return err
			}
			binPath, err := os.Executable()
			if err != nil {
				return fmt.Errorf("resolve binary path: %w", err)
			}
			binDir := filepath.Dir(binPath)
			snippet := shellSnippet(detected, binDir)

			if !install {
				fmt.Print(snippet)
				return nil
			}
			already, err := rcHasBlock(rc)
			if err != nil {
				return err
			}
			if already {
				fmt.Printf("✓ %s already contains a devner block — nothing to do\n", rc)
				fmt.Printf("  (rerun with `devner shell uninstall` to remove, then `devner shell init --install` to reinstall)\n")
				return nil
			}
			if err := appendBlock(rc, snippet); err != nil {
				return err
			}
			fmt.Printf("✓ appended devner integration to %s\n", rc)
			fmt.Printf("  → reload your shell:  source %s\n", rc)
			return nil
		},
	}
	cmd.Flags().StringVar(&shell, "shell", "auto", "target shell: auto | zsh | bash | fish")
	cmd.Flags().BoolVar(&install, "install", false, "append the snippet to the shell rc file instead of printing")
	return cmd
}

func newShellUninstallCmd() *cobra.Command {
	var shell string
	cmd := &cobra.Command{
		Use:   "uninstall",
		Short: "Remove the devner block from the shell rc file",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, rc, err := resolveShell(shell)
			if err != nil {
				return err
			}
			removed, err := removeBlock(rc)
			if err != nil {
				return err
			}
			if !removed {
				fmt.Printf("— no devner block found in %s\n", rc)
				return nil
			}
			fmt.Printf("✓ removed devner block from %s\n", rc)
			fmt.Printf("  → reload your shell:  source %s\n", rc)
			return nil
		},
	}
	cmd.Flags().StringVar(&shell, "shell", "auto", "target shell: auto | zsh | bash | fish")
	return cmd
}

// resolveShell returns (shellName, rcFilePath, err). When shell is "auto",
// it's inferred from the SHELL env var (basename) with a macOS bash
// fallback to ~/.bash_profile.
func resolveShell(shell string) (string, string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", "", err
	}
	if shell == "" || shell == "auto" {
		shell = filepath.Base(os.Getenv("SHELL"))
	}
	switch shell {
	case "zsh":
		return "zsh", filepath.Join(home, ".zshrc"), nil
	case "bash":
		rc := filepath.Join(home, ".bashrc")
		// On macOS, interactive login shells read ~/.bash_profile, not
		// ~/.bashrc. Prefer whichever the user already has; fall back to
		// ~/.bash_profile on darwin.
		if runtime.GOOS == "darwin" {
			bp := filepath.Join(home, ".bash_profile")
			if _, err := os.Stat(bp); err == nil {
				rc = bp
			} else if _, err := os.Stat(rc); err != nil {
				rc = bp
			}
		}
		return "bash", rc, nil
	case "fish":
		return "fish", filepath.Join(home, ".config", "fish", "config.fish"), nil
	default:
		return "", "", fmt.Errorf("unsupported shell %q (use --shell=zsh|bash|fish)", shell)
	}
}

func shellSnippet(shell, binDir string) string {
	var b strings.Builder
	b.WriteString(shellMarkerBegin + "\n")
	switch shell {
	case "fish":
		fmt.Fprintf(&b, "set -gx PATH %s $PATH\n", binDir)
		b.WriteString("devner completion fish | source\n")
	default: // zsh + bash
		fmt.Fprintf(&b, "export PATH=\"%s:$PATH\"\n", binDir)
		fmt.Fprintf(&b, "command -v devner >/dev/null 2>&1 && source <(devner completion %s)\n", shell)
	}
	b.WriteString(shellMarkerEnd + "\n")
	return b.String()
}

func rcHasBlock(path string) (bool, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return strings.Contains(string(data), shellMarkerBegin), nil
}

func appendBlock(path, block string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	// Ensure exactly one blank line before the block.
	existing, _ := os.ReadFile(path)
	sep := "\n"
	if len(existing) > 0 && !strings.HasSuffix(string(existing), "\n") {
		sep = "\n\n"
	} else if strings.HasSuffix(string(existing), "\n\n") {
		sep = ""
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := f.WriteString(sep + block); err != nil {
		return err
	}
	return nil
}

// removeBlock strips the marked block (inclusive) from the rc file.
// Returns (removed, err).
func removeBlock(path string) (bool, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	s := string(data)
	start := strings.Index(s, shellMarkerBegin)
	end := strings.Index(s, shellMarkerEnd)
	if start < 0 || end < 0 || end < start {
		return false, nil
	}
	endIdx := end + len(shellMarkerEnd)
	// Consume the trailing newline of the marker line if present.
	if endIdx < len(s) && s[endIdx] == '\n' {
		endIdx++
	}
	// Consume the optional leading blank line before the block.
	cutStart := start
	if cutStart > 0 && s[cutStart-1] == '\n' {
		cutStart--
		if cutStart > 0 && s[cutStart-1] == '\n' {
			cutStart--
		}
	}
	newContent := s[:cutStart] + s[endIdx:]
	if err := os.WriteFile(path, []byte(newContent), 0o644); err != nil {
		return false, err
	}
	return true, nil
}
