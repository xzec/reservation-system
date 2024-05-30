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
		ctx := context.Background()

		tx, err := pool.Begin(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		defer tx.Rollback(ctx)

		sql1 := "select id, session_token, user_id, expires from sessions where session_token=$1"

		var session models.Session
		err = tx.QueryRow(ctx, sql1, sessionToken).Scan(&session.Id, &session.SessionToken, &session.UserId, &session.Expires)
		if errors.Is(err, pgx.ErrNoRows) {
			w.WriteHeader(http.StatusNotFound)
			_, err = w.Write([]byte("null"))
			return
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		sql2 := "delete from sessions where session_token=$1"

		commandTag, err := tx.Exec(ctx, sql2, sessionToken)
		if err != nil {
			http.Error(w, "Failed to delete the session.", http.StatusInternalServerError)
			return
		}
		if commandTag.RowsAffected() == 0 {
			http.Error(w, "Session was not found when deleting.", http.StatusNotFound)
			return
		}

		if err = tx.Commit(ctx); err != nil {
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
