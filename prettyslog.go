package prettyslog

import (
	"context"
	"encoding/json"
	"os"
	"io"
	"log"
	"log/slog"

	"github.com/fatih/color"
)

type prettyLoggerOptions struct {
	slogOpts *slog.HandlerOptions
}

type prettyLoggerHandler struct {
	opts    prettyLoggerOptions
	handler slog.Handler
	l       *log.Logger
	attrs   []slog.Attr
}

func NewPrettyLogger() *slog.Logger {
	opts := prettyLoggerOptions{
		slogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	return slog.New(
		opts.configurePrettyLogger(os.Stdout),
	)
}

func (opts prettyLoggerOptions) configurePrettyLogger(out io.Writer) *prettyLoggerHandler {
	return &prettyLoggerHandler{
		opts:    opts,
		handler: slog.NewJSONHandler(out, opts.slogOpts),
		l:       log.New(out, "", 0),
	}
}

func (h *prettyLoggerHandler) Enabled(_ context.Context, l slog.Level) bool {
	return l >= slog.LevelDebug
}

func (h *prettyLoggerHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &prettyLoggerHandler{
		handler: h.handler,
		l:       h.l,
		attrs:   attrs,
	}
}

func (h *prettyLoggerHandler) WithGroup(name string) slog.Handler {
	return &prettyLoggerHandler{
		handler: h.handler.WithGroup(name),
		l:       h.l,
	}
}

func (h *prettyLoggerHandler) Handle(_ context.Context, r slog.Record) error {
	level := r.Level.String()

	switch r.Level {
	case slog.LevelDebug:
		level = color.BlueString(level)
	case slog.LevelInfo:
		level = color.YellowString(level)
	case slog.LevelWarn:
		level = color.MagentaString(level)
	case slog.LevelError:
		level = color.RedString(level)
	}

	extra := make(map[string]any, r.NumAttrs())
	// r.Attrs(
	// 	func(attr slog.Attr) bool {
	// 		extra[attr.Key] = attr.Value.Any()
	// 		return true
	// 	},
	// )

	for _, a := range h.attrs {
		extra[a.Key] = a.Value.Any()
	}

	var b []byte
	var err error

	if len(extra) > 0 {
		b, err = json.MarshalIndent(extra, "", "\t")
		if err != nil {
			return err
		}
	}

	h.l.Println(
		r.Time.Format("[2006-01-02 15:04:05.000]"),
		level,
		color.CyanString(r.Message),
		color.WhiteString(string(b)),
	)

	return nil
}
