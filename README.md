#### CLI for interacting with LocalAI / Ollama
------

**Версия App**
```text
v0.1
```

**Версия Go**
```text
1.25
```

**Компиляция, настройка и запуск**
```text
go init aicli
go mod tidy
go build -o aicli cmd/aicli/main.go

sudo cp aicli /usr/local/bin/

mkdir -p path/to/project
aicli --init

nano .aicli/config.yaml

aicli
```

- READY
  - Подключается к LocalAI
  - Смотри /help

- TODO
  - ⊂(◉‿◉)つ
