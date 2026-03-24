package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

func Init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".aicli")

	if home, _ := os.UserHomeDir(); home != "" {
		viper.AddConfigPath(filepath.Join(home, ".config", "aicli"))
	}

	// Дефолтные значения
	viper.SetDefault("current_provider", "localai")
	viper.SetDefault("providers.localai.type", "localai")
	viper.SetDefault("providers.localai.host", "http://localhost:8080")
	viper.SetDefault("providers.localai.default_model", "qwen_qwen3.5-2b")
	viper.SetDefault("providers.localai.context_length", 131072)

	viper.SetDefault("providers.ollama.type", "ollama")
	viper.SetDefault("providers.ollama.host", "http://127.0.0.1:11434")
	viper.SetDefault("providers.ollama.default_model", "qwen3.5:0.8b")
	viper.SetDefault("providers.ollama.context_length", 32768)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			fmt.Fprintf(os.Stderr, "Warning: config read error: %v\n", err)
		}
	}
}

func CurrentProvider() string {
	return viper.GetString("current_provider")
}

func ProviderType(name string) string {
	return viper.GetString("providers." + name + ".type")
}

func ProviderHost(name string) string {
	return viper.GetString("providers." + name + ".host")
}

func ProviderAPIKey(name string) string {
	return viper.GetString("providers." + name + ".api_key")
}

func ProviderDefaultModel(name string) string {
	m := viper.GetString("providers." + name + ".default_model")
	if m == "" {
		return "qwen_qwen3.5-2b"
	}
	return m
}

// ProviderContextLength возвращает максимальную длину контекста для провайдера
func ProviderContextLength(name string) int {
	length := viper.GetInt("providers." + name + ".context_length")
	if length > 0 {
		return length
	}
	return 8192 // fallback по умолчанию
}

func SystemPrompt() string {
	return viper.GetString("system_prompt")
}

func GetAllProviders() map[string]interface{} {
	return viper.GetStringMap("providers")
}
