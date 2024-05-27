package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http"
	"rs/pkg/server/models"
)

func GetUserByEmailHandler(pool *pgxpool.Pool) (handler func(http.ResponseWriter, *http.Request)) {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User

		sql := "select id, email, email_verified, name, image from users where email=$1"

		err := pool.QueryRow(context.Background(), sql, r.PathValue("email")).Scan(&user.Id, &user.Email, &user.EmailVerified, &user.Name, &user.Image)
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

		res, err := json.Marshal(user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		_, err = w.Write(res)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
