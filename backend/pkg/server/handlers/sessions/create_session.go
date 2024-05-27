package handlers

import (
	"context"
	"encoding/json"
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
		var body createSessionRequest

		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			http.Error(w, "Failed to parse the request body:"+err.Error(), http.StatusBadRequest)
			return
		}

		if err = validateCreateSessionRequest(body); err != nil {
			http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
			return
		}

		var session models.Session
		sql := `
insert into sessions(user_id, expires, session_token)
values ($1, $2, $3)
returning id, user_id, expires, session_token`

		err = pool.QueryRow(context.Background(), sql,
			body.UserId, body.Expires, body.SessionToken,
		).Scan(
			&session.Id, &session.UserId, &session.Expires, &session.SessionToken,
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		res, err := json.Marshal(session)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if _, err = w.Write(res); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
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
