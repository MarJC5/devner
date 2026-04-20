package cli

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/term"
	"github.com/devner/devner/internal/agent"
	"github.com/devner/devner/internal/llm"
	"github.com/devner/devner/internal/tools"
	"github.com/spf13/cobra"
)

func newAgentCmd() *cobra.Command {
	var providerOverride string
	var autoConfirm bool
	var rawOutput bool

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

			// Piped output or --raw: skip all styling so grep/jq/etc. work.
			isTTY := !rawOutput && term.IsTerminal(os.Stdout.Fd())
			renderer := newAgentRenderer(isTTY)

			events := make(chan agent.Event, 32)
			go loop.Run(cmd.Context(), nil, prompt, events)

			for ev := range events {
				switch {
				case ev.Err != nil:
					return ev.Err
				case ev.TextDelta != "":
					renderer.bufferText(ev.TextDelta)
				case ev.ToolCall != nil:
					renderer.flushText()
					renderer.toolCall(ev.ToolCall.Name, string(ev.ToolCall.Arguments))
				case ev.ToolResult != nil:
					renderer.toolResult(ev.ToolResult.OK, ev.ToolResult.Content)
				case ev.NeedConfirm != nil:
					if autoConfirm {
						ev.NeedConfirm.Reply <- true
					} else {
						renderer.cancelled(ev.NeedConfirm.Tool.Name())
						ev.NeedConfirm.Reply <- false
					}
				case ev.Done:
					renderer.flushText()
				}
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&providerOverride, "provider", "", "override active LLM provider (infomaniak|anthropic|ollama)")
	cmd.Flags().BoolVar(&autoConfirm, "yes", false, "auto-approve destructive tools")
	cmd.Flags().BoolVar(&rawOutput, "raw", false, "disable markdown rendering and styling (for piping)")
	return cmd
}

// agentRenderer formats the agent event stream for a human terminal:
// - tool calls and results get colored prefixes,
// - assistant text is buffered and rendered through glamour so markdown
//   (bold, tables, code fences) actually looks like markdown, and
// - when stdout isn't a TTY, everything falls back to plain text.
type agentRenderer struct {
	tty      bool
	textBuf  strings.Builder
	gRender  *glamour.TermRenderer
	toolSty  lipgloss.Style
	okSty    lipgloss.Style
	errSty   lipgloss.Style
	warnSty  lipgloss.Style
}

// Strip <think>...</think> blocks that some reasoning models (qwen3,
// gpt-oss, Kimi) emit. These are internal chain-of-thought, not intended
// for the user, and they wreck markdown rendering.
var thinkBlockRe = regexp.MustCompile(`(?s)<think>.*?</think>\s*`)

func newAgentRenderer(tty bool) *agentRenderer {
	r := &agentRenderer{tty: tty}
	if tty {
		width := 100
		if w, _, err := term.GetSize(os.Stdout.Fd()); err == nil && w > 20 {
			width = w - 2
		}
		gr, err := glamour.NewTermRenderer(
			glamour.WithAutoStyle(),
			glamour.WithWordWrap(width),
		)
		if err == nil {
			r.gRender = gr
		}
		r.toolSty = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFCC00")).Bold(true)
		r.okSty = lipgloss.NewStyle().Foreground(lipgloss.Color("#00C853"))
		r.errSty = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF5577"))
		r.warnSty = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF5577")).Bold(true)
	}
	return r
}

func (r *agentRenderer) bufferText(delta string) {
	r.textBuf.WriteString(delta)
}

func (r *agentRenderer) flushText() {
	if r.textBuf.Len() == 0 {
		return
	}
	text := thinkBlockRe.ReplaceAllString(r.textBuf.String(), "")
	text = strings.TrimSpace(text)
	r.textBuf.Reset()
	if text == "" {
		return
	}
	if r.gRender != nil {
		if out, err := r.gRender.Render(text); err == nil {
			fmt.Print(out)
			return
		}
	}
	// Fallback — non-TTY or renderer failed.
	fmt.Println(text)
}

func (r *agentRenderer) toolCall(name, args string) {
	shortArgs := truncateCLI(args, 200)
	if r.tty {
		fmt.Printf("%s %s%s\n",
			r.toolSty.Render("⏺"),
			r.toolSty.Render(name),
			lipgloss.NewStyle().Faint(true).Render("("+shortArgs+")"),
		)
		return
	}
	fmt.Printf("⏺ %s(%s)\n", name, shortArgs)
}

func (r *agentRenderer) toolResult(ok bool, content string) {
	short := truncateCLI(strings.TrimSpace(content), 500)
	if r.tty {
		mark := r.okSty.Render("  ✓ ")
		if !ok {
			mark = r.errSty.Render("  ✗ ")
		}
		fmt.Print(mark)
		fmt.Println(lipgloss.NewStyle().Faint(true).Render(short))
		return
	}
	mark := "  ✓ "
	if !ok {
		mark = "  ✗ "
	}
	fmt.Printf("%s%s\n", mark, short)
}

func (r *agentRenderer) cancelled(tool string) {
	msg := fmt.Sprintf("⚠ Destructive tool %s cancelled (pass --yes to auto-approve)", tool)
	if r.tty {
		fmt.Fprintln(os.Stderr, r.warnSty.Render(msg))
		return
	}
	fmt.Fprintln(os.Stderr, msg)
}

func truncateCLI(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "…"
}
