package chat

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"aicli/internal/config"
	"aicli/internal/types"
)

const sessionsDir = ".aicli/sessions"

func init() {
	_ = os.MkdirAll(sessionsDir, 0755)
}

func safeName(name string) string {
	name = strings.TrimSpace(name)
	name = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '-' || r == '_' {
			return r
		}
		return '_'
	}, name)
	if name == "" {
		return "default"
	}
	return name
}

func Path(name string) string {
	return filepath.Join(sessionsDir, safeName(name)+".json")
}

func Save(history []types.Message, name string) error {
	data, err := json.MarshalIndent(history, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(Path(name), data, 0644)
}

func Load(name string) ([]types.Message, error) {
	data, err := os.ReadFile(Path(name))
	if err != nil {
		return []types.Message{{Role: "system", Content: config.SystemPrompt()}}, nil
	}
	var history []types.Message
	if err := json.Unmarshal(data, &history); err != nil {
		return nil, err
	}
	return history, nil
}

func List() ([]string, error) {
	files, err := filepath.Glob(filepath.Join(sessionsDir, "*.json"))
	if err != nil {
		return nil, err
	}
	var names []string
	for _, f := range files {
		names = append(names, strings.TrimSuffix(filepath.Base(f), ".json"))
	}
	return names, nil
}

func Delete(name string) error {
	path := Path(name)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("чат '%s' не найден", name)
	}
	return os.Remove(path)
}
