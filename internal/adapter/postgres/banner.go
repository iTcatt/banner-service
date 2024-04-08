package postgres

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5"

	"banner-service/internal/config"
	"banner-service/internal/model"
)

type Storage struct {
	conn *pgx.Conn
}

func NewStorage(cfg config.PostgresConfig) (*Storage, error) {
	var (
		conn *pgx.Conn
		err  error
	)
	dbPath := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database)

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	timeout := time.After(cfg.Timeout)

	for {
		select {
		case <-ticker.C:
			conn, err = pgx.Connect(context.Background(), dbPath)
			if err != nil {
				continue
			}
			err = conn.Ping(context.Background())
			if err != nil {
				continue
			}
			log.Println("successful database connection")

			return &Storage{conn: conn}, nil
		case <-timeout:
			return nil, fmt.Errorf("timed out waiting for database to become available")
		}
	}
}

func (s *Storage) CreateBanner(context.Context, model.Banner) error {
	return nil
}

func (s *Storage) CreateBannerTagsLocks(context.Context, int, []int) error {
	return nil
}

func (s *Storage) GetUserBanner(context.Context, int, int) (string, error) {
	return "", nil
}

func (s *Storage) GetFilteredBanners(context.Context, model.GetFilteredBannersParams) ([]model.Banner, error) {
	return nil, nil
}

// TODO: понять как поменять содержимое баннера, при этом не затерев его дату создания и id
func (s *Storage) PatchBanner(context.Context, model.Banner) error {
	return nil
}

func (s *Storage) PatchBannerTagsLocks(context.Context, int, []int) error {
	return nil
}

func (s *Storage) DeleteBanner(context.Context, int) error {
	return nil
}

func (s *Storage) Close(ctx context.Context) error {
	return s.conn.Close(ctx)
}
