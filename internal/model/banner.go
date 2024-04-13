package model

import "time"

type Banner struct {
	ID        int         `json:"id"`
	FeatureID int         `json:"feature_id"`
	Content   interface{} `json:"content"`
	IsActive  bool        `json:"is_active"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

type BannerWithTags struct {
	Banner
	Tags []int `json:"tag_ids"`
}
