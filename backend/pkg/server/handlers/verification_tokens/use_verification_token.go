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

type useVerificationTokenRequest struct {
	Identifier *string `json:"identifier"`
	Token      *string `json:"token"`
}

func UseVerificationTokenHandler(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := utils.Decode[useVerificationTokenRequest](r)
		if err != nil {
			utils.HttpFormattedError(w, r, http.StatusBadRequest, err.Error(), "failed to parse the request body")
			return
		}

		if err = validateUseVerificationTokenRequest(body); err != nil {
			utils.HttpFormattedError(w, r, http.StatusBadRequest, err.Error(), "invalid request body")
			return
		}

		sql := `delete
from verification_tokens
where identifier = $1
  and token = $2
returning identifier, expires, token`

		var verificationToken models.VerificationToken
		err = pool.QueryRow(context.Background(), sql, body.Identifier, body.Token).Scan(
			&verificationToken.Identifier, &verificationToken.Expires, &verificationToken.Token,
		)
		if errors.Is(err, pgx.ErrNoRows) {
			utils.HttpFormattedError(w, r, http.StatusNotFound, err.Error(), nil)
			return
		}
		if err != nil {
			utils.HttpInternalServerError(w, r, err.Error())
			return
		}

		if err = utils.Encode(w, http.StatusOK, verificationToken); err != nil {
			utils.HttpInternalServerError(w, r, err.Error())
			return
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
