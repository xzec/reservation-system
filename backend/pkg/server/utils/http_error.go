package utils

import (
	"log/slog"
	"net/http"
)

func HttpInternalServerError(w http.ResponseWriter, r *http.Request, internalErrorMessage string) {
	slog.ErrorContext(r.Context(), "http error",
		"status", http.StatusInternalServerError,
		"method", r.Method,
		"path", r.URL.Path,
		"address", r.RemoteAddr,
		"error", internalErrorMessage)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
