package handlers

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http"
	"rs/pkg/server/models"
	"rs/pkg/server/utils"
)

type createSessionRequest struct {
	UserId       *string `json:"userId"`
	Expires      *string `json:"expires"`
	SessionToken *string `json:"sessionToken"`
}

func CreateSessionHandler(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := utils.Decode[createSessionRequest](r)
		if err != nil {
			utils.HttpFormattedError(w, r, http.StatusBadRequest, err.Error(), "failed to parse the request body")
			return
		}

		if err = validateCreateSessionRequest(body); err != nil {
			utils.HttpFormattedError(w, r, http.StatusBadRequest, err.Error(), "invalid request body")
			return
		}

		var session models.Session
		sql := `insert into sessions(user_id, expires, session_token)
values ($1, $2, $3)
returning id, user_id, expires, session_token`

		err = pool.QueryRow(context.Background(), sql,
			body.UserId, body.Expires, body.SessionToken,
		).Scan(
			&session.Id, &session.UserId, &session.Expires, &session.SessionToken,
		)
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

func validateCreateSessionRequest(body createSessionRequest) error {
	if body.UserId == nil {
		return errors.New("userId is required")
	}
	if !utils.IsValidUUID(*body.UserId) {
		return errors.New("invalid user id")
	}
	if body.Expires == nil {
		return errors.New("expires is required")
	}
	if body.SessionToken == nil {
		return errors.New("sessionToken is required")
	}
	return nil
}
