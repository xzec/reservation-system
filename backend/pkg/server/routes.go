package server

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http"
	"rs/pkg/server/handlers"
	accounts "rs/pkg/server/handlers/accounts"
	sessions "rs/pkg/server/handlers/sessions"
	users "rs/pkg/server/handlers/users"
	verificationTokens "rs/pkg/server/handlers/verification_tokens"
)

func addRoutes(
	mux *http.ServeMux,
	pool *pgxpool.Pool,
) {
	mux.HandleFunc("POST /auth/users", users.CreateUserHandler(pool))
	mux.HandleFunc("GET /auth/users/{id}", users.GetUserHandler(pool))
	mux.HandleFunc("GET /auth/users/email/{email}", users.GetUserByEmailHandler(pool))
	mux.HandleFunc("GET /auth/users/account/{provider}/{providerAccountId}", users.GetUserByAccountHandler(pool))
	mux.HandleFunc("PATCH /auth/users/{id}", users.UpdateUserHandler(pool))
	mux.HandleFunc("DELETE /auth/users/{id}", users.DeleteUserHandler(pool))
	mux.HandleFunc("POST /auth/accounts", accounts.LinkAccountHandler(pool))
	mux.HandleFunc("DELETE /auth/accounts/{provider}/{providerAccountId}", accounts.UnlinkAccountHandler(pool))
	mux.HandleFunc("POST /auth/sessions", sessions.CreateSessionHandler(pool))
	mux.HandleFunc("GET /auth/sessions/{sessionToken}", sessions.GetSessionAndUserHandler(pool))
	mux.HandleFunc("PATCH /auth/sessions/{sessionToken}", sessions.UpdateSessionHandler(pool))
	mux.HandleFunc("DELETE /auth/sessions/{sessionToken}", sessions.DeleteSessionHandler(pool))
	mux.HandleFunc("POST /auth/verification-tokens", verificationTokens.CreateVerificationTokenHandler(pool))
	mux.HandleFunc("POST /auth/verification-tokens/use", verificationTokens.UseVerificationTokenHandler(pool))
	mux.HandleFunc("/", handlers.NotFoundHandler())
}
