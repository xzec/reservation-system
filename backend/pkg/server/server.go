package server

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http"
	"rs/pkg/server/middleware"
)

func NewServer(pool *pgxpool.Pool) http.Handler {
	mux := http.NewServeMux()

	addRoutes(mux, pool)

	var handler http.Handler = mux
	handler = middleware.LoggingMiddleware(handler)
	handler = middleware.CheckAuthHeaders(handler)
	return handler
}
