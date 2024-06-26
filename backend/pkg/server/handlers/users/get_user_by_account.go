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

func GetUserByAccountHandler(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		provider := r.PathValue("provider")
		providerAccountId := r.PathValue("providerAccountId")

		sql := `select u.id,
       u.email,
       u.email_verified,
       u.name,
       u.image
from users u
         join
     accounts a
     on u.id = a.user_id
where a.provider = $1
  and a.provider_account_id = $2`

		var user models.User
		err := pool.QueryRow(context.Background(), sql, provider, providerAccountId).Scan(
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
