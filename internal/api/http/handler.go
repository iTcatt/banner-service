package http

import (
	"log"
	"net/http"

	"banner-service/internal/actions"
)

type Handler struct {
	service *actions.BannerService
}

func NewHandler(service *actions.BannerService) *Handler {
	return &Handler{service: service}
}

// useLastRevision убрал из routing
func (h *Handler) getUserBanner(w http.ResponseWriter, r *http.Request) error {
	tagID := r.URL.Query().Get("tag_id")
	featureID := r.URL.Query().Get("feature_id")
	useLastRevision := r.URL.Query().Get("use_last_revision")
	log.Printf("tagID: %s; featureID: %s; useLastRevision: %s\n", tagID, featureID, useLastRevision)

	result, err := h.service.GetUserBannerAction(r.Context(), actions.GetUserBannerParams{
		TagID:           tagID,
		FeatureID:       featureID,
		UseLastRevision: useLastRevision == "true",
	})
	if err != nil {
		return err
	}
	return sendJSONResponse(w, result, http.StatusOK)
}

func (h *Handler) getBannersWithFilter(w http.ResponseWriter, r *http.Request) error {
	tagID := r.URL.Query().Get("tag_id")
	featureID := r.URL.Query().Get("feature_id")
	limit := r.URL.Query().Get("limit")
	offset := r.URL.Query().Get("offset")

	result, err := h.service.GetBannerWithFiltersAction(r.Context(), actions.GetBannersWithFiltersParams{
		TagID:     tagID,
		FeatureID: featureID,
		Limit:     limit,
		Offset:    offset,
	})
	if err != nil {
		return err
	}
	return sendJSONResponse(w, result, http.StatusOK)
}

func (h *Handler) createBanner(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (h *Handler) updateBanner(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (h *Handler) deleteBanner(w http.ResponseWriter, r *http.Request) error {
	return nil
}
