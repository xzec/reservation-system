package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type User struct {
	Id    string `json:"id"`
	Email string `json:"email"`
}

func main() {
	fmt.Println("Starting server at 8080...")
	router := http.NewServeMux()

	router.HandleFunc("GET /users", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("get all users")
	})

	router.HandleFunc("GET /users/{id}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("get a user by id:", r.PathValue("id"))
	})

	router.HandleFunc("POST /users", func(w http.ResponseWriter, r *http.Request) {
		var user User
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		fmt.Println("create a user:", user)
	})

	log.Fatal(http.ListenAndServe(":8080", router))
}
