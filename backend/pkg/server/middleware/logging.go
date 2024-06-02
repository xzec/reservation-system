package middleware

import (
	"context"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
	"rs/pkg/server/utils"
	"time"
)

type wrappedResponseWriter struct {
	http.ResponseWriter
	status int
}

func (w *wrappedResponseWriter) WriteHeader(status int) {
	w.ResponseWriter.WriteHeader(status)
	w.status = status
}

func LoggingMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wrapped := &wrappedResponseWriter{
			ResponseWriter: w,
			status:         http.StatusOK,
		}

		traceId := r.Header.Get("X-Trace-Id")
		if traceId == "" {
			newTraceId, err := uuid.NewRandom()
			if err != nil {
				utils.HttpInternalServerError(w, r, "failed to generate a traceId for the request")
				return
			}
			traceId = newTraceId.String()
		}

		ctx := context.WithValue(r.Context(), "traceId", traceId)

		start := time.Now()
		handler.ServeHTTP(wrapped, r.WithContext(ctx))
		duration := time.Since(start)

		slog.InfoContext(ctx, "request",
			"status", wrapped.status,
			"method", r.Method,
			"path", r.URL.Path,
			"address", r.RemoteAddr,
			"durationMs", duration.Milliseconds(),
		)
	})
}
