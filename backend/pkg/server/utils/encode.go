package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func Encode[T any](w http.ResponseWriter, status int, v T) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		return fmt.Errorf("failed to encode json: %v", err)
	}
	return nil
}

func EncodeNullStatusOK(w http.ResponseWriter) error {
	return Encode[any](w, http.StatusOK, nil)
}
