package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http"
	"net/mail"
	"rs/pkg/server/models"
)

type createUserRequest struct {
	Email         *string `json:"email"`
	EmailVerified *string `json:"emailVerified"`
	Name          *string `json:"name"`
	Image         *string `json:"image"`
}

func CreateUserHandler(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body createUserRequest

		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			http.Error(w, "Failed to parse the request body:"+err.Error(), http.StatusBadRequest)
			return
		}

		if err = validateCreateUserRequest(body); err != nil {
			http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		}

		var newUser models.User
		sql := `insert into users(email, email_verified, name, image)
values ($1, $2, $3, $4)
returning id, email, email_verified, name, image`

		err = pool.QueryRow(context.Background(), sql,
			body.Email, body.EmailVerified, body.Name, body.Image,
		).Scan(
			&newUser.Id, &newUser.Email, &newUser.EmailVerified, &newUser.Name, &newUser.Image,
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		res, err := json.Marshal(newUser)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if _, err = w.Write(res); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func validateCreateUserRequest(body createUserRequest) error {
	if body.Email == nil {
		return fmt.Errorf("email is a required field")
	}
	if _, err := mail.ParseAddress(*body.Email); err != nil {
		return fmt.Errorf("email is not valid")
	}
	return nil
}
