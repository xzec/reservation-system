package handlers

import (
	"context"
	"encoding/json"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http"
	"rs/pkg/server/utils"
)

func DeleteUserHandler(pool *pgxpool.Pool) (handler func(http.ResponseWriter, *http.Request)) {
	return func(w http.ResponseWriter, r *http.Request) {
		userId := r.PathValue("id")
		if !utils.IsValidUUID(userId) {
			http.Error(w, "Invalid user id.", http.StatusBadRequest)
			return
		}

		sql := "delete from users where id=$1"

		commandTag, err := pool.Exec(context.Background(), sql, userId)
		if err != nil {
			http.Error(w, "Failed to delete user.", http.StatusInternalServerError)
			return
		}

		if commandTag.RowsAffected() == 0 {
			http.Error(w, "User not found.", http.StatusNotFound)
			return
		}

		nilResponse, _ := json.Marshal(nil)
		_, err = w.Write(nilResponse)
	}
}
