package service

import (
	"context"

	"banner-service/internal/model"
)

type BannerStorage interface {
	Close(ctx context.Context) error
}

type Service struct {
	repo BannerStorage
}

func NewService(repo BannerStorage) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetUserBannerAction(ctx context.Context, params model.GetUserBannerParams) (model.UserBanner, error) {
	return model.UserBanner{}, nil
}

func (s *Service) GetBannerWithFiltersAction(ctx context.Context, params model.GetBannerWithFiltersParams) (model.UserBanner, error) {
	return model.UserBanner{}, nil
}
