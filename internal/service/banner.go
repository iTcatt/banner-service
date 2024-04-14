package service

import (
	"context"
	"errors"
	"log"
	"math/rand/v2"
	"time"

	"github.com/jackc/pgx/v5"

	"banner-service/internal/model"
)

type BannerStorage interface {
	GetUserBanner(ctx context.Context, tagID int, featureID int) (model.Banner, error)
	GetBannersByFeature(context.Context, int) ([]model.Banner, error)
	GetBannersByTag(context.Context, int) ([]model.Banner, error)
	GetTagsByBannerID(context.Context, int) ([]int, error)
	GetAllBanners(context.Context) ([]model.Banner, error)
	GetAllTags(context.Context) ([]int, error)

	CreateBanner(context.Context, model.Banner) error
	CreateTag(context.Context, int) error
	CreateBannerTagLock(context.Context, int, int) error

	PatchBanner(context.Context, model.Banner) error

	DeleteBanner(context.Context, int) error
	DeleteBannerTagsLocks(context.Context, int) error

	Close(context.Context) error
}

type Service struct {
	repo BannerStorage
}

func NewService(repo BannerStorage) *Service {
	return &Service{repo: repo}
}

// Сделать стуктуру, которая является map[pair(feature_id, tag_id)]model.Banner и время, на которое он был актуален
// Далее я делаю time.Now() и вычисляю разницу и если она меньше 5 минут, то возращаю из map, если больше, иду в базу
// и кэширую, должна быть таска на очищение неактуальных баннеров
// можно сделать pkg/casher в котором будет хранится интерфейс
func (s *Service) GetUserBannerAction(ctx context.Context, p model.GetUserBannerParams) (interface{}, error) {
	log.Println("running GetUserBannerAction")

	var (
		banner model.Banner
		err    error
	)

	log.Println("using last revision")
	banner, err = s.repo.GetUserBanner(ctx, p.TagID, p.FeatureID)
	if err != nil {
		return nil, err
	}

	if !banner.IsActive && !p.IsAdmin {
		return nil, nil
	}
	return banner.Content, nil
}

func (s *Service) GetFilteredBannersAction(
	ctx context.Context,
	p model.GetFilteredBannersParams,
) ([]model.BannerWithTags, error) {
	log.Println("running GetFilteredBannersAction")

	var banners []model.BannerWithTags
	if p.TagID != 0 && p.FeatureID != 0 {
		result, err := s.repo.GetUserBanner(ctx, p.TagID, p.FeatureID)
		if err != nil {
			return nil, err
		}
		banners = append(banners, model.BannerWithTags{
			Banner: result,
			Tags:   []int{p.TagID},
		})
	} else if p.TagID != 0 {
		result, err := s.repo.GetBannersByTag(ctx, p.TagID)
		if err != nil {
			return nil, err
		}
		for _, b := range result {
			banners = append(banners, model.BannerWithTags{
				Banner: b,
				Tags:   []int{p.TagID},
			})
		}
	} else if p.FeatureID != 0 {
		result, err := s.repo.GetBannersByFeature(ctx, p.FeatureID)
		if err != nil {
			return nil, err
		}
		for _, b := range result {
			tags, err := s.repo.GetTagsByBannerID(ctx, b.ID)
			if err != nil {
				return nil, err
			}
			banners = append(banners, model.BannerWithTags{
				Banner: b,
				Tags:   tags,
			})
		}
	} else {
		result, err := s.repo.GetAllBanners(ctx)
		if err != nil {
			return nil, err
		}
		for _, b := range result {
			tags, err := s.repo.GetTagsByBannerID(ctx, b.ID)
			if err != nil {
				return nil, err
			}
			banners = append(banners, model.BannerWithTags{
				Banner: b,
				Tags:   tags,
			})
		}
	}
	if len(banners) > p.Limit && p.Limit != -1 {
		banners = banners[:p.Limit]
	}
	if len(banners) < p.Offset {
		return []model.BannerWithTags{}, nil
	}
	return banners[p.Offset:], nil
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
	for _, tag := range p.TagIDs {
		if _, err := s.repo.GetUserBanner(ctx, tag, p.FeatureID); !errors.Is(err, pgx.ErrNoRows) {
			return 0, ErrAlreadyExists
		}
	}

	if err := s.repo.CreateBanner(ctx, banner); err != nil {
		return 0, err
	}
	for _, tag := range p.TagIDs {
		if err := s.repo.CreateTag(ctx, tag); err != nil {
			return 0, err
		}
		if err := s.repo.CreateBannerTagLock(ctx, banner.ID, tag); err != nil {
			return 0, err
		}
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
	if err := s.repo.PatchBanner(ctx, banner); err != nil {
		return err
	}
	if err := s.repo.DeleteBannerTagsLocks(ctx, id); err != nil {
		return err
	}
	for _, tag := range p.TagIDs {
		if err := s.repo.CreateTag(ctx, tag); err != nil {
			return err
		}
		if err := s.repo.CreateBannerTagLock(ctx, banner.ID, tag); err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) DeleteBannerAction(ctx context.Context, id int) error {
	log.Println("running DeleteBannerAction")
	if err := s.repo.DeleteBannerTagsLocks(ctx, id); err != nil {
		return err
	}

	if err := s.repo.DeleteBanner(ctx, id); err != nil {
		return err
	}
	return nil
}

func (s *Service) AuthAction(ctx context.Context) (int, error) {
	log.Println("running AuthAction")

	result, err := s.repo.GetAllTags(ctx)
	if err != nil {
		return 0, err
	}
	if len(result) == 0 {
		return rand.Int(), nil
	}
	return result[rand.IntN(len(result))], nil
}
