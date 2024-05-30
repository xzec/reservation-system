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

func GetSessionAndUserHandler(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		result := struct {
			Session models.Session `json:"session"`
			User    models.User    `json:"user"`
		}{
			Session: models.Session{},
			User:    models.User{},
		}

		sql1 := `select s.id as session_id,
       s.user_id,
       s.expires,
       s.session_token,
       u.id as user_id,
       u.email,
       u.email_verified,
       u.name,
       u.image
from sessions s
         join users u on s.user_id = u.id
where s.session_token = $1`

		err := pool.QueryRow(
			context.Background(), sql1, r.PathValue("sessionToken"),
		).Scan(
			&result.Session.Id, &result.Session.UserId, &result.Session.Expires, &result.Session.SessionToken, &result.User.Id, &result.User.Email, &result.User.EmailVerified, &result.User.Name, &result.User.Image,
		)
		if errors.Is(err, pgx.ErrNoRows) {
			w.WriteHeader(http.StatusNotFound)
			_, err = w.Write([]byte("null"))
			return
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		res, err := json.Marshal(result)
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
