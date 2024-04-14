package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func InitRouter(h *Handler) http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.Logger)

	router.Get("/get_token", errorsMiddleware(h.getUserToken))
	router.Get("/admin", errorsMiddleware(h.getAdminToken))
	router.Get("/user_banner", errorsMiddleware(h.getUserBanner))
	router.Get("/banner", errorsMiddleware(h.getFilteredBanners))
	router.Post("/banner", errorsMiddleware(h.createBanner))
	router.Patch("/banner/{id}", errorsMiddleware(h.patchBanner))
	router.Delete("/banner/{id}", errorsMiddleware(h.deleteBanner))

	return router
}
