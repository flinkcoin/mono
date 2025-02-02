package pkg

import (
	"log/slog"
	"os"
	"sync"
)

var loggerOnce sync.Once
var logger *slog.Logger

func NewLogger() *slog.Logger {

	if logger == nil {
		loggerOnce.Do(func() {
			logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
		})
	}

	return logger
}
