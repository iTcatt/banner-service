package actions

import (
	"banner-service/internal/storage"
	"context"

	"banner-service/internal/models"
)

type GetUserBannerParams struct {
	TagID           string
	FeatureID       string
	UseLastRevision bool
}

type GetBannersWithFiltersParams struct {
	TagID     string
	FeatureID string
	Limit     string
	Offset    string
}

type BannerService struct {
	repo storage.BannerStorage
}

func NewBannerService(repo storage.BannerStorage) *BannerService {
	return &BannerService{repo: repo}
}

func (s *BannerService) GetUserBannerAction(ctx context.Context, params GetUserBannerParams) (models.UserBanner, error) {
	return models.UserBanner{}, nil
}

func (s *BannerService) GetBannerWithFiltersAction(ctx context.Context, params GetBannersWithFiltersParams) (models.UserBanner, error) {
	return models.UserBanner{}, nil
}
