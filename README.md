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
```

```text
sudo cp aicli /usr/local/bin/
```

```text
mkdir -p path/to/project
aicli --version
aicli --init
```

```text
nano .aicli/config.yaml
```

```text
aicli
```

### TODO
------

- READY
  - Подключается к LocalAI и Ollama
  - Смотри /help

- TODO
  - Сохранение провайдера
  - Автовыбор модели если не указана
  - Добавить pwd в /status
  - ⊂(◉‿◉)つ
