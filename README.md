# logmsglint — linter for Go log messages in slog and zap

![logmsglint-logo](https://raw.githubusercontent.com/timur-developer/logmsglint/refs/heads/main/logmsglint-logo.png)

[](https://github.com/timur-developer/logmsglint/actions/workflows/ci.yml)
![Go](https://img.shields.io/badge/go-1.22%2B-00ADD8?logo=go&logoColor=white)
![golangcilint](https://img.shields.io/badge/golangci--lint-v2-181717?logo=go)
![License MIT](https://img.shields.io/badge/license-MIT-blue.svg)

`logmsglint` is a Go static analyzer for validating log messages in `log/slog` and `go.uber.org/zap` calls.

It helps keep log messages consistent before code review: lowercase style, English-only messages, no emoji or unexpected special characters, and configurable checks for sensitive data patterns.

The analyzer can be used as a standalone tool or integrated into `golangci-lint` through the custom module plugin system.

## Contents

- [Why](#why)
- [What It Checks](#what-it-checks)
- [Supported Loggers](#supported-loggers)
- [Install](#install)
- [Usage](#usage)
  - [Standalone](#standalone)
  - [With golangci-lint](#with-golangci-lint)
- [Configuration](#configuration)
- [Examples](#examples)
- [Autofix](#autofix)
- [CI](#ci)
- [Development](#development)
- [License](#license)

## Why

Logs are part of the public interface of a backend service. They are used during debugging, incident response, monitoring, and support.

Without automated checks, log messages often become inconsistent over time:

- some start with uppercase letters, others with lowercase letters
- some use Russian or mixed-language text
- some contain emojis or punctuation that makes logs harder to search
- some accidentally include sensitive values or suspicious key names

`logmsglint` catches these issues statically, before they reach code review or production.

## What It Checks

| Rule | Description | Autofix |
| --- | --- | --- |
| Lowercase start | Log message should start with a lowercase letter | yes |
| English-only message | Log message should use English text | no  |
| No special characters / emoji | Disallows unexpected punctuation and emoji | no  |
| Sensitive data patterns | Detects configured keywords and regular expressions | no  |

Sensitive data checks are configurable, so teams can tune them for their own conventions.

## Supported Loggers

`logmsglint` checks message arguments in common Go logging calls.

Supported packages:

- `log/slog`
- `go.uber.org/zap`

Example calls:

```go
slog.Info("server started")
slog.Error("request failed", "err", err)

logger.Info("server started")
logger.Error("request failed", zap.Error(err))
```

## Install

### Standalone binary

Clone and build:

```bash
git clone https://github.com/timur-developer/logmsglint.git
cd logmsglint
go build -o logmsglint ./cmd/logmsglint
```

Run:

```bash
./logmsglint ./...
```

### golangci-lint plugin

`logmsglint` can also be embedded into a custom `golangci-lint` binary.

Build the custom linter binary:

```bash
golangci-lint custom -v
```

Then run it:

```bash
./custom-gcl run ./...
```

## Usage

### Standalone

Run the analyzer for the current module:

```bash
./logmsglint ./...
```

Run for a specific package:

```bash
./logmsglint ./internal/service
```

### With golangci-lint

Create a `.golangci.yml` configuration and enable `logmsglint`.

```bash
./custom-gcl run -c .golangci.yml ./...
```

Run with fixes enabled:

```bash
./custom-gcl run -c .golangci.yml --fix ./...
```

## Configuration

Example `.golangci.yml`:

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

Configuration ideas:

- use `allowed_punct` if your team allows selected punctuation in log messages
- add project-specific words to `sensitive_keywords`
- use `sensitive_regexps` for patterns such as `token=...`, `api_key=...`, or similar forms
- enable `enable_fixes` if you want automatic lowercase fixes

## Examples

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

After running with `--fix`, the first message becomes:

```go
slog.Info("starting server")
```

## Autofix

Autofix is intentionally conservative.

Currently, `logmsglint` can fix messages that only violate the lowercase-start rule:

```go
slog.Info("Starting server")
```

becomes:

```go
slog.Info("starting server")
```

The analyzer does not automatically rewrite non-English messages, remove emoji, or edit potentially sensitive content because those changes require human judgment.

## CI

Example GitHub Actions step for standalone usage:

```yaml
- name: Run logmsglint
  run: go run ./cmd/logmsglint ./...
```

Example step for a custom `golangci-lint` binary:

```yaml
- name: Build custom golangci-lint
  run: golangci-lint custom -v

- name: Run linters
  run: ./custom-gcl run -c .golangci.yml ./...
```

## Development

Run tests:

```bash
go test ./...
```

Build:

```bash
go build ./...
```

Run the analyzer from source:

```bash
go run ./cmd/logmsglint ./...
```


## License

MIT. See [LICENSE](LICENSE).
