package utils

import (
	"log/slog"
	"os"
)

func LogFatalf(logger *slog.Logger, msg string, args ...interface{}) {
	logger.Error(msg, args...)
	os.Exit(1)
}
