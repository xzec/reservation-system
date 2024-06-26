package handlers

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http"
	"rs/pkg/server/models"
	"rs/pkg/server/utils"
)

func GetSessionAndUserHandler(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionToken := r.PathValue("sessionToken")

		sessionAndUser := struct {
			Session models.Session `json:"session"`
			User    models.User    `json:"user"`
		}{
			Session: models.Session{},
			User:    models.User{},
		}

		sql := `select s.id as session_id,
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

		err := pool.QueryRow(context.Background(), sql, sessionToken).Scan(
			&sessionAndUser.Session.Id,
			&sessionAndUser.Session.UserId,
			&sessionAndUser.Session.Expires,
			&sessionAndUser.Session.SessionToken,
			&sessionAndUser.User.Id,
			&sessionAndUser.User.Email,
			&sessionAndUser.User.EmailVerified,
			&sessionAndUser.User.Name,
			&sessionAndUser.User.Image,
		)
		if errors.Is(err, pgx.ErrNoRows) {
			utils.HttpFormattedError(w, r, http.StatusNotFound, err.Error(), nil)
			return
		}
		if err != nil {
			utils.HttpInternalServerError(w, r, err.Error())
			return
		}

		if err = utils.Encode(w, http.StatusOK, sessionAndUser); err != nil {
			utils.HttpInternalServerError(w, r, err.Error())
			return
		}
	}
}
