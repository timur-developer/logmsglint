package fix

import "log/slog"

func f() {
	slog.Info("Starting server") // want "log message must start with a lowercase letter"
}
