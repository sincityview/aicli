package commands

import (
	"fmt"
	"strings"

	"aicli/internal/chat"
	"aicli/internal/client"
	"aicli/internal/config"
	"aicli/internal/types"
	"aicli/internal/ui"
)

type State struct {
	History      []types.Message
	CurrentModel string
	CurrentChat  string
	LastResponse string
}

func NewState() *State {
	return &State{
		CurrentModel: config.Model(),
		CurrentChat:  "default",
		History:      []types.Message{{Role: "system", Content: config.SystemPrompt()}},
	}
}

func Handle(s *State, input string) bool {
	input = strings.TrimSpace(input)
	if input == "" {
		return true
	}

	switch {
	case input == "/exit" || input == "exit":
		return false

	case input == "/help":
		fmt.Print(`
  Доступные команды:
  /status        		- информация о подключении и текущем состоянии
  /models        		- список доступных моделей
  /list          		- список сохранённых чатов
  /new           		- начать новый чат
  /open <name>   		- открыть существующий чат
  /delete <name>		- удалить чат
  /rename <old> <new>		- переименовать чат
  /model <name>  		- сменить модель
  /clear         		- очистить историю текущего чата
  /read <file>   		- добавить содержимое файла в контекст
  /save <file>   		- сохранить последний ответ в файл
  /help          		- показать эту справку
  /exit          		- выйти
`)
		return true

	case input == "/status":
		fmt.Printf("%s=== Статус подключения ===%s\n", ui.ColorCyan, ui.ColorReset)
		fmt.Printf("Host           : %s\n", config.Host())
		fmt.Printf("Текущая модель : %s\n", s.CurrentModel)
		fmt.Printf("Текущий чат    : %s\n", s.CurrentChat)
		fmt.Printf("Сообщений      : %d\n", len(s.History))

		// Проверка соединения
		if models, err := client.ListModels(); err == nil {
			fmt.Printf("Статус сервера : %sПодключено%s (%d моделей)\n", ui.ColorGreen, ui.ColorReset, len(models))
		} else {
			fmt.Printf("Статус сервера : %sОшибка подключения%s\n", ui.ColorRed, ui.ColorReset)
			fmt.Printf("  → %v\n", err)
		}
		return true

	case input == "/models":
		fmt.Printf("%sПолучение списка моделей...%s\n", ui.ColorCyan, ui.ColorReset)
		models, err := client.ListModels()
		if err != nil {
			fmt.Printf("%sОшибка: %v%s\n", ui.ColorRed, err, ui.ColorReset)
			return true
		}
		fmt.Printf("\n%sДоступные модели (%d):%s\n", ui.ColorCyan, len(models), ui.ColorReset)
		for _, m := range models {
			if m == s.CurrentModel {
				fmt.Printf(" • %s%s (текущая)%s\n", ui.ColorGreen, m, ui.ColorReset)
			} else {
				fmt.Printf(" • %s\n", m)
			}
		}
		return true

	case input == "/list":
		names, _ := chat.List()
		fmt.Printf("\n%sСохранённые чаты:%s\n", ui.ColorCyan, ui.ColorReset)
		for _, n := range names {
			fmt.Printf(" • %s\n", n)
		}
		return true

	default:
		// Обработка обычного сообщения пользователя (как было раньше)
		s.History = append(s.History, types.Message{Role: "user", Content: input})

		response, err := client.Call(s.History, s.CurrentModel)
		if err != nil {
			fmt.Printf("\n%sОшибка: %v%s\n", ui.ColorRed, err, ui.ColorReset)
			return true
		}

		s.LastResponse = response
		s.History = append(s.History, types.Message{Role: "assistant", Content: response})
		_ = chat.Save(s.History, s.CurrentChat)
		return true
	}
}
