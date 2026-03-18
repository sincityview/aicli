package types

// Message — стандартное сообщение OpenAI формата
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatCompletionChunk — структура для streaming ответа
type ChatCompletionChunk struct {
	Choices []struct {
		Delta struct {
			Content string `json:"content"`
		} `json:"delta"`
	} `json:"choices"`
}

// ModelList — ответ от /v1/models
type ModelList struct {
	Data []struct {
		ID string `json:"id"`
	} `json:"data"`
}
