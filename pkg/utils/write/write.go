package write

import (
	"encoding/json"
	"net/http"
)

func Write(obj interface{}, w http.ResponseWriter, status int) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(obj)
	if err != nil {
		return err
	}

	return nil
}
