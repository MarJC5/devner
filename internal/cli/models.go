package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/devner/devner/internal/config"
	"github.com/spf13/cobra"
)

// newModelsCmd hits the active provider's /models endpoint and prints the
// list. Handy as a smoke test after setting up API keys + product_id.
func newModelsCmd() *cobra.Command {
	var providerOverride string
	cmd := &cobra.Command{
		Use:   "models",
		Short: "List available models from the active LLM provider",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return err
			}
			name := cfg.LLM.ActiveProvider
			if providerOverride != "" {
				name = providerOverride
			}
			pc, ok := cfg.LLM.Providers[name]
			if !ok {
				return fmt.Errorf("provider %q not configured", name)
			}
			if pc.Kind != "openai_compat" {
				return fmt.Errorf("models listing only supported for openai_compat providers (got %s)", pc.Kind)
			}

			base := pc.BaseURL
			if strings.Contains(base, "{product_id}") {
				if pc.ProductID == "" {
					return fmt.Errorf("provider %q requires product_id in config", name)
				}
				base = strings.ReplaceAll(base, "{product_id}", pc.ProductID)
			}
			url := strings.TrimRight(base, "/") + "/models"

			req, _ := http.NewRequestWithContext(cmd.Context(), http.MethodGet, url, nil)
			key := pc.APIKey
			if key == "" && pc.APIKeyEnv != "" {
				key = os.Getenv(pc.APIKeyEnv)
			}
			if key != "" {
				req.Header.Set("Authorization", "Bearer "+key)
			}
			client := &http.Client{Timeout: 15 * time.Second}
			resp, err := client.Do(req)
			if err != nil {
				return fmt.Errorf("fetch: %w", err)
			}
			defer resp.Body.Close()

			body, _ := io.ReadAll(resp.Body)
			if resp.StatusCode >= 300 {
				return fmt.Errorf("HTTP %d: %s", resp.StatusCode, truncateCLI(string(body), 400))
			}

			var out struct {
				Data []struct {
					ID      string `json:"id"`
					OwnedBy string `json:"owned_by"`
				} `json:"data"`
			}
			if err := json.Unmarshal(body, &out); err != nil {
				return fmt.Errorf("parse: %w\n---\n%s", err, string(body))
			}
			if len(out.Data) == 0 {
				fmt.Println("(no models returned)")
				return nil
			}
			fmt.Printf("Provider: %s    URL: %s\n\n", name, url)
			for _, m := range out.Data {
				fmt.Printf("  %-45s  %s\n", m.ID, m.OwnedBy)
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&providerOverride, "provider", "", "override active LLM provider")
	return cmd
}
