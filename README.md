# logmsglint

`logmsglint` - кастомный линтер для Go, который проверяет сообщения логов в `log/slog` и `go.uber.org/zap`.

Линтер помогает поддерживать единый стиль лог-сообщений и снижает риск случайного попадания чувствительных данных в логи. Его можно запускать как standalone-анализатор или подключать к `golangci-lint` v2 через Module Plugin System.

## Что проверяет

- сообщение лога должно начинаться со строчной буквы;
- сообщение должно быть на английском языке;
- сообщение не должно содержать специальные символы или emoji;
- сообщение не должно содержать потенциально чувствительные данные по ключевым словам и регулярным выражениям.

Для правила про строчную букву поддерживается autofix: линтер может автоматически заменить первую заглавную букву на строчную.

## Требования

- Go 1.22+
- git
- golangci-lint v2, проверялось на `v2.10.1`

## Установка

Склонируйте репозиторий:

```bash
git clone https://github.com/timur-developer/logmsglint.git
cd logmsglint
```

При необходимости скачайте зависимости:

```bash
go mod download
```

## Сборка

Сборка standalone-анализатора:

```bash
go build -o logmsglint ./cmd/logmsglint
```

На Windows:

```bash
go build -o logmsglint.exe ./cmd/logmsglint
```

## Запуск как standalone-анализатор

Проверить текущий проект:

```bash
go run ./cmd/logmsglint ./...
```

Проверить отдельный пакет:

```bash
go run ./cmd/logmsglint ./scratch
```

## Запуск через golangci-lint

Проект использует Module Plugin System для подключения линтера к `golangci-lint`.

Сначала соберите кастомный бинарник `golangci-lint` с плагином:

```bash
golangci-lint custom -v
```

Команда использует конфигурацию из `.custom-gcl.yml` и создает бинарник `custom-gcl` в корне проекта. На Windows будет создан `custom-gcl.exe`.

Запустите проверку через собранный бинарник:

```bash
./custom-gcl run -c .golangci.yml ./...
```

Проверить отдельный пакет:

```bash
./custom-gcl run -c .golangci.yml ./scratch
```

Запуск с autofix:

```bash
./custom-gcl run -c .golangci.yml --fix ./...
```

Autofix применяется к правилу про первую букву сообщения. Например:

```go
slog.Info("Starting server")
```

будет заменено на:

```go
slog.Info("starting server")
```

## Конфигурация

Линтер настраивается в `.golangci.yml` в секции:

```text
linters.settings.custom.logmsglint.settings
```

Минимальный пример конфигурации:

```yaml
version: "2"

linters:
  default: none
  enable:
    - logmsglint

  settings:
    custom:
      logmsglint:
        type: module
        description: Checks slog/go.uber.org/zap log messages to basic style and safety rules.
        settings:
          enable_fixes: true
          allowed_punct: ""
          sensitive_keywords:
            - "secret_token"
          sensitive_regexps:
            - "(?i)api[_-]?key\\s*="
```

Параметры:

- `enable_fixes` - включает `SuggestedFix` для правила про строчную букву;
- `allowed_punct` - список разрешенных знаков пунктуации в сообщениях;
- `sensitive_keywords` - дополнительные ключевые слова, которые считаются чувствительными;
- `sensitive_regexps` - пользовательские регулярные выражения для поиска чувствительных данных.

## Примеры

Папка `scratch/` не хранится в репозитории. Чтобы повторить демо-команды, создайте локальный файл `scratch/main.go`.

Пример файла:

```go
package scratch

import "log/slog"

func main() {
	slog.Info("Starting server on port 8080")
	slog.Info("запуск сервера!")
	slog.Info("server started!🚀")
	slog.Info("token: " + "abc")
}
```

Запуск standalone-анализатора:

```bash
go run ./cmd/logmsglint ./scratch
```

Ожидаемый результат:

```text
scratch/main.go:...: log message must start with a lowercase letter
scratch/main.go:...: log message must contain only English letters; log message must not contain special characters or emoji
scratch/main.go:...: log message must not contain special characters or emoji
scratch/main.go:...: log message must not contain special characters or emoji; log message may contain sensitive data
```

Пример запуска autofix:

```bash
./custom-gcl run -c .golangci.yml --fix ./...
```

Он заменит:

```go
slog.Info("Starting server on port 8080")
```

на:

```go
slog.Info("starting server on port 8080")
```

## CI

В GitHub Actions запускаются unit-тесты, сборка плагина и проверка линтера на каждый push и pull request.

## Для чего полезен

`logmsglint` полезен в проектах, где важно держать логи в едином формате и заранее отсекать потенциально опасные сообщения. Он особенно хорошо подходит для командных Go-проектов с общим стандартом логирования через `slog` или `zap`.
