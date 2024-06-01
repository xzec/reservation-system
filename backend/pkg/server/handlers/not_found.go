package handlers

import (
	"log/slog"
	"net/http"
)

func NotFoundHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("route not found",
			"path", r.URL.Path)
		http.NotFound(w, r)
	}
}
