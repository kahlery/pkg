package util

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"runtime"
)

// ANSI color codes for background colors
const (
	bgRed    = "\033[41m" // ERROR
	bgYellow = "\033[43m" // WARN
	bgGreen  = "\033[42m" // INFO/SUCCESS
	bgBlue   = "\033[44m" // DEBUG/TASK
	bgReset  = "\033[0m"  // Reset
	fgBlack  = "\033[30m" // Black text for better contrast
)

// ColoredHandler is a custom slog.Handler that adds background colors to levels
type ColoredHandler struct {
	writer io.Writer
	opts   *slog.HandlerOptions
}

// NewColoredHandler creates a new colored handler
func NewColoredHandler(w io.Writer, opts *slog.HandlerOptions) *ColoredHandler {
	if opts == nil {
		opts = &slog.HandlerOptions{}
	}
	return &ColoredHandler{
		writer: w,
		opts:   opts,
	}
}

// Enabled reports whether the handler handles records at the given level
func (h *ColoredHandler) Enabled(ctx context.Context, level slog.Level) bool {
	minLevel := slog.LevelInfo
	if h.opts.Level != nil {
		minLevel = h.opts.Level.Level()
	}
	return level >= minLevel
}

// Handle formats and writes the log record
func (h *ColoredHandler) Handle(ctx context.Context, r slog.Record) error {
	// Get the appropriate background color for the level
	var bgColor string
	switch r.Level {
	case slog.LevelError:
		bgColor = bgRed
	case slog.LevelWarn:
		bgColor = bgYellow
	case slog.LevelInfo:
		bgColor = bgGreen
	case slog.LevelDebug:
		bgColor = bgBlue
	default:
		bgColor = bgGreen
	}

	// Create the colored level string
	coloredLevel := fmt.Sprintf("%s%s %-5s %s", bgColor, fgBlack, r.Level.String(), bgReset)

	// Build the log entry with custom time format
	logData := map[string]interface{}{
		"time":    r.Time.Format("2006-01-02 15:04:05"),
		"level":   coloredLevel,
		"message": r.Message,
	}

	// Add all attributes
	r.Attrs(func(a slog.Attr) bool {
		logData[a.Key] = a.Value.Any()
		return true
	})

	// Convert to JSON and write
	jsonBytes, err := json.Marshal(logData)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintln(h.writer, string(jsonBytes))
	return err
}

// WithAttrs returns a new handler with the given attributes
func (h *ColoredHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	// For simplicity, return the same handler
	// In a full implementation, you'd want to store these attrs
	return h
}

// WithGroup returns a new handler with the given group name
func (h *ColoredHandler) WithGroup(name string) slog.Handler {
	// For simplicity, return the same handler
	// In a full implementation, you'd want to handle grouping
	return h
}

// --------------------------------------------------------------------
var logger *slog.Logger

// --------------------------------------------------------------------
func Initialize() {
	handler := NewColoredHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug, // Allow all levels
	})
	logger = slog.New(handler)
}

// --------------------------------------------------------------------
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

// --------------------------------------------------------------------
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
