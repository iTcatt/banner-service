package service

import (
	"context"
	"log"
	"math/rand/v2"
	"time"

	"banner-service/internal/model"
)

type BannerStorage interface {
	GetUserBanner(context.Context, int, int) (string, error)
	GetFilteredBanners(context.Context, model.GetFilteredBannersParams) ([]model.Banner, error)

	CreateBanner(context.Context, model.Banner) error
	CreateBannerTagsLocks(context.Context, int, []int) error

	PatchBanner(context.Context, model.Banner) error
	PatchBannerTagsLocks(context.Context, int, []int) error

	DeleteBanner(context.Context, int) error

	Close(context.Context) error
}

type Service struct {
	repo BannerStorage
}

func NewService(repo BannerStorage) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetUserBannerAction(ctx context.Context, p model.GetUserBannerParams) (string, error) {
	log.Println("running GetUserBannerAction")

	var (
		content string
		err     error
	)
	if p.UseLastRevision {
		log.Println("using last revision")
		content, err = s.repo.GetUserBanner(ctx, p.TagID, p.FeatureID)
		if err != nil {
			return "", err
		}
	}
	return content, nil
}

func (s *Service) GetFilteredBannersAction(ctx context.Context, p model.GetFilteredBannersParams) ([]model.Banner, error) {
	log.Println("running GetFilteredBannersAction")
	result, err := s.repo.GetFilteredBanners(ctx, p)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *Service) CreateBannerAction(ctx context.Context, p model.BannerParams) (int, error) {
	log.Println("running CreateBannerAction")
	banner := model.Banner{
		ID:        rand.Int(),
		FeatureID: p.FeatureID,
		Content:   p.Content,
		IsActive:  p.IsActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err := s.repo.CreateBanner(ctx, banner)
	if err != nil {
		return 0, err
	}
	err = s.repo.CreateBannerTagsLocks(ctx, banner.ID, p.TagIDs)
	if err != nil {
		return 0, err
	}
	return banner.ID, nil
}

func (s *Service) PatchBannerAction(ctx context.Context, id int, p model.BannerParams) error {
	log.Println("running CreateBannerAction")
	banner := model.Banner{
		ID:        id,
		FeatureID: p.FeatureID,
		Content:   p.Content,
		IsActive:  p.IsActive,
		UpdatedAt: time.Now(),
	}
	err := s.repo.PatchBanner(ctx, banner)
	if err != nil {
		return err
	}
	err = s.repo.PatchBannerTagsLocks(ctx, id, p.TagIDs)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) DeleteBannerAction(ctx context.Context, id int) error {
	log.Println("running DeleteBannerAction")
	if err := s.repo.DeleteBanner(ctx, id); err != nil {
		return err
	}
	return nil
}
