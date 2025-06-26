package logger

import (
	"log/slog"
	"os"

	"github.com/phsym/console-slog"
	"github.com/walnuts1018/s3-oauth2-proxy/config"
)

func CreateLogger(logLevel slog.Level, logType config.LogType) *slog.Logger {
	var hander slog.Handler
	switch logType {
	case config.LogTypeText:
		hander = console.NewHandler(os.Stdout, &console.HandlerOptions{
			Level:     logLevel,
			AddSource: logLevel == slog.LevelDebug,
		})
	case config.LogTypeJSON:
		hander = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level:     logLevel,
			AddSource: logLevel == slog.LevelDebug,
		})
	}

	return slog.New(hander)
}
