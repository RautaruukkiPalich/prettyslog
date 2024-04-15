package prettyslog

import (
	"context"
	"encoding/json"
	"reflect"

	"io"
	"log"
	"log/slog"
	"os"

	"github.com/fatih/color"
	"github.com/fatih/structs"
)

type prettyLoggerOptions struct {
	slogOpts *slog.HandlerOptions
}

type prettyLoggerHandler struct {
	opts    prettyLoggerOptions
	handler slog.Handler
	l       *log.Logger
	attrs   []slog.Attr
	indent  string
}

func NewPrettyLogger(indent string) *slog.Logger {
	opts := prettyLoggerOptions{
		slogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	return slog.New(
		opts.configurePrettyLogger(
			os.Stdout,
			indent,
		),
	)
}

func (opts prettyLoggerOptions) configurePrettyLogger(out io.Writer, indent string) *prettyLoggerHandler {
	return &prettyLoggerHandler{
		opts: opts,
		handler: slog.NewJSONHandler(
			out,
			opts.slogOpts,
		),
		l:      log.New(out, "", 0),
		indent: indent,
	}
}

func (h *prettyLoggerHandler) Enabled(_ context.Context, l slog.Level) bool {
	return l >= slog.LevelDebug
}

func (h *prettyLoggerHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &prettyLoggerHandler{
		handler: h.handler.WithAttrs(attrs),
		l:       h.l,
		attrs:   attrs,
		indent:  h.indent,
	}
}

func (h *prettyLoggerHandler) WithGroup(name string) slog.Handler {
	return &prettyLoggerHandler{
		handler: h.handler.WithGroup(name),
		l:       h.l,
		indent:  h.indent,
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

	r.Attrs(
		func(attr slog.Attr) bool {
			key, val := h.getPreparedPair(&attr)
			extra[key] = val

			return true
		},
	)

	for _, a := range h.attrs {
		key, value := h.getPreparedPair(&a)
		extra[key] = value
	}

	var b []byte
	var err error

	if len(extra) > 0 {
		b, err = json.MarshalIndent(extra, "", h.indent)
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

func (h *prettyLoggerHandler) getPreparedPair(a *slog.Attr) (string, any) {
	switch a.Value.Kind() {

	case slog.KindBool:
		return a.Key, a.Value.Bool()
	case slog.KindDuration:
		return a.Key, a.Value.Duration()
	case slog.KindFloat64:
		return a.Key, a.Value.Float64()
	case slog.KindInt64:
		return a.Key, a.Value.Int64()
	case slog.KindString:
		return a.Key, a.Value.String()
	case slog.KindTime:
		return a.Key, a.Value.Time()
	case slog.KindUint64:
		return a.Key, a.Value.Uint64()
	case slog.KindLogValuer:
		return a.Key, a.Value.LogValuer()
	case slog.KindGroup:
		dict := map[string]any{}
		for _, k := range a.Value.Group() {
			key, val := h.getPreparedPair(&k)
			dict[key] = val
		}
		return a.Key, dict
	case slog.KindAny:
		res := a.Value.Any()
		typeOf := reflect.TypeOf(res)

		if typeOf.Kind() == reflect.Struct || typeOf.Kind() == reflect.Pointer {
			res = recursionFields(res)
			res = a.Value.Any()
		}

		return a.Key, res
	}
	return a.Key, a.Value.Any()
}

func recursionFields(fields any) map[string]any {
	res := map[string]any{}
	for _, f := range structs.Fields(fields) {
		if !f.IsExported() {
			continue
		}
		name := f.Name()
		val := f.Value()
		typeOf := reflect.TypeOf(val)

		if typeOf.Kind() == reflect.Struct || typeOf.Kind() == reflect.Pointer {
			nestedStruct, ok := f.FieldOk(name)
			if ok {
				res[name] = recursionFields(nestedStruct)
				continue
			}
		}
		res[name] = val
	}
	return res
}
