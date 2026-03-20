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

func getBaseURL(provider string) string {
	host := config.ProviderHost(provider)
	if host == "" {
		host = "http://localhost:8080"
	}
	if !strings.Contains(host, "/v1") {
		host = strings.TrimRight(host, "/") + "/v1"
	}
	return host
}

func getAuthHeader(provider string) (string, string) {
	if config.ProviderType(provider) == "ollama" {
		return "", ""
	}
	key := config.ProviderAPIKey(provider)
	if key == "" {
		return "", ""
	}
	return "Authorization", "Bearer " + key
}

func Call(messages []types.Message, model, provider string) (string, error) {
	url := getBaseURL(provider) + "/chat/completions"

	payload := map[string]interface{}{
		"model":       model,
		"messages":    messages,
		"stream":      true,
		"temperature": 0.7,
		"top_p":       0.9,
	}

	data, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")

	if name, value := getAuthHeader(provider); name != "" {
		req.Header.Set(name, value)
	}

	stop := ui.StartSpinner()
	defer stop()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		stop()
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		stop()
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API %d: %s", resp.StatusCode, body)
	}

	var sb strings.Builder
	r := bufio.NewReader(resp.Body)
	first := true

	for {
		line, err := r.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				fmt.Fprintf(os.Stderr, "stream err: %v\n", err)
			}
			break
		}
		line = strings.TrimSpace(line)
		if line == "" || line == "data: [DONE]" || !strings.HasPrefix(line, "data: ") {
			continue
		}

		var chunk types.ChatCompletionChunk
		if json.Unmarshal([]byte(strings.TrimPrefix(line, "data: ")), &chunk) != nil || len(chunk.Choices) == 0 {
			continue
		}

		content := chunk.Choices[0].Delta.Content
		if content == "" {
			continue
		}

		if first {
			stop()
			fmt.Printf("%s%s\n %s↳ %s",
				ui.RedBold(model),
				ui.Reset(),
				ui.Cyan("↳ "),
				ui.Reset(),
			)
			first = false
		}

		fmt.Print(content)
		sb.WriteString(content)
	}

	fmt.Println()

	if first {
		stop()
		return "", fmt.Errorf("empty response")
	}

	return sb.String(), nil
}

func ListModels(provider string) ([]string, error) {
	url := getBaseURL(provider) + "/models"

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("status %d", resp.StatusCode)
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
