package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http"
	"rs/pkg/server/models"
	"rs/pkg/server/utils"
	"time"
)

type updateSessionRequest struct {
	UserId  utils.Optional[string]    `json:"userId,omitempty"`
	Expires utils.Optional[time.Time] `json:"expires,omitempty"`
}

func UpdateSessionHandler(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body updateSessionRequest

		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			http.Error(w, "Failed to parse the request body: "+err.Error(), http.StatusBadRequest)
			return
		}

		if err = validateUpdateSessionRequest(body); err != nil {
			http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
			return
		}

		ctx := context.Background()
		transaction, err := pool.Begin(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		defer transaction.Rollback(ctx)

		sql1 := "select id, user_id, expires, session_token from sessions where session_token=$1"

		var oldSession models.Session
		err = transaction.QueryRow(
			ctx, sql1, r.PathValue("sessionToken"),
		).Scan(
			&oldSession.Id, &oldSession.UserId, &oldSession.Expires, &oldSession.SessionToken,
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

		toUpdate := oldSession
		if body.UserId.Defined {
			toUpdate.UserId = body.UserId.Value
		}
		if body.Expires.Defined {
			toUpdate.Expires = body.Expires.Value
		}

		var updatedSession models.Session
		sql2 := `
update sessions
set
    user_id=$2,
    expires=$3
where
    session_token=$1
returning
    id, user_id, expires, session_token`

		if err = transaction.QueryRow(
			ctx, sql2, r.PathValue("sessionToken"), toUpdate.UserId, toUpdate.Expires,
		).Scan(
			&updatedSession.Id, &updatedSession.UserId, &updatedSession.Expires, &updatedSession.SessionToken,
		); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err = transaction.Commit(ctx); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		res, err := json.Marshal(updatedSession)
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

func validateUpdateSessionRequest(body updateSessionRequest) error {
	if body.UserId.Defined == true {
		if body.UserId.Value == nil {
			return errors.New("user id cannot be null")
		}
		if !utils.IsValidUUID(*body.UserId.Value) {
			return errors.New("invalid user id")
		}
	}
	if body.Expires.Defined == true && body.Expires.Value == nil {
		return errors.New("expires cannot be null")
	}
	return nil
}
