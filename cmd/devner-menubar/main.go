// devner-menubar is a tiny standalone process whose only job is to
// own a macOS menubar item (NSStatusItem) for devner. Runs in its own
// process so fyne.io/systray can own the main run loop (Wails v2
// wouldn't yield it otherwise).
//
// Actions shell out to the `devner` CLI. We resolve its path
// ourselves rather than relying on $PATH: users often install devner
// as a shell alias (e.g. in ~/.zshrc), and exec.Command bypasses
// shell aliasing. Looking for a sibling binary in our own bin/ dir
// matches the layout produced by `make build`.
package main

import (
	_ "embed"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"fyne.io/systray"

	"github.com/devner/devner/internal/sysenv"
)

//go:embed tray_icon.png
var trayIcon []byte

// devnerBin is resolved once at startup; nil on failure.
var devnerBin string

// logFile is tail-able at ~/.devner/menubar.log — the only way to see
// what the detached daemon is doing, since stderr goes nowhere once
// `devner menubar` forks us.
func setupLogging() {
	home, err := os.UserHomeDir()
	if err != nil {
		return
	}
	dir := filepath.Join(home, ".devner")
	_ = os.MkdirAll(dir, 0o755)
	f, err := os.OpenFile(filepath.Join(dir, "menubar.log"),
		os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return
	}
	log.SetOutput(f)
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
}

func main() {
	setupLogging()
	// macOS-only: if the daemon is registered as a Login Item, launchd
	// gives it a stripped PATH that excludes /usr/local/bin where the
	// Docker Desktop `docker` symlink lives. statusLoop()'s
	// `docker ps` poll would otherwise always print "docker
	// unavailable".
	sysenv.EnsureDevPath()
	devnerBin = resolveDevnerBin()
	log.Printf("menubar start: devnerBin=%q", devnerBin)
	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetTemplateIcon(trayIcon, trayIcon)
	systray.SetTooltip("Devner — local dev orchestrator")

	mOpen := systray.AddMenuItem("Open Devner", "Launch the main window")
	mStatus := systray.AddMenuItem("Stack: …", "")
	mStatus.Disable()

	systray.AddSeparator()
	mStart := systray.AddMenuItem("Start stack", "devner up")
	mStop := systray.AddMenuItem("Stop stack", "devner down")
	mRestart := systray.AddMenuItem("Restart stack", "devner restart")

	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Quit menubar", "Remove the tray icon")

	go statusLoop(mStatus)

	for {
		select {
		case <-mOpen.ClickedCh:
			runDevner("gui")
		case <-mStart.ClickedCh:
			runDevner("up")
		case <-mStop.ClickedCh:
			runDevner("down")
		case <-mRestart.ClickedCh:
			runDevner("restart")
		case <-mQuit.ClickedCh:
			systray.Quit()
			return
		}
	}
}

func onExit() {}

// runDevner invokes the devner CLI with the given args, detached from
// this process. Logs to stderr if the binary isn't found or fails —
// the user can run `devner menubar` in a terminal to see those
// messages; a GUI dialog would be nicer but isn't worth another
// dependency.
func runDevner(args ...string) {
	log.Printf("runDevner: args=%v bin=%q", args, devnerBin)
	if devnerBin == "" {
		log.Printf("runDevner: devner binary not found")
		fmt.Fprintf(os.Stderr, "devner-menubar: devner binary not found — rebuild with `make build`\n")
		return
	}
	// Capture output so failures aren't silent. `devner gui` is
	// quick — a combined pipe is fine; we read it after the child
	// exits. Long-running subcommands would need pipe + goroutine,
	// but the menubar only ever invokes short one-shots.
	cmd := exec.Command(devnerBin, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("runDevner: %s %v failed: %v — output: %s", devnerBin, args, err, string(out))
		return
	}
	log.Printf("runDevner: %s %v ok — output: %s", devnerBin, args, string(out))
}

// resolveDevnerBin checks the same locations as the CLI's `gui` and
// `menubar` lookups: sibling binary, repo bin/, $PATH.
func resolveDevnerBin() string {
	if env := os.Getenv("DEVNER_BIN"); env != "" {
		if _, err := os.Stat(env); err == nil {
			return env
		}
	}
	self, err := os.Executable()
	if err == nil {
		selfDir := filepath.Dir(self)
		for _, c := range []string{
			filepath.Join(selfDir, "devner"),
			filepath.Join(selfDir, "..", "bin", "devner"),
			filepath.Join(selfDir, "..", "..", "bin", "devner"),
		} {
			if _, err := os.Stat(c); err == nil {
				abs, _ := filepath.Abs(c)
				return abs
			}
		}
	}
	// $PATH last resort.
	if p, err := exec.LookPath("devner"); err == nil {
		return p
	}
	return ""
}

// statusLoop polls docker every 3 s and updates the disabled
// "Stack: N/M running" line. We query docker directly rather than
// shelling out to `devner ps` because devner ps prints a human-
// readable table we'd have to parse; docker ps --format is
// machine-friendly.
func statusLoop(item *systray.MenuItem) {
	update := func() {
		out, err := exec.Command("docker", "ps", "--filter", "name=_devner",
			"--format", "{{.State}}").Output()
		if err != nil {
			item.SetTitle("Stack: docker unavailable")
			return
		}
		running := 0
		total := 0
		s := string(out)
		start := 0
		for i := 0; i < len(s); i++ {
			if s[i] == '\n' {
				line := s[start:i]
				if line != "" {
					total++
					if line == "running" {
						running++
					}
				}
				start = i + 1
			}
		}
		if total < 6 {
			total = 6 // canonical stack size
		}
		item.SetTitle(fmt.Sprintf("Stack: %d/%d running", running, total))
	}
	update()
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		update()
	}
}
