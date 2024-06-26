package handlers

import (
	"context"
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
		sessionToken := r.PathValue("sessionToken")

		body, err := utils.Decode[updateSessionRequest](r)
		if err != nil {
			utils.HttpFormattedError(w, r, http.StatusBadRequest, err.Error(), "failed to parse the request body")
			return
		}

		if err = validateUpdateSessionRequest(body); err != nil {
			utils.HttpFormattedError(w, r, http.StatusBadRequest, err.Error(), "invalid request body")
			return
		}

		ctx := context.Background()
		transaction, err := pool.Begin(ctx)
		if err != nil {
			utils.HttpInternalServerError(w, r, err.Error())
			return
		}

		defer transaction.Rollback(ctx)

		sql1 := "select id, user_id, expires, session_token from sessions where session_token=$1"

		var oldSession models.Session
		err = transaction.QueryRow(ctx, sql1, sessionToken).Scan(
			&oldSession.Id, &oldSession.UserId, &oldSession.Expires, &oldSession.SessionToken,
		)
		if errors.Is(err, pgx.ErrNoRows) {
			utils.HttpFormattedError(w, r, http.StatusNotFound, err.Error(), nil)
			return
		}
		if err != nil {
			utils.HttpInternalServerError(w, r, err.Error())
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
		sql2 := `update sessions
set user_id=$2,
    expires=$3
where session_token = $1
returning
    id, user_id, expires, session_token`

		if err = transaction.QueryRow(ctx, sql2, sessionToken, toUpdate.UserId, toUpdate.Expires).Scan(
			&updatedSession.Id, &updatedSession.UserId, &updatedSession.Expires, &updatedSession.SessionToken,
		); err != nil {
			utils.HttpInternalServerError(w, r, err.Error())
			return
		}

		if err = transaction.Commit(ctx); err != nil {
			utils.HttpInternalServerError(w, r, err.Error())
			return
		}

		if err = utils.Encode(w, http.StatusOK, updatedSession); err != nil {
			utils.HttpInternalServerError(w, r, err.Error())
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
