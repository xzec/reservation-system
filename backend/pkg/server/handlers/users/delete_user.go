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
			http.Error(w, "Invalid user id.", http.StatusBadRequest)
			return
		}

		sql := "delete from users where id=$1"

		commandTag, err := pool.Exec(context.Background(), sql, userId)
		if err != nil {
			http.Error(w, "Failed to delete the user.", http.StatusInternalServerError)
			return
		}

		if commandTag.RowsAffected() == 0 {
			http.Error(w, "User not found.", http.StatusNotFound)
			return
		}

		w.Write([]byte("null"))
	}
}
