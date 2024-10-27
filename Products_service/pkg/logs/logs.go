package logs

import (
	"log"
	"log/slog"
	"os"
)

func NewLogger() *slog.Logger {
	logFile, err := os.OpenFile("pkg/logs/app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	handle := slog.NewJSONHandler(logFile, nil)

	logger := slog.New(handle)

	return logger
}
