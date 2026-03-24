package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"aicli/internal/chat"
	"aicli/internal/client"
	"aicli/internal/config"
	"aicli/internal/types"
	"aicli/internal/ui"
)

type State struct {
	History         []types.Message
	CurrentModel    string
	CurrentChat     string
	LastResponse    string
	CurrentProvider string
}

var CurrentState *State

func NewState() *State {
	prov := config.CurrentProvider()
	if prov == "" {
		prov = "localai-1"
	}
	model := config.ProviderDefaultModel(prov)
	if model == "" {
		model = "qwen_qwen3.5-2b"
	}

	s := &State{
		CurrentProvider: prov,
		CurrentModel:    model,
		CurrentChat:     "default",
		History:         []types.Message{{Role: "system", Content: config.SystemPrompt()}},
	}
	CurrentState = s
	return s
}

func Handle(s *State, input string) bool {
	input = strings.TrimSpace(input)
	if input == "" {
		return true
	}

	parts := strings.Fields(input)
	cmd := parts[0]
	args := parts[1:]

	switch cmd {
	case "/exit", "exit", "/quit", "quit":
		return false

	case "/help":
		fmt.Println(`Команды:
  /help               эта справка
  /status             текущее состояние
  /model list         список моделей
  /model use <имя>    выбрать модель
  /provider list      список провайдеров
  /provider use <имя> выбрать провайдер
  /chat list          список чатов
  /chat use <имя>     открыть чат или создать если не существует
  /clear              очистить текущий чат
  /read <файл>        добавить содержимое файла
  /save <файл>        сохранить последний ответ
  /quit               выход из программы`)
		return true

	case "/status":
		ctxLen := config.ProviderContextLength(s.CurrentProvider)

		fmt.Printf("%sСостояние:%s\n", ui.Cyan(""), ui.Reset())
		fmt.Printf("  Провайдер   : %s%s%s\n", ui.Green(s.CurrentProvider), " ", ui.Reset())
		fmt.Printf("  Модель      : %s%s%s\n", ui.Green(s.CurrentModel), " ", ui.Reset())
		fmt.Printf("  Чат         : %s%s%s\n", ui.Green(s.CurrentChat), " ", ui.Reset())
		fmt.Printf("  Сообщений   : %s%d%s\n", ui.Green(""), len(s.History), ui.Reset())
		fmt.Printf("  Контекст    : %s%d токенов%s\n", ui.Green(""), ctxLen, ui.Reset())
		return true

	case "/model":
		if len(args) == 0 || args[0] == "help" {
			fmt.Println("Использование:")
			fmt.Println("  /model list                — список моделей")
			fmt.Println("  /model use <имя или номер> — выбрать модель")
			return true
		}
		sub := args[0]
		subArgs := args[1:]

		switch sub {
		case "list":
			models, err := client.ListModels(s.CurrentProvider)
			if err != nil {
				fmt.Printf("%sОшибка: %v%s\n", ui.Red(""), err, ui.Reset())
				return true
			}
			fmt.Printf("%sМодели (%d):%s\n", ui.Cyan(""), len(models), ui.Reset())
			for i, m := range models {
				prefix := "   "
				if m == s.CurrentModel {
					prefix = " • "
				}
				numColor := ""
				if m == s.CurrentModel {
					numColor = ui.Green("")
				}
				fmt.Printf("%s%s%d)%s %s\n", prefix, numColor, i+1, ui.Reset(), m)
			}
			return true

		case "use":
			if len(subArgs) == 0 {
				fmt.Println("Использование: /model use <имя или номер>")
				return true
			}
			arg := subArgs[0]
			if n, err := strconv.Atoi(arg); err == nil {
				models, _ := client.ListModels(s.CurrentProvider)
				if n > 0 && n <= len(models) {
					s.CurrentModel = models[n-1]
					fmt.Printf("%sМодель выбрана: %s%s\n", ui.Green(""), s.CurrentModel, ui.Reset())
					return true
				}
			}
			s.CurrentModel = arg
			fmt.Printf("%sМодель выбрана: %s%s\n", ui.Green(""), s.CurrentModel, ui.Reset())
			return true

		default:
			fmt.Println("Неизвестная подкоманда. Используйте /model list или /model use")
			return true
		}

	case "/provider":
		if len(args) == 0 || args[0] == "help" {
			fmt.Println("Использование:")
			fmt.Println("  /provider list       — список провайдеров")
			fmt.Println("  /provider use <имя>  — выбрать провайдер")
			return true
		}
		sub := args[0]
		subArgs := args[1:]

		switch sub {
		case "list":
			fmt.Printf("%sПровайдеры:%s\n", ui.Cyan(""), ui.Reset())
			provs := config.GetAllProviders()
			i := 1
			for p := range provs {
				prefix := "   "
				if p == s.CurrentProvider {
					prefix = " • "
				}
				numColor := ""
				if p == s.CurrentProvider {
					numColor = ui.Green("")
				}
				fmt.Printf("%s%s%d)%s %s\n", prefix, numColor, i, ui.Reset(), p)
				i++
			}
			return true

		case "use":
			if len(subArgs) == 0 {
				fmt.Println("Использование: /provider use <имя>")
				return true
			}
			name := subArgs[0]
			if _, ok := config.GetAllProviders()[name]; !ok {
				fmt.Printf("%sПровайдер не найден: %s%s\n", ui.Red(""), name, ui.Reset())
				return true
			}
			s.CurrentProvider = name
			s.CurrentModel = config.ProviderDefaultModel(name)
			fmt.Printf("%sПровайдер выбран: %s (модель: %s)%s\n",
				ui.Green(""), name, s.CurrentModel, ui.Reset())
			return true

		default:
			fmt.Println("Неизвестная подкоманда")
			return true
		}

	case "/chat":
		if len(args) == 0 {
			fmt.Println("Использование:")
			fmt.Println("  /chat list       — список чатов")
			fmt.Println("  /chat use <имя>  — открыть чат")
			fmt.Println("  /chat delete <имя> — удалить чат")
			return true
		}
		sub := args[0]
		subArgs := args[1:]

		switch sub {
		case "list":
			names, _ := chat.List()
			fmt.Printf("%sЧаты:%s\n", ui.Cyan(""), ui.Reset())
			for i, n := range names {
				prefix := "   " // три пробела без точки
				if n == s.CurrentChat {
					prefix = " • "
				}
				numColor := ""
				if n == s.CurrentChat {
					numColor = ui.Green("")
				}
				fmt.Printf("%s%s%d)%s %s\n", prefix, numColor, i+1, ui.Reset(), n)
			}
			return true

		case "use":
			if len(subArgs) == 0 {
				fmt.Println("Использование: /chat use <имя>")
				return true
			}
			name := subArgs[0]
			if h, err := chat.Load(name); err == nil {
				s.History = h
				s.CurrentChat = name
				fmt.Printf("%sЧат открыт: %s%s\n", ui.Cyan(""), name, ui.Reset())
			} else {
				fmt.Printf("%sЧат не найден: %s%s\n", ui.Red(""), name, ui.Reset())
			}
			return true

		case "delete":
			if len(subArgs) == 0 {
				fmt.Println("Использование: /chat delete <имя_чата>")
				return true
			}
			name := subArgs[0]

			if err := chat.Delete(name); err != nil {
				fmt.Printf("%sОшибка: %v%s\n", ui.Red(""), err, ui.Reset())
				return true
			}

			fmt.Printf("%sЧат '%s' удалён%s\n", ui.Green(""), name, ui.Reset())

			// Если удалили текущий чат — переключаемся на default
			if name == s.CurrentChat {
				s.CurrentChat = "default"
				if h, err := chat.Load("default"); err == nil {
					s.History = h
				} else {
					s.History = []types.Message{{Role: "system", Content: config.SystemPrompt()}}
				}
				fmt.Printf("%sПереключено на чат: default%s\n", ui.Cyan(""), ui.Reset())
			}

			return true

		default:
			fmt.Println("Неизвестная подкоманда для /chat")
			return true
		}

	case "/clear":
		s.History = []types.Message{{Role: "system", Content: config.SystemPrompt()}}
		_ = chat.Save(s.History, s.CurrentChat)
		fmt.Println("История очищена")
		return true

	case "/read":
		if len(args) == 0 {
			fmt.Println("Использование: /read <путь_к_файлу>")
			return true
		}
		path := args[0]
		data, err := os.ReadFile(path)
		if err != nil {
			fmt.Printf("%sОшибка: %v%s\n", ui.Red(""), err, ui.Reset())
			return true
		}
		content := fmt.Sprintf("Содержимое %s:\n%s", filepath.Base(path), string(data))
		s.History = append(s.History, types.Message{Role: "user", Content: content})
		goto sendMessage

	case "/save":
		if s.LastResponse == "" {
			fmt.Println("Нет ответа для сохранения")
			return true
		}
		if len(args) == 0 {
			fmt.Println("Использование: /save <путь_к_файлу>")
			return true
		}
		path := args[0]
		err := os.WriteFile(path, []byte(s.LastResponse), 0644)
		if err == nil {
			fmt.Printf("%sСохранено: %s%s\n", ui.Green(""), path, ui.Reset())
		} else {
			fmt.Printf("%sОшибка: %v%s\n", ui.Red(""), err, ui.Reset())
		}
		return true

	default:
		s.History = append(s.History, types.Message{Role: "user", Content: input})
	}

sendMessage:
	response, err := client.Call(s.History, s.CurrentModel, s.CurrentProvider)
	if err != nil {
		fmt.Printf("%sОшибка: %v%s\n", ui.Red(""), err, ui.Reset())
		return true
	}
	s.LastResponse = response
	s.History = append(s.History, types.Message{Role: "assistant", Content: response})
	_ = chat.Save(s.History, s.CurrentChat)
	return true
}
