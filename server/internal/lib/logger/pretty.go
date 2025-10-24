package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"strings"

	"github.com/fatih/color"
)

type PrettyHandler struct {
	writer io.Writer
	level  slog.Level
}

func NewPrettyHandler(w io.Writer, level slog.Level) *PrettyHandler {
	return &PrettyHandler{
		writer: w,
		level:  level,
	}
}

func (h *PrettyHandler) Enabled(_ context.Context, lvl slog.Level) bool {
	return lvl >= h.level
}

func (h *PrettyHandler) Handle(_ context.Context, r slog.Record) error {
	var b strings.Builder

	// Time
	timestamp := r.Time.Format("2006-01-02 15:04:05.000")
	b.WriteString(color.New(color.FgHiBlack).Sprintf("[%s] ", timestamp))

	// Level
	var levelColor *color.Color
	switch r.Level {
	case slog.LevelDebug:
		levelColor = color.New(color.FgBlue)
	case slog.LevelInfo:
		levelColor = color.New(color.FgGreen)
	case slog.LevelWarn:
		levelColor = color.New(color.FgYellow)
	case slog.LevelError:
		levelColor = color.New(color.FgRed)
	default:
		levelColor = color.New(color.Reset)
	}
	b.WriteString(levelColor.Sprintf("[%s] ", strings.ToUpper(r.Level.String())))

	// Message
	b.WriteString(r.Message)

	// Attributes (key=value)
	r.Attrs(func(a slog.Attr) bool {
		val := fmt.Sprintf("%v", a.Value.Any())
		b.WriteString(color.New(color.FgHiBlack).Sprintf(" %s=", a.Key))
		b.WriteString(color.New(color.FgWhite).Sprintf("%s", val))
		return true
	})

	b.WriteRune('\n')

	_, err := h.writer.Write([]byte(b.String()))
	return err
}

func (h *PrettyHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h // ignoring for simplicity
}

func (h *PrettyHandler) WithGroup(name string) slog.Handler {
	return h // ignoring for simplicity
}
