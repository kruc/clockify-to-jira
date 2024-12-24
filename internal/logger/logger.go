package logger

import (
	"log/slog"
)

var logger *slog.Logger

func InitializeLogger() *slog.Logger {
	prettyHandler := NewHandler(&slog.HandlerOptions{
		Level:       slog.LevelInfo,
		AddSource:   false,
		ReplaceAttr: nil,
	})

	logger = slog.New(prettyHandler)

	return logger
}
