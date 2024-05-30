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

type useVerificationTokenRequest struct {
	Identifier *string `json:"identifier"`
	Token      *string `json:"token"`
}

func UseVerificationTokenHandler(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body useVerificationTokenRequest

		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			http.Error(w, "Failed to parse the request body:"+err.Error(), http.StatusBadRequest)
			return
		}

		if err = validateUseVerificationTokenRequest(body); err != nil {
			http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
			return
		}

		sql := `
delete from verification_tokens
where identifier=$1 and token=$2
returning identifier, expires, token`

		var verificationToken models.VerificationToken
		err = pool.QueryRow(
			context.Background(), sql, body.Identifier, body.Token,
		).Scan(
			&verificationToken.Identifier, &verificationToken.Expires, &verificationToken.Token,
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

		res, err := json.Marshal(verificationToken)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		if _, err = w.Write(res); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func validateUseVerificationTokenRequest(body useVerificationTokenRequest) error {
	if body.Identifier == nil {
		return errors.New("identifier is required")
	}
	if body.Token == nil {
		return errors.New("token is required")
	}
	return nil
}
