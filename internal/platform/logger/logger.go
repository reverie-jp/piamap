package logger

import (
	"log/slog"
	"os"

	"github.com/lmittmann/tint"

	"github.com/reverie-jp/piamap/internal/config"
)

func parseLevel(level config.LogLevel) slog.Level {
	switch level {
	case config.LogLevelDebug:
		return slog.LevelDebug
	case config.LogLevelInfo:
		return slog.LevelInfo
	case config.LogLevelWarn:
		return slog.LevelWarn
	case config.LogLevelError:
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func Init(cfg *config.Config) {
	var handler slog.Handler

	switch cfg.Env {
	case config.EnvDevelopment:
		handler = tint.NewHandler(os.Stdout, &tint.Options{
			Level:      parseLevel(cfg.Log.Level),
			TimeFormat: "15:04:05.000",
		})
	default:
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: parseLevel(cfg.Log.Level),
		})
	}

	slog.SetDefault(slog.New(handler))
}
