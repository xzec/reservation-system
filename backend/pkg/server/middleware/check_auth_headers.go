package middleware

import (
	"log/slog"
	"net/http"
	"os"
)

func CheckAuthHeaders(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Api-Key") != os.Getenv("API_KEY") {
			slog.ErrorContext(r.Context(), "unauthorized access")
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		handler.ServeHTTP(w, r)
	})
}
