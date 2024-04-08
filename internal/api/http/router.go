package http

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"banner-service/internal/adapter"
)

var (
	ErrValidationFailed = errors.New("validation failed")
)

func InitRouter(h *Handler) http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Get("/user_banner", registerHandler(h.getUserBanner))
	router.Get("/banner", registerHandler(h.getFilteredBanners))
	router.Post("/banner", registerHandler(h.createBanner))
	router.Patch("/banner/{id}", registerHandler(h.patchBanner))
	router.Delete("/banner/{id}", registerHandler(h.deleteBanner))

	return router
}

func registerHandler(handler func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := handler(w, r)

		switch {
		case err == nil:
			w.WriteHeader(http.StatusOK)
		case errors.Is(err, adapter.ErrNotFound):
			w.WriteHeader(http.StatusNotFound)
		case errors.Is(err, ErrValidationFailed):
			w.WriteHeader(http.StatusBadRequest)
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
