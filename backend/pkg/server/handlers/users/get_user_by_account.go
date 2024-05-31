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
		err := pool.QueryRow(
			context.Background(), sql, provider, providerAccountId,
		).Scan(
			&user.Id, &user.Email, &user.EmailVerified, &user.Name, &user.Image,
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

		res, err := json.Marshal(user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if _, err = w.Write(res); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
