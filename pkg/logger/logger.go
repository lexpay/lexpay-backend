package logger

import (
	"log/slog"
	"os"
)

const (
	EnvDevelopment = "development"
	EnvProduction  = "production"
)

func SetupLogger(env string) {
	var logger *slog.Logger

	switch env {
	case EnvProduction:
		// Production use JSON for machine readability
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))
	case EnvDevelopment:
		// Development use Text for human readability
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
	default:
		// Fallback to basic JSON
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))
	}

	// Set as default logger
	slog.SetDefault(logger)
}
