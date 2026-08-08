package httpx

import (
	"encoding/json"
	"net/http"
)

func Write(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func Decode(w http.ResponseWriter, r *http.Request, target any) bool {
	if err := json.NewDecoder(r.Body).Decode(target); err != nil {
		Write(w, http.StatusBadRequest, map[string]string{"message": "El JSON enviado no es valido."})
		return false
	}
	return true
}

func Message(w http.ResponseWriter, status int, message string) {
	Write(w, status, map[string]string{"message": message})
}
