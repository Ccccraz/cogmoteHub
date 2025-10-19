package logger

import (
	"log/slog"
)

func Init() {
	logger := slog.Default()
	slog.SetDefault(logger)
}
