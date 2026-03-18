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

const appVerion = "v0.1"

func initWorkdir() {
	_ = os.MkdirAll(".aicli/sessions", 0755)
	cfg := []byte(`host: "http://localhost:8080"
model: "qwen_qwen3.5-2b"
system_prompt: "You are a helpful assistant."
api_key: ""
`)
	_ = os.WriteFile(".aicli/config.yaml", cfg, 0644)
	fmt.Println("Проект инициализирован в .aicli/")
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "--init" {
		initWorkdir()
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
		resp, err := client.Call(history, config.Model())
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s%s%s\n", ui.ColorRed, err, ui.ColorReset)
			os.Exit(1)
		}
		fmt.Println(resp)
		return
	}

	// Интерактивный режим
	state := commands.NewState()
	// Загружаем историю чата (при ошибке будет использован новый чат с system prompt)
	if h, err := chat.Load(state.CurrentChat); err == nil {
		state.History = h
	}

	rl, err := readline.NewEx(&readline.Config{
		Prompt:          fmt.Sprintf(" %s↳ %s", ui.ColorCyan, ui.ColorReset),
		HistoryFile:     ".aicli/history.tmp",
		AutoComplete:    ui.Completer{},
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to init readline: %v\n", err)
		os.Exit(1)
	}
	defer rl.Close()

	fmt.Printf("%s--- aicli %s (Model: %s | Chat: %s) ---%s\n\n",
		ui.ColorCyan, appVerion, state.CurrentModel, state.CurrentChat, ui.ColorReset)

	for {
		fmt.Printf("\n%s%s%s\n", ui.ColorGreen+ui.ColorBold, os.Getenv("USER"), ui.ColorReset)

		line, err := rl.Readline()
		if err != nil {
			break
		}

		if !commands.Handle(state, line) {
			break
		}
	}
}
