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

func DeleteSessionHandler(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionToken := r.PathValue("sessionToken")

		sql := `delete
from sessions
where session_token = $1
returning id, session_token, user_id, expires`

		var session models.Session
		err := pool.QueryRow(context.Background(), sql, sessionToken).Scan(&session.Id, &session.SessionToken, &session.UserId, &session.Expires)
		if errors.Is(err, pgx.ErrNoRows) {
			utils.HttpFormattedError(w, r, http.StatusNotFound, err.Error(), nil)
			return
		}
		if err != nil {
			utils.HttpInternalServerError(w, r, err.Error())
			return
		}

		if err = utils.Encode(w, http.StatusOK, session); err != nil {
			utils.HttpInternalServerError(w, r, err.Error())
			return
		}
	}
}
