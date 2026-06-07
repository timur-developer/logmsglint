# logmsglint — линтер Go-сообщений логирования в slog и zap

![logmsglintlogo](https://raw.githubusercontent.com/timur-developer/logmsglint/refs/heads/main/logmsglint-logo.png)

[![CI](https://github.com/timur-developer/logmsglint/actions/workflows/ci.yml/badge.svg)](https://github.com/timur-developer/logmsglint/actions/workflows/ci.yml)
![Go](https://img.shields.io/badge/go-1.22%2B-00ADD8?logo=go&logoColor=white)
![golangcilint](https://img.shields.io/badge/golangci--lint-v2-181717?logo=go)
![License MIT](https://img.shields.io/badge/license-MIT-blue.svg)

Read this in other languages: [English](README.md)

`logmsglint` — статический анализатор для Go, который проверяет сообщения логирования в вызовах `log/slog` и `go.uber.org/zap`.

Он помогает поддерживать единый стиль логов ещё до code review: сообщения с маленькой буквы, английский текст, отсутствие emoji и неожиданных специальных символов, а также настраиваемые проверки на потенциально чувствительные данные.

Анализатор можно использовать как самостоятельный инструмент или встроить в `golangci-lint` через custom module plugin system.

## Содержание

- [Зачем](#зачем)
- [Что проверяет](#что-проверяет)
- [Поддерживаемые логгеры](#поддерживаемые-логгеры)
- [Установка](#установка)
- [Использование](#использование)
  - [Standalone](#standalone)
  - [С golangci-lint](#с-golangci-lint)
- [Конфигурация](#конфигурация)
- [Примеры](#примеры)
- [Autofix](#autofix)
- [CI](#ci)
- [Разработка](#разработка)
- [Лицензия](#лицензия)

## Зачем

Логи — часть публичного интерфейса backend-сервиса. Их используют при debugging, incident response, мониторинге и поддержке.

Без автоматических проверок сообщения логирования со временем становятся неоднородными:

- одни начинаются с заглавной буквы, другие — с маленькой
- часть сообщений написана на русском или смешивает несколько языков
- где-то появляются emoji или пунктуация, из-за которых сложнее искать по логам
- в сообщениях случайно оказываются чувствительные значения или подозрительные имена ключей

`logmsglint` ловит такие проблемы статически — до code review и до попадания кода в production.

## Что проверяет

| Правило | Описание | Autofix |
| --- | --- | --- |
| Начало с маленькой буквы | Сообщение лога должно начинаться с маленькой буквы | да  |
| Только английский текст | Сообщение лога должно быть написано на английском | нет |
| Без специальных символов и emoji | Запрещает неожиданные знаки пунктуации и emoji | нет |
| Паттерны чувствительных данных | Находит настроенные keywords и регулярные выражения | нет |

Проверки чувствительных данных настраиваются, поэтому команду можно адаптировать под собственные правила и соглашения.

## Поддерживаемые логгеры

`logmsglint` проверяет аргументы сообщений в распространённых Go logging calls.

Поддерживаемые пакеты:

- `log/slog`
- `go.uber.org/zap`

Примеры вызовов:

```go
slog.Info("server started")
slog.Error("request failed", "err", err)

logger.Info("server started")
logger.Error("request failed", zap.Error(err))
```

## Установка

### Standalone binary

Склонируйте репозиторий и соберите бинарник:

```bash
git clone https://github.com/timur-developer/logmsglint.git
cd logmsglint
go build -o logmsglint ./cmd/logmsglint
```

Запуск:

```bash
./logmsglint ./...
```

### golangci-lint plugin

`logmsglint` можно встроить в custom `golangci-lint` binary.

Соберите custom linter binary:

```bash
golangci-lint custom -v
```

Затем запустите его:

```bash
./custom-gcl run ./...
```

## Использование

### Standalone

Запустить анализатор для текущего Go module:

```bash
./logmsglint ./...
```

Запустить для конкретного пакета:

```bash
./logmsglint ./internal/service
```

### С golangci-lint

Создайте `.golangci.yml` и включите `logmsglint`.

```bash
./custom-gcl run -c .golangci.yml ./...
```

Запуск с автоисправлениями:

```bash
./custom-gcl run -c .golangci.yml --fix ./...
```

## Конфигурация

Пример `.golangci.yml`:

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
        description: Checks slog/zap log messages.
        settings:
          enable_fixes: true
          allowed_punct: ""
          sensitive_keywords:
            - "password"
            - "secret"
            - "secret_token"
          sensitive_regexps:
            - "(?i)api[_-]?key\\s*="
            - "(?i)token\\s*="
```

Идеи для настройки:

- используйте `allowed_punct`, если в вашей команде разрешена определённая пунктуация в log messages
- добавляйте project-specific слова в `sensitive_keywords`
- используйте `sensitive_regexps` для паттернов вроде `token=...`, `api_key=...` и похожих форм
- включайте `enable_fixes`, если хотите автоматически исправлять начало сообщения на маленькую букву

## Примеры

Input:

```go
package main

import "log/slog"

func main() {
    slog.Info("Starting server")
    slog.Info("запуск сервера")
    slog.Info("server started 🚀")
    slog.Info("token=abc123")
}
```

Output:

```text
main.go:6: log message must start with a lowercase letter
main.go:7: log message must contain only English letters
main.go:8: log message must not contain emoji or disallowed special characters
main.go:9: log message may contain sensitive data
```

После запуска с `--fix` первое сообщение станет таким:

```go
slog.Info("starting server")
```

## Autofix

Autofix намеренно сделан консервативным.

Сейчас `logmsglint` умеет автоматически исправлять только сообщения, которые нарушают правило lowercase-start:

```go
slog.Info("Starting server")
```

становится:

```go
slog.Info("starting server")
```

Анализатор не переписывает неанглийские сообщения, не удаляет emoji и не редактирует потенциально чувствительный контент автоматически, потому что такие изменения требуют человеческого решения.

## CI

Пример шага GitHub Actions для standalone usage:

```yaml
- name: Run logmsglint
  run: go run ./cmd/logmsglint ./...
```

Пример для custom `golangci-lint` binary:

```yaml
- name: Build custom golangci-lint
  run: golangci-lint custom -v

- name: Run linters
  run: ./custom-gcl run -c .golangci.yml ./...
```

## Разработка

Запустить тесты:

```bash
go test ./...
```

Собрать проект:

```bash
go build ./...
```

Запустить анализатор из исходников:

```bash
go run ./cmd/logmsglint ./...
```

## Лицензия

MIT. См. [LICENSE](LICENSE).
