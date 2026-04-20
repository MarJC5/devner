package llm

import (
	"fmt"
	"os"

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
	switch pc.Kind {
	case "openai_compat":
		if apiKey == "" && name != "ollama" {
			return nil, fmt.Errorf("env var %s is empty (set your %s API key)", pc.APIKeyEnv, name)
		}
		return NewOpenAICompat(name, pc.BaseURL, apiKey, pc.Model), nil
	case "anthropic":
		if apiKey == "" {
			return nil, fmt.Errorf("env var %s is empty", pc.APIKeyEnv)
		}
		return NewAnthropic(name, apiKey, pc.Model), nil
	default:
		return nil, fmt.Errorf("unknown provider kind %q for %s", pc.Kind, name)
	}
}
