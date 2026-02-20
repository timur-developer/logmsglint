package lintcases

import (
	"context"
	"log/slog"

	"go.uber.org/zap"
)

func slogLowercase() {
	slog.Info("Starting server") // want "log message must start with a lowercase letter"
}

func slogEnglishOnly() {
	slog.Info("запуск сервера") // want "log message must contain only English letters"
}

func slogSpecials() {
	slog.Info("server started!") // want "log message must not contain special characters or emoji"
}

func slogSensitiveStatic() {
	slog.Info("token abc") // want "log message may contain sensitive data"
}

func slogSensitiveDynamic() {
	slog.Info("token " + token()) // want "log message may contain sensitive data"
}

func slogContextMsgIndex() {
	slog.InfoContext(context.Background(), "Hello") // want "log message must start with a lowercase letter"
}

func zapLowercase() {
	l, _ := zap.NewProduction()
	l.Info("Failed to connect") // want "log message must start with a lowercase letter"
}

func token() string { return "x" }
