package postgres

import (
	"context"
	"database/sql"
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

// TODO: Заменить на миграции
func (s *Storage) Init(ctx context.Context) error {
	createBanner := `
		CREATE TABLE IF NOT EXISTS banner (
			banner_id BIGINT PRIMARY KEY,
			feature_id BIGINT NOT NULL,
			content JSON, 
			is_active BOOLEAN NOT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
		);`

	createTag := `CREATE TABLE IF NOT EXISTS tag (tag_id BIGINT PRIMARY KEY);`

	createBannerTag := `
		CREATE TABLE IF NOT EXISTS banner_tag (
		    banner_id BIGINT NOT NULL,
		    tag_id BIGINT NOT NULL,
		    FOREIGN KEY(banner_id) REFERENCES banner(banner_id) ON DELETE CASCADE,
		    FOREIGN KEY(tag_id) REFERENCES tag(tag_id)
		);`

	if _, err := s.conn.Exec(ctx, createBanner); err != nil {
		return err
	}
	log.Println("banner table created")
	if _, err := s.conn.Exec(ctx, createTag); err != nil {
		return err
	}
	log.Println("tag table created")
	if _, err := s.conn.Exec(ctx, createBannerTag); err != nil {
		return err
	}
	log.Println("banner_tag table created")
	return nil
}

func (s *Storage) CreateBanner(ctx context.Context, b model.Banner) error {
	log.Println("[DEBUG] db: create banner")
	q := `
		INSERT INTO banner (banner_id, feature_id, content, is_active, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6);`
	_, err := s.conn.Exec(ctx, q, b.ID, b.FeatureID, b.Content, b.IsActive, b.CreatedAt, b.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) CreateTag(ctx context.Context, id int) error {
	log.Println("[DEBUG] db: create tag")

	q := `INSERT INTO tag (tag_id) VALUES ($1);`
	if _, err := s.conn.Exec(ctx, q, id); err != nil {
		return err
	}
	return nil
}

func (s *Storage) CreateBannerTagLock(ctx context.Context, bannerID, tagID int) error {
	log.Println("[DEBUG] db: create banner tags locks")

	q := `INSERT INTO banner_tag (banner_id, tag_id) VALUES ($1, $2);`
	if _, err := s.conn.Exec(ctx, q, bannerID, tagID); err != nil {
		return err
	}
	return nil
}

func (s *Storage) GetUserBanner(ctx context.Context, tagID, featureID int) (model.Banner, error) {
	log.Println("[DEBUG] db: get user banner")
	q := `
		SELECT b.banner_id, b.feature_id, b.content, b.is_active, b.created_at, b.updated_at
		FROM banner b 
		JOIN banner_tag bt USING (banner_id) 
		JOIN tag t USING (tag_id)
		WHERE t.tag_id = $1 AND b.feature_id = $2;`

	row := s.conn.QueryRow(ctx, q, tagID, featureID)
	var b model.Banner
	if err := row.Scan(&b.ID, &b.FeatureID, &b.Content, &b.IsActive, &b.CreatedAt, &b.UpdatedAt); err != nil {
		return model.Banner{}, err
	}

	return b, nil
}

func (s *Storage) GetBannersByTag(ctx context.Context, tagID int) ([]model.Banner, error) {
	log.Println("[DEBUG] db: get user banners by tag")
	q := `
		SELECT b.banner_id, b.feature_id, b.content, b.is_active, b.created_at, b.updated_at
		FROM banner b 
		JOIN banner_tag bg USING (banner_id) 
		JOIN tag t USING (tag_id)
		WHERE t.tag_id = $1;`

	rows, err := s.conn.Query(ctx, q, tagID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	banners := make([]model.Banner, 0)
	for rows.Next() {
		var b model.Banner
		if err := rows.Scan(&b.ID, &b.FeatureID, &b.Content, &b.IsActive, &b.CreatedAt, &b.UpdatedAt); err != nil {
			return nil, err
		}
		banners = append(banners, b)
	}
	return banners, nil
}

func (s *Storage) GetBannersByFeature(ctx context.Context, featureID int) ([]model.Banner, error) {
	log.Println("[DEBUG] db: get filtered banners")

	q := `
		SELECT b.banner_id, b.feature_id, b.content, b.is_active, b.created_at, b.updated_at
		FROM banner b 
		WHERE b.feature_id = $1`

	rows, err := s.conn.Query(ctx, q, featureID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	banners := make([]model.Banner, 0)
	for rows.Next() {
		var b model.Banner
		if err := rows.Scan(&b.ID, &b.FeatureID, &b.Content, &b.IsActive, &b.CreatedAt, &b.UpdatedAt); err != nil {
			return nil, err
		}
		banners = append(banners, b)
	}

	return banners, nil
}

func (s *Storage) GetTagsByBannerID(ctx context.Context, bannerID int) ([]int, error) {
	log.Println("[DEBUG] db: get tags by id")

	q := `SELECT tag_id FROM tag t JOIN banner_tag bt USING (tag_id) WHERE bt.banner_id = $1`
	rows, err := s.conn.Query(ctx, q, bannerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	tags := make([]int, 0)
	for rows.Next() {
		var tag int
		if err := rows.Scan(&tag); err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}
	return tags, nil
}

func (s *Storage) GetAllBanners(ctx context.Context) ([]model.Banner, error) {
	log.Println("[DEBUG] db: get all banners")

	rows, err := s.conn.Query(ctx, `SELECT * FROM banner;`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	banners := make([]model.Banner, 0)
	for rows.Next() {
		var b model.Banner
		if err := rows.Scan(&b.ID, &b.FeatureID, &b.Content, &b.IsActive, &b.CreatedAt, &b.UpdatedAt); err != nil {
			return nil, err
		}
		banners = append(banners, b)
	}
	return banners, nil
}

func (s *Storage) GetAllTags(ctx context.Context) ([]int, error) {
	log.Println("[DEBUG] db: get all tags")

	rows, err := s.conn.Query(ctx, `SELECT tag_id FROM tag;`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	tags := make([]int, 0)
	for rows.Next() {
		var tag int
		if err := rows.Scan(&tag); err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}
	return tags, nil
}

func (s *Storage) PatchBanner(ctx context.Context, b model.Banner) error {
	log.Println("[DEBUG] db: patch banner")

	q := `
		UPDATE banner
		SET feature_id = $1,
			content = $2,
			is_active = $3,
			updated_at = $4
		WHERE banner_id = $5;`
	result, err := s.conn.Exec(ctx, q, b.FeatureID, b.Content, b.IsActive, b.UpdatedAt, b.ID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (s *Storage) DeleteBanner(ctx context.Context, id int) error {
	log.Println("[DEBUG] db: delete banner")

	result, err := s.conn.Exec(ctx, `DELETE FROM banner WHERE banner_id = $1;`, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (s *Storage) DeleteBannerTagsLocks(ctx context.Context, id int) error {
	log.Println("[DEBUG] db: delete banner tags locks")

	_, err := s.conn.Exec(ctx, `DELETE FROM banner_tag WHERE banner_id = $1;`, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) Close(ctx context.Context) error {
	log.Println("[DEBUG] db: close connection")

	return s.conn.Close(ctx)
}
