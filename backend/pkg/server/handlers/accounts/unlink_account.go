package handlers

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http"
)

func UnlinkAccountHandler(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		provider := r.PathValue("provider")
		providerAccountId := r.PathValue("providerAccountId")

		sql := "delete from accounts where provider=$1 and provider_account_id=$2"

		commandTag, err := pool.Exec(context.Background(), sql, provider, providerAccountId)
		if err != nil {
			http.Error(w, "Failed to unlink the account.", http.StatusInternalServerError)
			return
		}

		if commandTag.RowsAffected() == 0 {
			http.Error(w, "Account not found.", http.StatusNotFound)
			return
		}

		w.Write([]byte("null"))
	}
}
