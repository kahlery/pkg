package util

import (
	"context"
	"log/slog"
	"os"
	"runtime"
	"time"
)

var logger *slog.Logger

// Initialize sets up the slog logger with a JSON handler for structured logging
func Initialize() {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	logger = slog.New(handler)
}

// getCallerInfo returns the file and line number of the caller
func getCallerInfo() (file string, line int) {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		return "???", 0
	}
	return file, line
}

// commonAttrs returns common logging attributes
func commonAttrs(location, processID string) []slog.Attr {
	fields := []slog.Attr{
		slog.String("timestamp", time.Now().Format("2006-01-02 15:04:05")),
	}
	if location != "" {
		fields = append(fields, slog.String("location", location))
		fields = append(fields, slog.String("process_id", processID))
	}
	return fields
}

// attrsToArgs converts []slog.Attr to []any for variadic logging methods
func attrsToArgs(attrs []slog.Attr) []any {
	args := make([]any, len(attrs))
	for i, attr := range attrs {
		args[i] = attr
	}
	return args
}

// LogError logs an error message with caller info
func LogError(ctx context.Context, msg, location, processID string) {
	if logger == nil {
		Initialize()
	}
	file, line := getCallerInfo()
	logger.With(
		slog.String("file", file),
		slog.Int("line", line),
	).ErrorContext(ctx, msg, attrsToArgs(commonAttrs(location, processID))...)
}

// LogSuccess logs a success/info-level message
func LogSuccess(ctx context.Context, msg, location, processID string) {
	if logger == nil {
		Initialize()
	}
	logger.InfoContext(ctx, msg, attrsToArgs(commonAttrs(location, processID))...)
}

// LogWarn logs a warning message
func LogWarn(ctx context.Context, msg, location, processID string) {
	if logger == nil {
		Initialize()
	}
	logger.WarnContext(ctx, msg, attrsToArgs(commonAttrs(location, processID))...)
}

// LogTask logs a debug-level message
func LogTask(ctx context.Context, msg, location, processID string) {
	if logger == nil {
		Initialize()
	}
	logger.DebugContext(ctx, msg, attrsToArgs(commonAttrs(location, processID))...)
}

// LogInfo logs an info-level message
func LogInfo(ctx context.Context, msg, location, processID string) {
	if logger == nil {
		Initialize()
	}
	logger.InfoContext(ctx, msg, attrsToArgs(commonAttrs(location, processID))...)
}
