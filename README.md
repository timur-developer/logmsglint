# logmsglint
Custom golangci-lint linter (module plugin) that checks slog and go.uber.zap log messages to basic style and safety rules.

## Rules
 - Log messages must start with a lowercase letter.

 - Log messages must be English-only (letters).

 - Log messages must not contain special characters / emoji.

 - Log messages must not contain potentially sensitive data (keyword-based).

## Requirements
 - Go 1.22+
 - git
 - golangci-lint v2 (tested with v2.10.1)

## Build

Build standalone analyzer binary:
```bash
go build -o logmsglint ./cmd/logmsglint
```

On Windows you can do:
```bash
go build -o logmsglint.exe ./cmd/logmsglint
```

## Run as standalone analyzer
Run on a package / project:

```bash
go run ./cmd/logmsglint ./...
```

Example (any package with Go files, e.g. `./scratch`):
```bash
go run ./cmd/logmsglint ./scratch
```

## Run with golangci-lint (module plugin)
This repository uses the Module Plugin System.

1. Build a custom golangci-lint binary with the plugin (uses .custom-gcl.yml):

```bash
golangci-lint custom -v
```

This produces custom-gcl (custom-gcl.exe on Windows) in the project root by default.

2. Run the custom binary with this repository config:

```bash
./custom-gcl run -c .golangci.yml ./...
```

Example (any package with Go files, e.g. `./scratch`):

```bash
./custom-gcl run -c .golangci.yml ./scratch
```

To apply autofixes use `--fix` flag (will replace uppercase to lowercase)
```bash
./custom-gcl run -c .golangci.yml --fix ./...
```

## Configuration
The plugin is configured via `.golangci.yml` under:

`linters.settings.custom.logmsglint.settings`.

Minimal configuration example:

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
        description: Checks slog/go.uber.zap log messages to basic style and safety rules.
        settings:
          enable_fixes: true             # enable SuggestedFix for lowercase rule
          allowed_punct: ""              # allowed punctuation characters in log messages
          sensitive_keywords:            # extra keywords treated as sensitive
            - "secret_token"
          sensitive_regexps:             # custom regex patterns for sensitive data
            - "(?i)api[_-]?key\\s*="
```

## Examples

Example code:

File: scratch/main.go

```bash
slog.Info("Starting server on port 8080")
slog.Info("запуск сервера!")
slog.Info("server started!🚀")
slog.Info("token: " + "abc")
```

Expected output (standalone analyzer):

`go run ./cmd/logmsglint ./scratch`

```bash
scratch/main.go:...: log message must start with a lowercase letter
scratch/main.go:...: log message must contain only English letters; log message must not contain special characters or emoji
scratch/main.go:...: log message must not contain special characters or emoji
scratch/main.go:...: log message must not contain special characters or emoji; log message may contain sensitive data
```

Autofix example:

`./custom-gcl run -c .golangci.yml --fix ./...`

It will rewrite:
`"Starting server on port 8080"`
to:
`"starting server on port 8080"`

