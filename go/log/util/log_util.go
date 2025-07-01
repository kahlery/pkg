package util

import (
	"context"
	"log/slog"
	"os"
	"runtime"
	"time"
)

var logger *slog.Logger

// CustomTime wraps time.Time and implements slog.TimeMarshaler
type CustomTime time.Time

// MarshalTime formats time as "2006-01-02 15:04:05"
func (ct CustomTime) MarshalTime() string {
	t := time.Time(ct)
	return t.Format("2006-01-02 15:04:05")
}

func Initialize() {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
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
		slog.String("time", "w"), // Use custom formatted time
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
