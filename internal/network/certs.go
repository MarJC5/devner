package network

import (
	"context"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

type CertsStatus struct {
	MkcertInstalled bool
	CARootInstalled bool
	CAROOTPath      string
}

func DetectCerts(ctx context.Context) CertsStatus {
	var st CertsStatus
	if _, err := exec.LookPath("mkcert"); err == nil {
		st.MkcertInstalled = true
	}
	if st.MkcertInstalled {
		if out, err := exec.CommandContext(ctx, "mkcert", "-CAROOT").Output(); err == nil {
			st.CAROOTPath = string(out)
		}
	}
	return st
}

// InstallRootCA runs `mkcert -install` on Linux/Windows directly. On
// macOS the call is delegated to Terminal.app because the Security
// framework refuses to modify the System keychain from a context that
// can't present its own confirmation dialog (seen with
// "SecTrustSettingsSetTrustSettings: The authorization was denied
// since no user interaction was possible" when tried via osascript
// from the GUI). Running the command from a real terminal session
// gives mkcert a proper Aqua context and the Keychain prompt appears
// correctly.
//
// Returns a specific error type (*NeedsTerminalError) when the caller
// should display a "please run this in your terminal" hint instead of
// the raw subprocess failure.
func InstallRootCA(ctx context.Context) error {
	mkcert, err := exec.LookPath("mkcert")
	if err != nil {
		return fmt.Errorf("mkcert not found — install via `brew install mkcert` (macOS) or `apt install libnss3-tools mkcert` (Linux) or choco/scoop (Windows)")
	}

	if runtime.GOOS == "darwin" {
		return openMkcertInTerminalMac(ctx, mkcert)
	}

	cmd := exec.CommandContext(ctx, mkcert, "-install")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("mkcert -install: %w: %s", err, strings.TrimSpace(string(out)))
	}
	return nil
}

// openMkcertInTerminalMac launches the user's Terminal.app with
// `mkcert -install` ready to run. The user sees the command, hits
// Enter (or it auto-runs depending on the AppleScript), types their
// password natively via a Keychain prompt with proper Aqua context,
// and the CA is installed system-wide.
//
// We return a sentinel error so the UI can show a helpful "check your
// terminal" message rather than treating this as a failure.
func openMkcertInTerminalMac(_ context.Context, mkcertBin string) error {
	// Quote the binary path for safe insertion into the shell command
	// that Terminal.app will execute.
	cmdLine := fmt.Sprintf("'%s' -install && echo '' && echo '✓ CA installed. You can close this window.' ; echo ''; read -n 1 -s -r -p 'Press any key to close…'",
		strings.ReplaceAll(mkcertBin, `'`, `'"'"'`))

	// AppleScript: tell Terminal to activate + execute the command in a
	// new tab. The tab stays open after completion so the user can see
	// the result and close when ready.
	applescript := `tell application "Terminal"
	activate
	do script ` + appleScriptQuote(cmdLine) + `
end tell`

	if err := exec.Command("osascript", "-e", applescript).Run(); err != nil {
		return fmt.Errorf("couldn't open Terminal.app: %w", err)
	}
	return &NeedsTerminalError{Command: cmdLine}
}

// NeedsTerminalError signals that the install was handed off to a
// terminal and we can't observe the result from here. Callers should
// treat it as "success, please refresh later".
type NeedsTerminalError struct {
	Command string
}

func (e *NeedsTerminalError) Error() string {
	return "devner opened Terminal.app to run mkcert -install — confirm the Keychain password prompt there, then click Refresh"
}

func appleScriptQuote(s string) string {
	var b strings.Builder
	b.WriteByte('"')
	for _, r := range s {
		switch r {
		case '"', '\\':
			b.WriteByte('\\')
			b.WriteRune(r)
		default:
			b.WriteRune(r)
		}
	}
	b.WriteByte('"')
	return b.String()
}
