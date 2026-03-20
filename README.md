#### CLI for interacting with LocalAI / Ollama
------

**Версия App**
```text
v0.2
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

**Создание локального проекта**
```text
mkdir -p path/to/project
aicli --version
aicli --init
```

**Настройки провайдера и модели**
```text
nano .aicli/config.yaml
```

**Запуск**
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
