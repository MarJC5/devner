package cli

import (
	"fmt"

	"github.com/devner/devner/internal/network"
	"github.com/spf13/cobra"
)

func newCertsCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "certs", Short: "Manage HTTPS certificates for local dev"}

	cmd.AddCommand(
		&cobra.Command{
			Use:   "status",
			Short: "Show mkcert status and CA trust state",
			RunE: func(cmd *cobra.Command, args []string) error {
				st := network.DetectCerts(cmd.Context())
				if !st.MkcertInstalled {
					fmt.Println("mkcert: ✗ not installed")
					fmt.Println()
					fmt.Println("Devner uses Caddy's internal CA for browsers by default.")
					fmt.Println("Install mkcert if you also need curl/Postman to trust local certs:")
					fmt.Println("  macOS   :  brew install mkcert nss")
					fmt.Println("  Linux   :  apt install libnss3-tools && <install mkcert binary>")
					fmt.Println("  Windows :  choco install mkcert   OR   scoop install mkcert")
					return nil
				}
				fmt.Printf("mkcert: ✓ installed\n")
				if st.CAROOTPath != "" {
					fmt.Printf("CAROOT: %s", st.CAROOTPath)
				}
				return nil
			},
		},
		&cobra.Command{
			Use:   "install",
			Short: "Install mkcert's CA into the system trust store (may prompt for sudo)",
			RunE: func(cmd *cobra.Command, args []string) error {
				if err := network.InstallRootCA(cmd.Context()); err != nil {
					return err
				}
				fmt.Println("✓ mkcert CA installed")
				return nil
			},
		},
	)
	return cmd
}
