package handlers

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http"
	"net/mail"
	"rs/pkg/server/models"
	"rs/pkg/server/utils"
)

type createUserRequest struct {
	Email         *string `json:"email"`
	EmailVerified *string `json:"emailVerified"`
	Name          *string `json:"name"`
	Image         *string `json:"image"`
}

func CreateUserHandler(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := utils.Decode[createUserRequest](r)
		if err != nil {
			utils.HttpFormattedError(w, r, http.StatusBadRequest, err.Error(), "failed to parse the request body")
			return
		}

		if err = validateCreateUserRequest(body); err != nil {
			utils.HttpFormattedError(w, r, http.StatusBadRequest, err.Error(), "invalid request body")
			return
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
			utils.HttpInternalServerError(w, r, err.Error())
			return
		}

		if err = utils.Encode(w, http.StatusOK, newUser); err != nil {
			utils.HttpInternalServerError(w, r, err.Error())
			return
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
