package http

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"banner-service/internal/adapter"
)

type wrappedHandler func(w http.ResponseWriter, r *http.Request) error

func errorsHandler(handler wrappedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := handler(w, r)

		switch {
		case err == nil:
			w.WriteHeader(http.StatusOK)
		case errors.Is(err, adapter.ErrNotFound):
			w.WriteHeader(http.StatusNotFound)
		default:
			resp := struct {
				Error string `json:"error"`
			}{err.Error()}
			_ = sendJSONResponse(w, resp, http.StatusInternalServerError)
		}
	}
}

func sendJSONResponse(w http.ResponseWriter, result interface{}, status int) error {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(status)

	log.Printf("sending response")

	if err := json.NewEncoder(w).Encode(result); err != nil {
		return err
	}
	return nil
}
