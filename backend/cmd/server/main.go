package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"net/http"
	"os"
	"time"
)

type User struct {
	Id            string    `json:"id"`
	Email         string    `json:"email"`
	EmailVerified string    `json:"emailVerified"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

func main() {
	fmt.Println("DATABASE_URL:", os.Getenv("DATABASE_URL"))
	dbPool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer dbPool.Close()

	fmt.Println("Starting server at 8080...")
	router := http.NewServeMux()

	router.HandleFunc("GET /users", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("GET /users")
		rows, err := dbPool.Query(context.Background(), "select id, email from users")
		defer rows.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		var dbUsers []User
		for rows.Next() {
			var user User
			err = rows.Scan(&user.Id, &user.Email)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			dbUsers = append(dbUsers, user)
		}
		fmt.Println("dbUsers=", dbUsers)
		res, err := json.Marshal(dbUsers)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		_, err = w.Write(res)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	router.HandleFunc("GET /users/{id}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("GET /users/{id}, id=", r.PathValue("id"))
		var dbUser User
		err = dbPool.QueryRow(context.Background(), "select id, email, created_at from users where id=$1", r.PathValue("id")).Scan(&dbUser.Id, &dbUser.Email, &dbUser.CreatedAt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Println("dbUser=", dbUser)
		res, err := json.Marshal(dbUser)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		_, err = w.Write(res)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	router.HandleFunc("POST /users", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("POST /users")
		var user User
		err = json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Println("user=", user)

		err = dbPool.QueryRow(context.Background(), "insert into users(email) values($1)", user.Email).Scan()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var dbUser User
		err = dbPool.QueryRow(context.Background(), "select id, email, created_at from users where email=$1", user.Email).Scan(&dbUser.Id, &dbUser.Email, &dbUser.CreatedAt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Println("dbUser=", dbUser)

		res, err := json.Marshal(dbUser)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		_, err = w.Write(res)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	router.HandleFunc("DELETE /users/{id}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("DELETE /users/{id}", r.PathValue("id"))

		err = dbPool.QueryRow(context.Background(), "delete from users where id=$1", r.PathValue("id")).Scan()
		if err != nil {
			http.Error(w, "failed to delete user "+r.PathValue("id"), http.StatusInternalServerError)
			return
		}
		_, err = w.Write([]byte("user " + r.PathValue("id") + " was deleted"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	log.Fatal(http.ListenAndServe(":8080", router))
}
