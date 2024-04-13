package http

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5"

	"banner-service/internal/service"
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
			return
		case errors.Is(err, ErrValidationFailed) || errors.Is(err, service.ErrAlreadyExists):
			_ = sendJSONResponse(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		case errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows):
			w.WriteHeader(http.StatusNotFound)
		case errors.Is(err, service.ErrNoPermission):
			w.WriteHeader(http.StatusForbidden)
		default:
			_ = sendJSONResponse(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
		}
	}
}

func sendJSONResponse(w http.ResponseWriter, response interface{}, status int) error {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(status)

	log.Printf("sending response")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		return err
	}
	return nil
}
