package base

import (
	"log/slog"
	"os"
	"sync"
)

var logOnce sync.Once
var Log *slog.Logger

func newLogger() *slog.Logger {

	if Log == nil {
		logOnce.Do(func() {
			Log = slog.New(slog.NewJSONHandler(os.Stdout, nil))
		})
	}

	return Log
}

func SetLogger(logger *slog.Logger) {
	Log = logger
}

func init() {
	Log = newLogger()
}
