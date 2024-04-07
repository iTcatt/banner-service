package service

import (
	"context"

	"banner-service/internal/model"
)

type BannerStorage interface {
	Close(ctx context.Context) error
}

type BannerService struct {
	repo BannerStorage
}

func NewBannerService(repo BannerStorage) *BannerService {
	return &BannerService{repo: repo}
}

func (s *BannerService) GetUserBannerAction(ctx context.Context, params model.GetUserBannerParams) (model.UserBanner, error) {
	return model.UserBanner{}, nil
}

func (s *BannerService) GetBannerWithFiltersAction(ctx context.Context, params model.GetBannerWithFiltersParams) (model.UserBanner, error) {
	return model.UserBanner{}, nil
}
