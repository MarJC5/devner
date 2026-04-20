package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	Stack StackConfig `mapstructure:"stack"`
	LLM   LLMConfig   `mapstructure:"llm"`
}

type StackConfig struct {
	DataDir     string `mapstructure:"data_dir"`
	ProjectsDir string `mapstructure:"projects_dir"`
}

type LLMConfig struct {
	ActiveProvider string                    `mapstructure:"active_provider"`
	Providers      map[string]ProviderConfig `mapstructure:"providers"`
}

type ProviderConfig struct {
	BaseURL   string `mapstructure:"base_url"`
	APIKeyEnv string `mapstructure:"api_key_env"`
	Model     string `mapstructure:"model"`
	Kind      string `mapstructure:"kind"` // openai_compat | anthropic
}

// Load reads the config file from the OS-standard location
// (~/Library/Application Support/devner/ on macOS, ~/.config/devner/ on
// Linux, %AppData%\devner\ on Windows). If the file does not exist on
// first run, it is materialized from the defaults so the user has a
// discoverable and editable starting point.
func Load() (*Config, error) {
	cfgDir, err := os.UserConfigDir()
	if err != nil {
		return nil, fmt.Errorf("resolve user config dir: %w", err)
	}
	devnerCfgDir := filepath.Join(cfgDir, "devner")
	if err := os.MkdirAll(devnerCfgDir, 0o755); err != nil {
		return nil, fmt.Errorf("create config dir: %w", err)
	}
	cfgPath := filepath.Join(devnerCfgDir, "config.toml")

	home, _ := os.UserHomeDir()

	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("toml")
	v.AddConfigPath(devnerCfgDir)

	applyDefaults(v, home)

	// If the config file doesn't exist, write the defaults once so the
	// user can see and edit it. ViperWriteConfigAs is used rather than
	// WriteConfig so we control the exact path and don't depend on
	// Viper's internal state.
	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		if err := writeDefaultConfig(cfgPath); err != nil {
			return nil, fmt.Errorf("write default config: %w", err)
		}
	}

	if err := v.ReadInConfig(); err != nil {
		if _, notFound := err.(viper.ConfigFileNotFoundError); !notFound {
			return nil, fmt.Errorf("read config: %w", err)
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}
	return &cfg, nil
}

func (c *Config) ConfigPath() string {
	cfgDir, _ := os.UserConfigDir()
	return filepath.Join(cfgDir, "devner", "config.toml")
}

func applyDefaults(v *viper.Viper, home string) {
	v.SetDefault("stack.data_dir", filepath.Join(home, ".devner"))
	v.SetDefault("stack.projects_dir", filepath.Join(home, "devner", "projects"))
	v.SetDefault("llm.active_provider", "infomaniak")

	v.SetDefault("llm.providers.infomaniak.base_url", "https://api.infomaniak.com/1/ai/openai/v1")
	v.SetDefault("llm.providers.infomaniak.api_key_env", "INFOMANIAK_API_KEY")
	v.SetDefault("llm.providers.infomaniak.model", "mixtral")
	v.SetDefault("llm.providers.infomaniak.kind", "openai_compat")

	v.SetDefault("llm.providers.anthropic.api_key_env", "ANTHROPIC_API_KEY")
	v.SetDefault("llm.providers.anthropic.model", "claude-opus-4-7")
	v.SetDefault("llm.providers.anthropic.kind", "anthropic")

	v.SetDefault("llm.providers.ollama.base_url", "http://localhost:11434/v1")
	v.SetDefault("llm.providers.ollama.model", "qwen2.5-coder:14b")
	v.SetDefault("llm.providers.ollama.kind", "openai_compat")
}

// writeDefaultConfig writes a documented TOML file to `path`. Hand-written
// instead of letting Viper serialize because we want comments and a stable
// section order.
func writeDefaultConfig(path string) error {
	home, _ := os.UserHomeDir()
	content := fmt.Sprintf(`# Devner configuration.
# Regenerated automatically on first run — edit freely.

[stack]
# data_dir holds compose.yaml, store.db, materialized Dockerfiles, Caddyfile.
data_dir     = "%s/.devner"
# projects_dir is mounted into frankenphp at /var/www/html/.
# Point this at wherever you keep your WordPress / Laravel / Node projects.
projects_dir = "%s/devner/projects"

[llm]
# Pick which provider answers 'devner agent' and the TUI chat scene.
# Must match one of the [llm.providers.*] names below.
active_provider = "infomaniak"

[llm.providers.infomaniak]
base_url    = "https://api.infomaniak.com/1/ai/openai/v1"
api_key_env = "INFOMANIAK_API_KEY"
model       = "mixtral"
kind        = "openai_compat"

[llm.providers.anthropic]
api_key_env = "ANTHROPIC_API_KEY"
model       = "claude-opus-4-7"
kind        = "anthropic"

[llm.providers.ollama]
# Local Ollama server — no API key needed.
# Run: ollama serve & ; ollama pull qwen2.5-coder:14b
base_url    = "http://localhost:11434/v1"
model       = "qwen2.5-coder:14b"
kind        = "openai_compat"
`, home, home)

	return os.WriteFile(path, []byte(content), 0o644)
}
