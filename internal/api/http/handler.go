package http

import (
	"context"
	"log"
	"net/http"

	"banner-service/internal/model"
)

type Service interface {
	GetUserBannerAction(context.Context, model.GetUserBannerParams) (model.UserBanner, error)
	GetBannerWithFiltersAction(context.Context, model.GetBannerWithFiltersParams) (model.UserBanner, error)
}

type Handler struct {
	serv Service
}

func NewHandler(s Service) *Handler {
	return &Handler{serv: s}
}

func (h *Handler) getUserBanner(w http.ResponseWriter, r *http.Request) error {
	tagID := r.URL.Query().Get("tag_id")
	featureID := r.URL.Query().Get("feature_id")
	useLastRevision := r.URL.Query().Get("use_last_revision")
	log.Printf("tagID: %s; featureID: %s; useLastRevision: %s\n", tagID, featureID, useLastRevision)

	result, err := h.serv.GetUserBannerAction(r.Context(), model.GetUserBannerParams{
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

	result, err := h.serv.GetBannerWithFiltersAction(r.Context(), model.GetBannerWithFiltersParams{
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
