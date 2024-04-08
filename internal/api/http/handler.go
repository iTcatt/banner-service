package http

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"strconv"

	"banner-service/internal/model"
)

type BannerService interface {
	GetUserBannerAction(context.Context, model.GetUserBannerParams) (string, error)
	GetFilteredBannersAction(context.Context, model.GetFilteredBannersParams) ([]model.Banner, error)

	CreateBannerAction(context.Context, model.BannerParams) (int, error)

	PatchBannerAction(context.Context, int, model.BannerParams) error

	DeleteBannerAction(context.Context, int) error
}

type Handler struct {
	service BannerService
}

func NewHandler(s BannerService) *Handler {
	return &Handler{service: s}
}

func (h *Handler) getUserBanner(w http.ResponseWriter, r *http.Request) error {
	tagID, err := strconv.Atoi(r.URL.Query().Get("tag_id"))
	if err != nil {
		return fmt.Errorf("%w: %s", ErrValidationFailed, err)
	}
	featureID, err := strconv.Atoi(r.URL.Query().Get("feature_id"))
	if err != nil {
		return fmt.Errorf("%w: %s", ErrValidationFailed, err)
	}
	useLastRevision := r.URL.Query().Get("use_last_revision") == "true"
	log.Printf("tagID: %d; featureID: %d; useLastRevision: %d\n", tagID, featureID, useLastRevision)

	result, err := h.service.GetUserBannerAction(r.Context(), model.GetUserBannerParams{
		TagID:           tagID,
		FeatureID:       featureID,
		UseLastRevision: useLastRevision,
	})
	if err != nil {
		return err
	}
	return sendJSONResponse(w, result, http.StatusOK)
}

func (h *Handler) getFilteredBanners(w http.ResponseWriter, r *http.Request) error {
	tagID, err := strconv.Atoi(r.URL.Query().Get("tag_id"))
	if err != nil {
		return fmt.Errorf("%w: %s", ErrValidationFailed, err)
	}
	featureID, err := strconv.Atoi(r.URL.Query().Get("feature_id"))
	if err != nil {
		return fmt.Errorf("%w: %s", ErrValidationFailed, err)
	}
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		return fmt.Errorf("%w: %s", ErrValidationFailed, err)
	}
	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil {
		return fmt.Errorf("%w: %s", ErrValidationFailed, err)
	}
	result, err := h.service.GetFilteredBannersAction(r.Context(), model.GetFilteredBannersParams{
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
	var params model.BannerParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		return fmt.Errorf("%w: %s", ErrValidationFailed, err)
	}

	id, err := h.service.CreateBannerAction(r.Context(), params)
	if err != nil {
		return err
	}
	return sendJSONResponse(w, id, http.StatusOK)
}

func (h *Handler) patchBanner(_ http.ResponseWriter, r *http.Request) error {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		return fmt.Errorf("%w: %s", ErrValidationFailed, err)
	}

	var params model.BannerParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		return fmt.Errorf("%w: %s", ErrValidationFailed, err)
	}

	if err := h.service.PatchBannerAction(r.Context(), id, params); err != nil {
		return err
	}
	return nil
}

func (h *Handler) deleteBanner(_ http.ResponseWriter, r *http.Request) error {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		return fmt.Errorf("%w: %s", ErrValidationFailed, err)
	}

	if err := h.service.DeleteBannerAction(r.Context(), id); err != nil {
		return err
	}
	return nil
}
