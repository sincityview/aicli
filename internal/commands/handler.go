package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

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

	// ====================== СПЕЦИАЛЬНЫЕ КОМАНДЫ ======================
	switch {
	case input == "/exit" || input == "exit":
		return false

	case input == "/help":
		fmt.Print(`Доступные команды:
  /status        — информация о подключении
  /models        — список доступных моделей
  /model <name>  — сменить модель
  /list          — список сохранённых чатов
  /new           — новый чат
  /open <name>   — открыть чат
  /delete <name> — удалить чат
  /rename <old> <new> — переименовать чат
  /clear         — очистить текущий чат
  /read <file>   — добавить файл в контекст
  /save <file>   — сохранить последний ответ
  /help          — эта справка
  /exit          — выход
`)
		return true

	case input == "/status":
		fmt.Printf("%s=== Статус подключения ===%s\n", ui.ColorCyan, ui.ColorReset)
		fmt.Printf("Host           : %s\n", config.Host())
		fmt.Printf("Текущая модель : %s\n", s.CurrentModel)
		fmt.Printf("Текущий чат    : %s\n", s.CurrentChat)
		fmt.Printf("Сообщений      : %d\n", len(s.History))

		if models, err := client.ListModels(); err == nil {
			fmt.Printf("Статус сервера : %sПодключено%s (%d моделей)\n", ui.ColorGreen, ui.ColorReset, len(models))
		} else {
			fmt.Printf("Статус сервера : %sОшибка%s\n  → %v\n", ui.ColorRed, ui.ColorReset, err)
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
			mark := " "
			if m == s.CurrentModel {
				mark = ui.ColorGreen + "→" + ui.ColorReset + " "
			}
			fmt.Printf("%s• %s\n", mark, m)
		}
		return true

	case strings.HasPrefix(input, "/model "):
		newModel := strings.TrimSpace(strings.TrimPrefix(input, "/model "))
		if newModel == "" {
			fmt.Println("Использование: /model <имя_модели>")
			return true
		}
		s.CurrentModel = newModel
		fmt.Printf("%sМодель изменена на: %s%s\n", ui.ColorGreen, s.CurrentModel, ui.ColorReset)
		return true

	case input == "/list":
		names, _ := chat.List()
		fmt.Printf("\n%sСохранённые чаты:%s\n", ui.ColorCyan, ui.ColorReset)
		for _, n := range names {
			fmt.Printf(" • %s\n", n)
		}
		return true

	case input == "/new":
		s.CurrentChat = "chat-" + time.Now().Format("2006-01-02_15-04-05")
		s.History = []types.Message{{Role: "system", Content: config.SystemPrompt()}}
		fmt.Printf("%sНовая сессия: %s%s\n", ui.ColorCyan, s.CurrentChat, ui.ColorReset)
		return true

	case strings.HasPrefix(input, "/open "):
		name := strings.TrimSpace(strings.TrimPrefix(input, "/open "))
		if h, err := chat.Load(name); err == nil {
			s.History = h
			s.CurrentChat = name
			fmt.Printf("%sЧат загружен: %s%s\n", ui.ColorCyan, name, ui.ColorReset)
		} else {
			fmt.Printf("%sЧат '%s' не найден%s\n", ui.ColorRed, name, ui.ColorReset)
		}
		return true

	case input == "/clear":
		s.History = []types.Message{{Role: "system", Content: config.SystemPrompt()}}
		_ = chat.Save(s.History, s.CurrentChat)
		fmt.Println("История текущего чата очищена.")
		return true

	case strings.HasPrefix(input, "/read "):
		path := strings.TrimSpace(strings.TrimPrefix(input, "/read "))
		if data, err := os.ReadFile(path); err == nil {
			input = fmt.Sprintf("Содержимое файла %s:\n%s", filepath.Base(path), string(data))
			s.History = append(s.History, types.Message{Role: "user", Content: input})
			// сразу отправляем как сообщение
			goto sendMessage
		} else {
			fmt.Printf("%sОшибка чтения файла: %v%s\n", ui.ColorRed, err, ui.ColorReset)
			return true
		}

	case strings.HasPrefix(input, "/save "):
		if s.LastResponse == "" {
			fmt.Println("Нет последнего ответа для сохранения.")
			return true
		}
		path := strings.TrimSpace(strings.TrimPrefix(input, "/save "))
		if err := os.WriteFile(path, []byte(s.LastResponse), 0644); err == nil {
			fmt.Printf("%sСохранено в %s%s\n", ui.ColorGreen, path, ui.ColorReset)
		} else {
			fmt.Printf("%sОшибка сохранения: %v%s\n", ui.ColorRed, err, ui.ColorReset)
		}
		return true

	case strings.HasPrefix(input, "/delete ") || strings.HasPrefix(input, "/rename "):
		// Можно расширить позже
		fmt.Println("Команда пока не реализована полностью.")
		return true

	default:
		// Обычное сообщение пользователя
		s.History = append(s.History, types.Message{Role: "user", Content: input})
	}

sendMessage:
	// Отправка сообщения модели
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
