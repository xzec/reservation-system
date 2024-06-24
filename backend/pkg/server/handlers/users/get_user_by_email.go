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

func GetUserByEmailHandler(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User

		sql := "select id, email, email_verified, name, image from users where email=$1"

		err := pool.QueryRow(
			context.Background(), sql, r.PathValue("email"),
		).Scan(
			&user.Id, &user.Email, &user.EmailVerified, &user.Name, &user.Image,
		)
		if errors.Is(err, pgx.ErrNoRows) {
			utils.HttpFormattedError(w, r, http.StatusNotFound, err.Error(), nil)
			return
		}
		if err != nil {
			utils.HttpInternalServerError(w, r, err.Error())
			return
		}

		if err = utils.Encode(w, http.StatusOK, user); err != nil {
			utils.HttpInternalServerError(w, r, err.Error())
			return
		}
	}
}
