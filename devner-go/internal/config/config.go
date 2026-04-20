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

func Load() (*Config, error) {
	cfgDir, err := os.UserConfigDir()
	if err != nil {
		return nil, fmt.Errorf("resolve user config dir: %w", err)
	}
	devnerCfgDir := filepath.Join(cfgDir, "devner")
	if err := os.MkdirAll(devnerCfgDir, 0o755); err != nil {
		return nil, fmt.Errorf("create config dir: %w", err)
	}

	home, _ := os.UserHomeDir()

	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("toml")
	v.AddConfigPath(devnerCfgDir)

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
