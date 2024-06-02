package utils

import (
	"log/slog"
	"net/http"
)

func HttpInternalServerError(w http.ResponseWriter, r *http.Request, internalError string) {
	logHttpError(r, slog.LevelError, http.StatusInternalServerError, internalError)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func HttpFormattedError(w http.ResponseWriter, r *http.Request, status int, error string, response any) {
	var logLevel slog.Level
	switch {
	case status >= 500:
		logLevel = slog.LevelError
	case status >= 400:
		logLevel = slog.LevelWarn
	default:
		logLevel = slog.LevelInfo
	}

	logHttpError(r, logLevel, status, error)

	if err := Encode(w, status, response); err != nil {
		HttpInternalServerError(w, r, err.Error())
		return
	}
}

func logHttpError(r *http.Request, level slog.Level, status int, error string) {
	slog.Log(r.Context(), level, "http error",
		"status", status,
		"method", r.Method,
		"path", r.URL.Path,
		"address", r.RemoteAddr,
		"error", error)
}
