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

func GetUserHandler(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId := r.PathValue("id")
		if !utils.IsValidUUID(userId) {
			utils.HttpFormattedError(w, r, http.StatusBadRequest, "invalid user id", "invalid \"userId\"")
			return
		}

		sql := `select id, email, email_verified, name, image
from users
where id = $1`

		var user models.User
		err := pool.QueryRow(context.Background(), sql, userId).Scan(&user.Id, &user.Email, &user.EmailVerified, &user.Name, &user.Image)
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
