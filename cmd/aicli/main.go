package main

import (
	"fmt"
	"os"
	"strings"

	"aicli/internal/chat"
	"aicli/internal/client"
	"aicli/internal/commands"
	"aicli/internal/config"
	"aicli/internal/types"
	"aicli/internal/ui"

	"github.com/chzyer/readline"
)

const appVersion = "v0.2"

func initWorkdir() {
	_ = os.MkdirAll(".aicli/sessions", 0755)
	cfg := []byte(`current_provider: localai

providers:
  localai:
    type: localai
    host: http://localhost:8080
    api_key: ""
    default_model: qwen_qwen3.5-2b

  ollama:
    type: ollama
    host: http://127.0.0.1:11434
    api_key: ""
    default_model: qwen2.5:3b
`)
	_ = os.WriteFile(".aicli/config.yaml", cfg, 0644)
	fmt.Println("Проект инициализирован в .aicli/")
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "--init" {
		initWorkdir()
		return
	}
	if len(os.Args) > 1 && os.Args[1] == "--version" {
		fmt.Printf("aicli %s. CLI for interacting with LocalAI and Ollama.\n", appVersion)
		return
	}

	config.Init()

	// One-shot режим
	if len(os.Args) > 1 {
		prompt := strings.Join(os.Args[1:], " ")
		history := []types.Message{
			{Role: "system", Content: config.SystemPrompt()},
			{Role: "user", Content: prompt},
		}
		resp, err := client.Call(history, config.ProviderDefaultModel(config.CurrentProvider()), config.CurrentProvider())
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s%s%s\n", ui.Red(""), err, ui.Reset())
			os.Exit(1)
		}
		fmt.Println(resp)
		return
	}

	// Интерактивный режим
	state := commands.NewState()

	// Загружаем историю текущего чата
	if h, err := chat.Load(state.CurrentChat); err == nil {
		state.History = h
	}

	rl, err := readline.NewEx(&readline.Config{
		Prompt:          fmt.Sprintf(" %s↳ %s", ui.Cyan(""), ui.Reset()),
		HistoryFile:     ".aicli/history.tmp",
		AutoComplete:    ui.Completer{},
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "readline init failed: %v\n", err)
		os.Exit(1)
	}
	defer rl.Close()

	fmt.Printf("%s--- aicli %s (Provider: %s | Model: %s | Chat: %s) ---%s\n\n",
		ui.Cyan(""), appVersion, state.CurrentProvider, state.CurrentModel, state.CurrentChat, ui.Reset())

	for {
		fmt.Printf("\n%s%s%s\n", ui.Green("")+ui.Bold(""), os.Getenv("USER"), ui.Reset())

		line, err := rl.Readline()
		if err != nil {
			break
		}

		if !commands.Handle(state, line) {
			break
		}
	}
}
