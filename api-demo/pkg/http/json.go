package http

import (
	"encoding/json"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, payload interface{}) {
	encoder := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	if err := encoder.Encode(payload); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func WriteError(w http.ResponseWriter, err error, code int) {
	encoder := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	payload := struct {
		Error string `json:"error"`
	}{
		err.Error(),
	}

	// TODO: the error should be typed to figure if it was a server error or an user error
	if err := encoder.Encode(payload); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(code)
	}
}
