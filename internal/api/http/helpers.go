package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"banner-service/internal/storage"
)

type wrappedHandler func(w http.ResponseWriter, r *http.Request) error

func errorsHandler(handler wrappedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := handler(w, r)

		switch {
		case err == nil:
			w.WriteHeader(http.StatusOK)
		case errors.Is(err, storage.ErrNotFound):
			w.WriteHeader(http.StatusNotFound)
		default:
			msg := fmt.Sprintf(`{"error":"%s"}`, err.Error())
			_ = sendJSONResponse(w, msg, http.StatusInternalServerError)
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
