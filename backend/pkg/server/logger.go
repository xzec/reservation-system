package server

import (
	"context"
	"io"
	"log/slog"
)

type ContextHandler struct {
	slog.Handler
}

func (h ContextHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return ContextHandler{h.Handler.WithAttrs(attrs)}
}

func (h ContextHandler) WithGroup(name string) slog.Handler {
	return ContextHandler{h.Handler.WithGroup(name)}
}

func (h ContextHandler) Handle(ctx context.Context, r slog.Record) error {
	traceId, ok := ctx.Value("traceId").(string)
	if ok {
		attr := slog.Attr{
			Key:   "traceId",
			Value: slog.StringValue(traceId),
		}
		r.AddAttrs(attr)
	}
	return h.Handler.Handle(ctx, r)
}

func createLogger(w io.Writer) *slog.Logger {
	jsonHandler := slog.NewJSONHandler(w, nil)
	ctxHandler := ContextHandler{jsonHandler}

	logger := slog.New(ctxHandler)
	return logger
}
