package http

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"banner-service/internal/model"
)

type BannerService interface {
	GetUserBannerAction(context.Context, model.GetUserBannerParams) (interface{}, error)
	GetFilteredBannersAction(context.Context, model.GetFilteredBannersParams) ([]model.BannerWithTags, error)

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
	var (
		params model.GetUserBannerParams
		err    error
	)
	// TODO: проверка токенов авторизация и т.д.
	params.IsAdmin = true

	tagID := r.URL.Query().Get("tag_id")
	if tagID == "" {
		return fmt.Errorf("%w: tag_id is required", ErrValidationFailed)
	}
	if params.TagID, err = strconv.Atoi(tagID); err != nil {
		return fmt.Errorf("%w: %s", ErrValidationFailed, err)
	}

	featureID := r.URL.Query().Get("feature_id")
	if featureID == "" {
		return fmt.Errorf("%w: feature_id is required", ErrValidationFailed)
	}
	if params.FeatureID, err = strconv.Atoi(featureID); err != nil {
		return fmt.Errorf("%w: %s", ErrValidationFailed, err)
	}

	useLastRevision := r.URL.Query().Get("use_last_revision")
	if useLastRevision == "true" {
		params.UseLastRevision = true
	} else if useLastRevision == "false" {
		params.UseLastRevision = false
	} else {
		return fmt.Errorf("%w: use_last_revision is bool", ErrValidationFailed)
	}

	log.Printf("tagID: %s; featureID: %s; useLastRevision: %s\n", tagID, featureID, useLastRevision)

	result, err := h.service.GetUserBannerAction(r.Context(), params)
	if err != nil {
		return err
	}
	return sendJSONResponse(w, result, http.StatusOK)
}

func (h *Handler) getFilteredBanners(w http.ResponseWriter, r *http.Request) error {
	var (
		params model.GetFilteredBannersParams
		err    error
	)

	tagID := r.URL.Query().Get("tag_id")
	if tagID != "" {
		if params.TagID, err = strconv.Atoi(tagID); err != nil {
			return fmt.Errorf("%w: %v", ErrValidationFailed, err)
		}
	}

	featureID := r.URL.Query().Get("feature_id")
	if featureID != "" {
		if params.FeatureID, err = strconv.Atoi(featureID); err != nil {
			return fmt.Errorf("%w: %v", ErrValidationFailed, err)
		}
	}

	limit := r.URL.Query().Get("limit")
	if limit != "" {
		if params.Limit, err = strconv.Atoi(limit); err != nil || params.Limit < 0 {
			return fmt.Errorf("%w: %v", ErrValidationFailed, err)
		}
	} else {
		params.Limit = -1
	}

	offset := r.URL.Query().Get("offset")
	if offset != "" {
		if params.Offset, err = strconv.Atoi(offset); err != nil || params.Offset < 0 {
			return fmt.Errorf("%w: %v", ErrValidationFailed, err)
		}
	}
	log.Printf("tagID: %s; featureID: %s; limit: %s; offset: %s\n", tagID, featureID, limit, offset)
	result, err := h.service.GetFilteredBannersAction(r.Context(), params)
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

	return sendJSONResponse(w, id, http.StatusCreated)
}

func (h *Handler) patchBanner(_ http.ResponseWriter, r *http.Request) error {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		return fmt.Errorf("%w: %v", ErrValidationFailed, err)
	}

	var params model.BannerParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		return fmt.Errorf("%w: %v", ErrValidationFailed, err)
	}

	if err := h.service.PatchBannerAction(r.Context(), id, params); err != nil {
		return err
	}
	return nil
}

func (h *Handler) deleteBanner(w http.ResponseWriter, r *http.Request) error {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		return fmt.Errorf("%w: %v", ErrValidationFailed, err)
	}

	if err := h.service.DeleteBannerAction(r.Context(), id); err != nil {
		return err
	}
	w.WriteHeader(http.StatusNoContent)
	return nil
}
