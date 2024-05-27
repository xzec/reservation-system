package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http"
	"net/mail"
	"rs/pkg/server/models"
	"rs/pkg/server/utils"
	"time"
)

type updateUserRequest struct {
	Email         utils.Optional[string]    `json:"email"`
	EmailVerified utils.Optional[time.Time] `json:"emailVerified,omitempty"`
	Name          utils.Optional[string]    `json:"name,omitempty"`
	Image         utils.Optional[string]    `json:"image,omitempty"`
}

func UpdateUserHandler(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId := r.PathValue("id")
		if !utils.IsValidUUID(userId) {
			http.Error(w, "Invalid user id.", http.StatusBadRequest)
			return
		}

		var body updateUserRequest

		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			http.Error(w, "Failed to parse the request body: "+err.Error(), http.StatusBadRequest)
			return
		}

		if err = validateUpdateUserRequest(body); err != nil {
			http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
			return
		}

		ctx := context.Background()
		transaction, err := pool.Begin(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		defer transaction.Rollback(ctx)

		sql1 := "select email, email_verified, name, image from users where id=$1"

		var oldUser models.User
		err = transaction.QueryRow(ctx, sql1, userId).Scan(&oldUser.Email, &oldUser.EmailVerified, &oldUser.Name, &oldUser.Image)
		if errors.Is(err, pgx.ErrNoRows) {
			w.WriteHeader(http.StatusNotFound)
			nilResponse, _ := json.Marshal(nil)
			_, err = w.Write(nilResponse)
			return
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		toUpdate := oldUser
		if body.Email.Defined {
			toUpdate.Email = body.Email.Value
		}
		if body.EmailVerified.Defined {
			toUpdate.EmailVerified = body.EmailVerified.Value
		}
		if body.Name.Defined {
			toUpdate.Name = body.Name.Value
		}
		if body.Image.Defined {
			toUpdate.Image = body.Image.Value
		}

		var updatedUser models.User
		sql2 := `
update users
set email=$2,
    email_verified=$3,
    name=$4,
    image=$5
where id=$1
returning id, email, email_verified, name, image`

		if err = transaction.QueryRow(ctx, sql2, userId, toUpdate.Email, toUpdate.EmailVerified, toUpdate.Name, toUpdate.Image).Scan(&updatedUser.Id, &updatedUser.Email, &updatedUser.EmailVerified, &updatedUser.Name, &updatedUser.Image); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err = transaction.Commit(ctx); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		res, err := json.Marshal(updatedUser)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if _, err = w.Write(res); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func validateUpdateUserRequest(body updateUserRequest) error {
	if body.Email.Defined == false {
		return nil
	}
	if body.Email.Value == nil {
		return fmt.Errorf("email cannot be null")
	}
	if _, err := mail.ParseAddress(*body.Email.Value); err != nil {
		return fmt.Errorf("email is not valid")
	}
	return nil
}
