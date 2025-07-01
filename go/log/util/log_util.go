package util

import (
	"context"
	"os"
	"runtime"
	"strings"

	"log/slog"
)

var logger *slog.Logger

func Initialize() {
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == "level" {
				levelStr := strings.ToLower(a.Value.String())
				var coloredLevel string
				switch levelStr {
				case "debug":
					coloredLevel = "\033[44mDEBUG\033[0m" // Blue background
				case "info":
					coloredLevel = "\033[42mINFO\033[0m" // Green background
				case "warn", "warning":
					coloredLevel = "\033[43mWARN\033[0m" // Yellow background
				case "error":
					coloredLevel = "\033[41mERROR\033[0m" // Red background
				default:
					coloredLevel = levelStr
				}
				a.Value = slog.StringValue(coloredLevel)
			}
			return a
		},
	})

	logger = slog.New(handler)
}

func getCallerInfo() (file string, line int) {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		return "???", 0
	}
	return file, line
}

func commonAttrs(location, processID string) []slog.Attr {
	fields := []slog.Attr{
		slog.String("location", location),
		slog.String("process_id", processID),
	}

	return fields
}

func attrsToArgs(attrs []slog.Attr) []any {
	args := make([]any, len(attrs))
	for i, attr := range attrs {
		args[i] = attr
	}
	return args
}

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

func LogSuccess(ctx context.Context, msg, location, processID string) {
	if logger == nil {
		Initialize()
	}
	logger.InfoContext(ctx, msg, attrsToArgs(commonAttrs(location, processID))...)
}

func LogWarn(ctx context.Context, msg, location, processID string) {
	if logger == nil {
		Initialize()
	}
	logger.WarnContext(ctx, msg, attrsToArgs(commonAttrs(location, processID))...)
}

func LogTask(ctx context.Context, msg, location, processID string) {
	if logger == nil {
		Initialize()
	}
	logger.DebugContext(ctx, msg, attrsToArgs(commonAttrs(location, processID))...)
}

func LogInfo(ctx context.Context, msg, location, processID string) {
	if logger == nil {
		Initialize()
	}
	logger.InfoContext(ctx, msg, attrsToArgs(commonAttrs(location, processID))...)
}
