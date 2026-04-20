package llm

import (
	"fmt"
	"os"
	"strings"

	"github.com/devner/devner/internal/config"
)

// FromConfig builds the active provider from the Config. Callers can also
// pass a specific provider name (useful for per-session overrides).
func FromConfig(cfg *config.Config, override string) (Provider, error) {
	name := cfg.LLM.ActiveProvider
	if override != "" {
		name = override
	}
	pc, ok := cfg.LLM.Providers[name]
	if !ok {
		return nil, fmt.Errorf("provider %q not defined in config", name)
	}
	apiKey := os.Getenv(pc.APIKeyEnv)

	baseURL, err := resolveBaseURL(name, pc)
	if err != nil {
		return nil, err
	}

	switch pc.Kind {
	case "openai_compat":
		if apiKey == "" && name != "ollama" {
			return nil, fmt.Errorf("env var %s is empty (export your %s API key and retry)", pc.APIKeyEnv, name)
		}
		return NewOpenAICompat(name, baseURL, apiKey, pc.Model), nil
	case "anthropic":
		if apiKey == "" {
			return nil, fmt.Errorf("env var %s is empty", pc.APIKeyEnv)
		}
		return NewAnthropic(name, apiKey, pc.Model), nil
	default:
		return nil, fmt.Errorf("unknown provider kind %q for %s", pc.Kind, name)
	}
}

// resolveBaseURL substitutes provider-specific placeholders in BaseURL.
// Currently handles {product_id} for Infomaniak-style URLs.
func resolveBaseURL(name string, pc config.ProviderConfig) (string, error) {
	base := pc.BaseURL
	if strings.Contains(base, "{product_id}") {
		if pc.ProductID == "" {
			return "", fmt.Errorf("provider %q requires product_id in config (base_url contains {product_id})", name)
		}
		base = strings.ReplaceAll(base, "{product_id}", pc.ProductID)
	}
	return base, nil
}
