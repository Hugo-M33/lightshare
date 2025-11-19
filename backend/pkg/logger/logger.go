package logger

import (
	"log/slog"
	"os"
)

var defaultLogger *slog.Logger

// Init initializes the logger with the specified level
func Init(level string) {
	var logLevel slog.Level
	switch level {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level: logLevel,
	}

	handler := slog.NewJSONHandler(os.Stdout, opts)
	defaultLogger = slog.New(handler)
	slog.SetDefault(defaultLogger)
}

// Get returns the default logger
func Get() *slog.Logger {
	if defaultLogger == nil {
		Init("info")
	}
	return defaultLogger
}

// Debug logs at debug level
func Debug(msg string, args ...any) {
	Get().Debug(msg, args...)
}

// Info logs at info level
func Info(msg string, args ...any) {
	Get().Info(msg, args...)
}

// Warn logs at warn level
func Warn(msg string, args ...any) {
	Get().Warn(msg, args...)
}

// Error logs at error level
func Error(msg string, args ...any) {
	Get().Error(msg, args...)
}

// With returns a logger with the specified attributes
func With(args ...any) *slog.Logger {
	return Get().With(args...)
}
