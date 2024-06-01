package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

type wrappedResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *wrappedResponseWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.statusCode = statusCode
}

func LoggingMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wrapped := &wrappedResponseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		start := time.Now()
		handler.ServeHTTP(wrapped, r)
		slog.Info("request",
			"status_code", wrapped.statusCode,
			"method", r.Method,
			"path", r.URL.Path,
			"address", r.RemoteAddr,
			"duration_ms", time.Since(start).Milliseconds(),
		)
	})
}
