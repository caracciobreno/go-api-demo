package http

import (
	"encoding/json"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, payload interface{}) {
	encoder := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := encoder.Encode(payload); err != nil {
		panic(err)
	}
}

func WriteError(w http.ResponseWriter, err error, code int) {
	encoder := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	payload := struct {
		Error string `json:"error"`
	}{
		err.Error(),
	}

	// TODO: the error should be typed to figure if it was a server error or an user error
	if err := encoder.Encode(payload); err != nil {
		panic(err)
	}
}
