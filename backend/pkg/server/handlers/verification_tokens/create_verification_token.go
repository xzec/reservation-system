package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http"
	"rs/pkg/server/models"
	"time"
)

type createVerificationTokenRequest struct {
	Identifier *string    `json:"identifier"`
	Expires    *time.Time `json:"expires"`
	Token      *string    `json:"token"`
}

func CreateVerificationTokenHandler(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body createVerificationTokenRequest

		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			http.Error(w, "Failed to parse the request body:"+err.Error(), http.StatusBadRequest)
			return
		}

		if err = validateCreateVerificationTokenRequest(body); err != nil {
			http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
			return
		}

		var verificationToken models.VerificationToken
		sql := `
insert into verification_tokens(identifier, expires, token)
values ($1, $2, $3)
returning identifier, expires, token`

		err = pool.QueryRow(context.Background(), sql,
			body.Identifier, body.Expires, body.Token,
		).Scan(
			&verificationToken.Identifier, &verificationToken.Expires, &verificationToken.Token,
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		res, err := json.Marshal(verificationToken)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if _, err = w.Write(res); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func validateCreateVerificationTokenRequest(body createVerificationTokenRequest) error {
	if body.Identifier == nil {
		return errors.New("identifier is required")
	}
	if body.Expires == nil {
		return errors.New("expires is required")
	}
	if body.Token == nil {
		return errors.New("token is required")
	}
	return nil
}
