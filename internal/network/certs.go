package network

import (
	"context"
	"fmt"
	"os/exec"
)

// CertsStatus reports whether mkcert is installed and its CA root is trusted.
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

// InstallRootCA runs `mkcert -install` to add the local CA to the system
// trust store. Requires sudo on macOS/Linux; Keychain/pass prompt appears.
func InstallRootCA(ctx context.Context) error {
	if _, err := exec.LookPath("mkcert"); err != nil {
		return fmt.Errorf("mkcert not found — install via `brew install mkcert` (macOS) or `apt install libnss3-tools mkcert` (Linux) or choco/scoop (Windows)")
	}
	cmd := exec.CommandContext(ctx, "mkcert", "-install")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("mkcert -install: %w: %s", err, out)
	}
	return nil
}
