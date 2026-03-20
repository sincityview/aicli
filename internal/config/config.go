package config

import (
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

	viper.SetDefault("current_provider", "localai-1")
	viper.SetDefault("providers.localai-1.type", "localai")
	viper.SetDefault("providers.localai-1.host", "http://localhost:8080")
	viper.SetDefault("providers.localai-1.default_model", "qwen_qwen3.5-2b")

	viper.ReadInConfig() // ошибки игнорируем — используем дефолты
}

func CurrentProvider() string           { return viper.GetString("current_provider") }
func ProviderType(name string) string   { return viper.GetString("providers." + name + ".type") }
func ProviderHost(name string) string   { return viper.GetString("providers." + name + ".host") }
func ProviderAPIKey(name string) string { return viper.GetString("providers." + name + ".api_key") }
func ProviderDefaultModel(name string) string {
	return viper.GetString("providers." + name + ".default_model")
}
func SystemPrompt() string                    { return viper.GetString("system_prompt") }
func GetAllProviders() map[string]interface{} { return viper.GetStringMap("providers") }
