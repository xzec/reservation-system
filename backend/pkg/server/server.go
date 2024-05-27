package server

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"net/http"
	"os"
	handlers "rs/pkg/server/handlers/users"
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

	router.HandleFunc("POST /auth/users", handlers.CreateUserHandler(pool))

	router.HandleFunc("GET /auth/users/{id}", handlers.GetUserHandler(pool))

	router.HandleFunc("GET /auth/users/email/{email}", handlers.GetUserByEmailHandler(pool))

	router.HandleFunc("PATCH /auth/users/{id}", handlers.UpdateUserHandler(pool))

	router.HandleFunc("DELETE /auth/users/{id}", handlers.DeleteUserHandler(pool))

	fmt.Println("Serving and listening at port " + os.Getenv("PORT"))
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), router))
}
