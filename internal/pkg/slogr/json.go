package slogr

import (
	"log/slog"
	"os"
)

// SetDefaultJSON configures the default logger to use a JSON handler
func SetDefaultJSON(level slog.Level) {
	h := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:       level,
		AddSource:   level == slog.LevelDebug,
		ReplaceAttr: nil,
	})

	slog.SetDefault(slog.New(h))
}
