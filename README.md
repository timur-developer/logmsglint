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

Apply autofixes (lowercase)
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