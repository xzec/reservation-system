package handlers

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http"
	"rs/pkg/server/models"
	"rs/pkg/server/utils"
	"time"
)

type createVerificationTokenRequest struct {
	Identifier *string    `json:"identifier"`
	Expires    *time.Time `json:"expires"`
	Token      *string    `json:"token"`
}

func CreateVerificationTokenHandler(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := utils.Decode[createVerificationTokenRequest](r)
		if err != nil {
			utils.HttpFormattedError(w, r, http.StatusBadRequest, err.Error(), "failed to parse the request body")
			return
		}

		if err = validateCreateVerificationTokenRequest(body); err != nil {
			utils.HttpFormattedError(w, r, http.StatusBadRequest, err.Error(), "invalid request body")
			return
		}

		sql := `insert into verification_tokens(identifier, expires, token)
values ($1, $2, $3)
returning identifier, expires, token`

		var newVerificationToken models.VerificationToken
		err = pool.QueryRow(context.Background(), sql, body.Identifier, body.Expires, body.Token).Scan(
			&newVerificationToken.Identifier, &newVerificationToken.Expires, &newVerificationToken.Token,
		)
		if err != nil {
			utils.HttpInternalServerError(w, r, err.Error())
			return
		}

		if err = utils.Encode(w, http.StatusOK, newVerificationToken); err != nil {
			utils.HttpInternalServerError(w, r, err.Error())
			return
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
