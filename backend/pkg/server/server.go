package server

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http"
)

func NewServer(pool *pgxpool.Pool) http.Handler {
	mux := http.NewServeMux()

	addRoutes(mux, pool)

	return mux
}
