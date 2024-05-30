package server

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"net/http"
	"os"
	accounts "rs/pkg/server/handlers/accounts"
	sessions "rs/pkg/server/handlers/sessions"
	users "rs/pkg/server/handlers/users"
	verificationTokens "rs/pkg/server/handlers/verification_tokens"
	"time"
)

type User struct {
	Id            string    `json:"id"`
	Email         string    `json:"email"`
	EmailVerified bool      `json:"emailVerified"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

func Start() {
	pool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer pool.Close()

	fmt.Println("Starting server...")
	router := http.NewServeMux()

	router.HandleFunc("POST /auth/users", users.CreateUserHandler(pool))

	router.HandleFunc("GET /auth/users/{id}", users.GetUserHandler(pool))

	router.HandleFunc("GET /auth/users/email/{email}", users.GetUserByEmailHandler(pool))

	router.HandleFunc("PATCH /auth/users/{id}", users.UpdateUserHandler(pool))

	router.HandleFunc("DELETE /auth/users/{id}", users.DeleteUserHandler(pool))

	router.HandleFunc("POST /auth/accounts", accounts.LinkAccountHandler(pool))

	router.HandleFunc("DELETE /auth/accounts/{provider}/{providerAccountId}", accounts.UnlinkAccountHandler(pool))

	router.HandleFunc("POST /auth/sessions", sessions.CreateSessionHandler(pool))

	router.HandleFunc("GET /auth/sessions/{sessionToken}", sessions.GetSessionAndUserHandler(pool))

	router.HandleFunc("PATCH /auth/sessions/{sessionToken}", sessions.UpdateSessionHandler(pool))

	router.HandleFunc("DELETE /auth/sessions/{sessionToken}", sessions.DeleteSessionHandler(pool))

	router.HandleFunc("POST /auth/verification-tokens", verificationTokens.CreateVerificationTokenHandler(pool))

	router.HandleFunc("POST /auth/verification-tokens/use", verificationTokens.UseVerificationTokenHandler(pool))

	fmt.Println("Serving and listening at port " + os.Getenv("PORT"))
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), router))
}
