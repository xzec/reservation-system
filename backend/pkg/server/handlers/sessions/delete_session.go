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

func DeleteSessionHandler(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionToken := r.PathValue("sessionToken")

		sql := `
delete from sessions
where session_token=$1
returning id, session_token, user_id, expires`

		var session models.Session
		err := pool.QueryRow(context.Background(), sql, sessionToken).Scan(&session.Id, &session.SessionToken, &session.UserId, &session.Expires)
		if errors.Is(err, pgx.ErrNoRows) {
			w.WriteHeader(http.StatusNotFound)
			_, err = w.Write([]byte("null"))
			return
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		res, err := json.Marshal(session)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		if _, err = w.Write(res); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
