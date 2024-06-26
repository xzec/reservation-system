package handlers

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http"
	"rs/pkg/server/models"
	"rs/pkg/server/utils"
)

type linkAccountRequest struct {
	UserId            *string `json:"userId"`
	Type              *string `json:"type"`
	Provider          *string `json:"provider"`
	ProviderAccountId *string `json:"providerAccountId"`
	RefreshToken      *string `json:"refreshToken"`
	AccessToken       *string `json:"accessToken"`
	ExpiresAt         *int64  `json:"expiresAt"`
	IdToken           *string `json:"idToken"`
	Scope             *string `json:"scope"`
	SessionState      *string `json:"sessionState"`
	TokenType         *string `json:"tokenType"`
}

func LinkAccountHandler(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := utils.Decode[linkAccountRequest](r)
		if err != nil {
			utils.HttpFormattedError(w, r, http.StatusBadRequest, err.Error(), "failed to parse the request body")
			return
		}

		if err = validateLinkAccountRequest(body); err != nil {
			utils.HttpFormattedError(w, r, http.StatusBadRequest, err.Error(), "invalid request body")
			return
		}

		var newAccount models.Account
		sql := `insert into accounts(user_id, type, provider, provider_account_id, refresh_token, access_token, expires_at, id_token,
                     scope, session_state, token_type)
values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
returning id, user_id, type, provider, provider_account_id, refresh_token, access_token, expires_at, id_token, scope, session_state, token_type`

		err = pool.QueryRow(context.Background(), sql,
			body.UserId, body.Type, body.Provider, body.ProviderAccountId, body.RefreshToken, body.AccessToken, body.ExpiresAt, body.IdToken, body.Scope, body.SessionState, body.TokenType,
		).Scan(
			&newAccount.Id, &newAccount.UserId, &newAccount.Type, &newAccount.Provider, &newAccount.ProviderAccountId, &newAccount.RefreshToken, &newAccount.AccessToken, &newAccount.ExpiresAt, &newAccount.IdToken, &newAccount.Scope, &newAccount.SessionState, &newAccount.TokenType,
		)
		if err != nil {
			utils.HttpInternalServerError(w, r, err.Error())
			return
		}

		if err = utils.Encode(w, http.StatusOK, newAccount); err != nil {
			utils.HttpInternalServerError(w, r, err.Error())
			return
		}
	}
}

func validateLinkAccountRequest(body linkAccountRequest) error {
	if body.UserId == nil {
		return fmt.Errorf("userId is required")
	}
	if !utils.IsValidUUID(*body.UserId) {
		return fmt.Errorf("userId is not valid")
	}
	if body.Type == nil {
		return fmt.Errorf("type is required")
	}
	if body.Provider == nil {
		return fmt.Errorf("provider is required")
	}
	if body.ProviderAccountId == nil {
		return fmt.Errorf("providerAccountId is required")
	}
	return nil
}
