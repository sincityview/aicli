package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

const (
	DefaultHost         = "http://localhost:8080"
	DefaultModel        = "qwen_qwen3.5-2b"
	DefaultSystemPrompt = "You are a helpful assistant."
)

func Init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".aicli")

	if home, err := os.UserHomeDir(); err == nil {
		viper.AddConfigPath(filepath.Join(home, ".config", "aicli-tool"))
	}

	viper.SetDefault("host", DefaultHost)
	viper.SetDefault("model", DefaultModel)
	viper.SetDefault("system_prompt", DefaultSystemPrompt)
	viper.SetDefault("api_key", "")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			fmt.Fprintf(os.Stderr, "Warning: failed to read config: %v\n", err)
		}
	}
}

func Host() string         { return viper.GetString("host") }
func Model() string        { return viper.GetString("model") }
func SystemPrompt() string { return viper.GetString("system_prompt") }
func APIKey() string       { return viper.GetString("api_key") }
