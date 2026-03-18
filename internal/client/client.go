package client

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"aicli/internal/config"
	"aicli/internal/types"
	"aicli/internal/ui"
)

// Call отправляет запрос к LocalAI/OpenAI-совместимому серверу с streaming-ответом
func Call(messages []types.Message, model string) (string, error) {
	url := config.Host() + "/v1/chat/completions"

	payload := map[string]interface{}{
		"model":       model,
		"messages":    messages,
		"stream":      true,
		"temperature": 0.7,
		"top_p":       0.9,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if key := config.APIKey(); key != "" {
		req.Header.Set("Authorization", "Bearer "+key)
	}

	// Запускаем спиннер
	stopSpinner := ui.StartSpinner()
	defer stopSpinner() // безопасно благодаря sync.Once

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		stopSpinner()
		return "", fmt.Errorf("network error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		stopSpinner()
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API error %d: %s", resp.StatusCode, string(body))
	}

	var sb strings.Builder
	reader := bufio.NewReader(resp.Body)
	firstChunk := true

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				fmt.Fprintf(os.Stderr, "stream read error: %v\n", err)
			}
			break
		}

		line = strings.TrimSpace(line)
		if line == "" || line == "data: [DONE]" {
			continue
		}
		if !strings.HasPrefix(line, "data: ") {
			continue
		}

		data := strings.TrimPrefix(line, "data: ")
		var chunk types.ChatCompletionChunk
		if err := json.Unmarshal([]byte(data), &chunk); err != nil {
			continue
		}

		if len(chunk.Choices) == 0 {
			continue
		}

		content := chunk.Choices[0].Delta.Content
		if content == "" {
			continue
		}

		// Первый чанк — останавливаем спиннер и выводим заголовок модели
		if firstChunk {
			stopSpinner()
			fmt.Printf("%s%s%s\n %s↳ %s",
				ui.ColorRed+ui.ColorBold, model, ui.ColorReset,
				ui.ColorCyan, ui.ColorReset)
			firstChunk = false
		}

		fmt.Print(content)
		sb.WriteString(content)
	}

	fmt.Println() // новая строка после ответа

	if firstChunk {
		stopSpinner()
		return "", fmt.Errorf("server returned empty response")
	}

	return sb.String(), nil
}

// ListModels возвращает список доступных моделей
func ListModels() ([]string, error) {
	url := config.Host() + "/v1/models"
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("не удалось подключиться к %s: %w", config.Host(), err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("сервер вернул статус %d", resp.StatusCode)
	}

	var list types.ModelList
	if err := json.NewDecoder(resp.Body).Decode(&list); err != nil {
		return nil, err
	}

	var models []string
	for _, m := range list.Data {
		models = append(models, m.ID)
	}
	return models, nil
}
