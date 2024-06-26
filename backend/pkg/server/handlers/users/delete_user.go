package handlers

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http"
	"rs/pkg/server/utils"
)

func DeleteUserHandler(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId := r.PathValue("id")
		if !utils.IsValidUUID(userId) {
			utils.HttpFormattedError(w, r, http.StatusBadRequest, "invalid user id", "invalid \"userId\"")
			return
		}

		sql := "delete from users where id=$1"

		commandTag, err := pool.Exec(context.Background(), sql, userId)
		if err != nil {
			utils.HttpFormattedError(w, r, http.StatusInternalServerError, "failed to delete the user", nil)
			return
		}

		if commandTag.RowsAffected() == 0 {
			utils.HttpFormattedError(w, r, http.StatusNotFound, "user not found", nil)
			return
		}
	}
}
