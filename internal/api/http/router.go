package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func InitRouter(h *Handler) http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Get("/user_banner", errorsHandler(h.getUserBanner))
	router.Get("/banner", errorsHandler(h.getBannersWithFilter))
	router.Post("/banner", errorsHandler(h.createBanner))
	router.Patch("/banner/{id}", errorsHandler(h.updateBanner))
	router.Delete("/banner/{id}", errorsHandler(h.deleteBanner))

	return router
}
